package commands

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

// CommandManager gère l'exécution des commandes Hytale.
type CommandManager struct {
	CredentialsPath string
	DownloadPath    string
	Patchline       string
	ServerPath      string
}

// NewCommandManager crée une nouvelle instance de CommandManager.
func NewCommandManager(credentialsPath, downloadPath, patchline, serverPath string) *CommandManager {
	return &CommandManager{
		CredentialsPath: credentialsPath,
		DownloadPath:    downloadPath,
		Patchline:       patchline,
		ServerPath:      serverPath,
	}
}

// ExecuteCommand exécute une commande Hytale.
func (cm *CommandManager) ExecuteCommand(args ...string) (string, error) {
	cmd := exec.Command("./hytale-downloader-linux-amd64", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("erreur lors de l'exécution de la commande: %v, sortie: %s", err, out.String())
	}
	return out.String(), nil
}

// UpdateServer vérifie et télécharge les mises à jour.
func (cm *CommandManager) UpdateServer() (string, error) {
	return cm.ExecuteCommand("-download-path", cm.DownloadPath, "-patchline", cm.Patchline)
}

// RestartServer redémarre le serveur Hytale.
func (cm *CommandManager) RestartServer() error {
	// Arrêter le serveur
	stopCmd := exec.Command("pkill", "-f", "java -jar "+cm.ServerPath+"/Server/HytaleServer.jar --assets "+cm.ServerPath+"/Assets.zip")
	err := stopCmd.Run()
	if err != nil {
		return fmt.Errorf("erreur lors de l'arrêt du serveur: %v", err)
	}

	// Redémarrer le serveur
	startCmd := exec.Command("java", "-jar", cm.ServerPath+"/HytaleServer.jar", "--assets", cm.ServerPath+"/Assets.zip")
	err = startCmd.Start()
	if err != nil {
		return fmt.Errorf("erreur lors du redémarrage du serveur: %v", err)
	}

	return nil
}

// CheckVersion vérifie la version actuelle et la version disponible.
func (cm *CommandManager) CheckVersion() (currentVersion string, latestVersion string, err error) {
	// Vérifier la version actuelle
	currentCmd := exec.Command("./hytale-downloader-linux-amd64", "-print-version")
	var currentOut bytes.Buffer
	currentCmd.Stdout = &currentOut
	err = currentCmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("erreur lors de la vérification de la version actuelle: %v", err)
	}
	currentVersion = currentOut.String()

	// Vérifier la dernière version disponible
	latestCmd := exec.Command("./hytale-downloader-linux-amd64", "-check-update")
	var latestOut bytes.Buffer
	latestCmd.Stdout = &latestOut
	err = latestCmd.Run()
	if err != nil {
		return currentVersion, "", fmt.Errorf("erreur lors de la vérification de la dernière version: %v", err)
	}
	latestVersion = latestOut.String()

	return currentVersion, latestVersion, nil
}

// StartServer démarre le serveur Hytale en mode console et retourne un canal pour les logs et un writer pour envoyer des commandes.
func (cm *CommandManager) StartServer() (*exec.Cmd, chan string, io.WriteCloser, error) {
	cmd := exec.Command("java", "-jar", cm.ServerPath+"/Server/HytaleServer.jar", "--assets", cm.ServerPath+"/Assets.zip")

	// Créer un pipe pour capturer le stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("erreur lors de la création du pipe stdout: %v", err)
	}

	// Créer un pipe pour envoyer des commandes via stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("erreur lors de la création du pipe stdin: %v", err)
	}

	// Démarrer le processus
	err = cmd.Start()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("erreur lors du démarrage du serveur: %v", err)
	}

	// Créer un canal pour envoyer les logs
	logChan := make(chan string)

	// Lire le stdout dans une goroutine
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			logChan <- line
		}
	}()

	return cmd, logChan, stdin, nil
}

// DownloadAndExtractServer télécharge et extrait le serveur Hytale.
func (cm *CommandManager) DownloadAndExtractServer() error {
	// Créer le dossier serverPath s'il n'existe pas
	if err := os.MkdirAll(cm.ServerPath, os.ModePerm); err != nil {
		return fmt.Errorf("erreur lors de la création du dossier %s: %v", cm.ServerPath, err)
	}

	// Vérifier si le serveur est déjà installé
	serverDir := filepath.Join(cm.ServerPath, "Server")
	if _, err := os.Stat(serverDir); !os.IsNotExist(err) {
		log.Println("Le serveur est déjà installé dans", serverDir)
		return nil
	}

	// Télécharger le fichier ZIP
	zipURL := "https://downloader.hytale.com/hytale-downloader.zip"
	resp, err := http.Get(zipURL)
	if err != nil {
		return fmt.Errorf("erreur lors du téléchargement du fichier ZIP: %v", err)
	}
	defer resp.Body.Close()

	// Créer un fichier temporaire pour le ZIP
	zipFile, err := os.CreateTemp("", "hytale-downloader-*.zip")
	if err != nil {
		return fmt.Errorf("erreur lors de la création du fichier temporaire: %v", err)
	}
	defer os.Remove(zipFile.Name())

	// Écrire le contenu du ZIP dans le fichier temporaire
	_, err = io.Copy(zipFile, resp.Body)
	if err != nil {
		return fmt.Errorf("erreur lors de l'écriture du fichier ZIP: %v", err)
	}
	zipFile.Close()

	// Extraire le ZIP dans le répertoire du serveur
	err = cm.unzip(zipFile.Name(), "")
	if err != nil {
		return fmt.Errorf("erreur lors de l'extraction du fichier ZIP: %v", err)
	}

	// Exécuter la commande pour télécharger le serveur Hytale
	err = cm.ExecuteCommandWithOutput("-download-path", filepath.Join(cm.ServerPath, "game.zip"), "-patchline", cm.Patchline)
	if err != nil {
		return fmt.Errorf("erreur lors du téléchargement du serveur Hytale: %v", err)
	}

	// Extraire le fichier game.zip
	err = cm.unzip(filepath.Join(cm.ServerPath, "game.zip"), cm.ServerPath)
	if err != nil {
		return fmt.Errorf("erreur lors de l'extraction du fichier game.zip: %v", err)
	}

	return nil
}

// unzip extrait un fichier ZIP dans un répertoire
// unzip extrait un fichier ZIP dans un répertoire avec une barre de progression
func (cm *CommandManager) unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// Créer une barre de progression
	bar := progressbar.Default(int64(len(r.File)))

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			bar.Add(1)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}

		// Mettre à jour la barre de progression
		bar.Add(1)
	}

	return nil
}

// ExecuteCommandWithOutput exécute une commande Hytale et affiche le retour en temps réel.
func (cm *CommandManager) ExecuteCommandWithOutput(args ...string) error {
	cmd := exec.Command("./hytale-downloader-linux-amd64", args...)

	// Créer des pipes pour capturer stdout et stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("erreur lors de la création du pipe stdout: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("erreur lors de la création du pipe stderr: %v", err)
	}

	// Démarrer le processus
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("erreur lors du démarrage de la commande: %v", err)
	}

	// Lire stdout en temps réel
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Lire stderr en temps réel
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Attendre la fin de la commande
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("erreur lors de l'exécution de la commande: %v", err)
	}

	return nil
}
