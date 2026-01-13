package log_monitor

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// LogMonitor écoute les logs du serveur Hytale.
type LogMonitor struct {
	LogPath     string
	DiscordChan chan<- string
	CommandChan chan<- string
}

// NewLogMonitor crée une nouvelle instance de LogMonitor.
func NewLogMonitor(logPath string, discordChan chan<- string, commandChan chan<- string) *LogMonitor {
	return &LogMonitor{
		LogPath:     logPath,
		DiscordChan: discordChan,
		CommandChan: commandChan,
	}
}

// Start commence à surveiller les logs.
func (lm *LogMonitor) Start() {
	file, err := os.Open(lm.LogPath)
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture du fichier de log: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line) // Afficher le log dans la console

		// Envoyer les logs à Discord
		lm.DiscordChan <- line

		// Détecter des commandes spécifiques dans les logs
		if strings.Contains(line, "Server started") {
			lm.CommandChan <- "server_started"
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Erreur lors de la lecture des logs: %v", err)
	}
}
