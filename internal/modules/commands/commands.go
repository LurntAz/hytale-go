package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
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
