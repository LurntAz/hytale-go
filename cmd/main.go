package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LurntAz/hytale-go/internal/domains"
	loghandler "github.com/LurntAz/hytale-go/internal/handler/log_handler"
	"github.com/LurntAz/hytale-go/internal/modules/commands"
	"github.com/LurntAz/hytale-go/internal/modules/discord"
	"github.com/LurntAz/hytale-go/internal/modules/server"
)

func main() {
	// Définition des flags
	credentialsPath := flag.String("credentials-path", "/home/lurnt/Téléchargements/hytale-downloader/.hytale-downloader-credentials.json", "Chemin vers le fichier de credentials")
	downloadPath := flag.String("download-path", "", "Chemin pour télécharger le fichier zip")
	patchline := flag.String("patchline", "release", "Patchline à utiliser")
	serverPath := flag.String("server-path", "", "Chemin vers le dossier du serveur Hytale")
	webhookURL := flag.String("webhook-url", "", "URL du webhook Discord")

	flag.Parse()

	// Initialisation des modules
	commandManager := commands.NewCommandManager(*credentialsPath, *downloadPath, *patchline, *serverPath)
	discordManager := discord.NewDiscordManager(*webhookURL)

	// Télécharger et extraire le serveur Hytale
	err := commandManager.DownloadAndExtractServer()
	if err != nil {
		log.Fatalf("Erreur lors du téléchargement et de l'extraction du serveur: %v\n", err)
	}

	// Vérifier les versions
	currentVersion, latestVersion, err := commandManager.CheckVersion()
	if err != nil {
		log.Printf("Erreur lors de la vérification des versions: %v\n", err)
	} else {
		log.Printf("Version actuelle: %s\n", currentVersion)
		log.Printf("Dernière version disponible: %s\n", latestVersion)
	}

	// Démarrer le serveur Hytale en mode console
	serverCmd, logChan, stdin, err := commandManager.StartServer()
	if err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v\n", err)
	}
	log.Println("Serveur Hytale démarré en mode console. Appuyez sur Ctrl+C pour arrêter.")

	// Initialiser le gestionnaire de logs
	logHandler := loghandler.NewLogHandler(discordManager, domains.InterestingLogs)

	// Simuler un canal de logs (à remplacer par votre source réelle)
	logHandler.HandleLogs(logChan)

	// Écouter les entrées utilisateur pour envoyer des commandes au serveur
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			command := scanner.Text()
			// Envoyer la commande au serveur
			_, err := fmt.Fprintf(stdin, "%s\n", command)
			if err != nil {
				log.Printf("Erreur lors de l'envoi de la commande: %v\n", err)
			}
		}
	}()

	// Démarrer le serveur HTTP pour interagir avec l'outil
	go server.StartServer(commandManager, discordManager)

	// Attendre un signal pour arrêter le programme
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Arrêter le serveur Hytale
	if err := serverCmd.Process.Kill(); err != nil {
		log.Printf("Erreur lors de l'arrêt du serveur: %v\n", err)
	}

	fmt.Println("\nArrêt du service...")
}
