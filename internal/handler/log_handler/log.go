package loghandler

import (
	"fmt"
	"log"
	"strings"

	"github.com/LurntAz/hytale-go/internal/modules/discord"
)

// LogHandler gère les logs du serveur Hytale
type LogHandler struct {
	DiscordManager  *discord.DiscordManager
	InterestingLogs []string
}

// NewLogHandler crée une nouvelle instance de LogHandler
func NewLogHandler(dm *discord.DiscordManager, interestingLogs []string) *LogHandler {
	return &LogHandler{
		DiscordManager:  dm,
		InterestingLogs: interestingLogs,
	}
}

// HandleLogs écoute les logs et envoie les messages pertinents à Discord
func (lh *LogHandler) HandleLogs(logChan <-chan string) {
	go func() {
		for line := range logChan {
			fmt.Println(line)
			if lh.isInterestingLog(line) {
				formattedMessage := discord.FormatLogAsEmbed(line)
				err := lh.DiscordManager.SendEmbed(formattedMessage)
				if err != nil {
					log.Printf("Erreur lors de l'envoi du message Discord: %v", err)
				}
			}
		}
	}()
}

// isInterestingLog vérifie si un log est intéressant
func (lh *LogHandler) isInterestingLog(line string) bool {
	for _, pattern := range lh.InterestingLogs {
		if strings.Contains(line, pattern) {
			return true
		}
	}
	return false
}
