package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// DiscordManager gère l'envoi de messages via les webhooks Discord.
type DiscordManager struct {
	WebhookURL string
}

// NewDiscordManager crée une nouvelle instance de DiscordManager.
func NewDiscordManager(webhookURL string) *DiscordManager {
	return &DiscordManager{WebhookURL: webhookURL}
}

// SendMessage envoie un message via un webhook Discord.
func (dm *DiscordManager) SendMessage(message string) error {
	if dm.WebhookURL == "" {
		log.Println("Aucune URL de webhook Discord configurée.")
		return nil
	}

	payload := map[string]string{"content": message}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = http.Post(dm.WebhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi du message Discord: %v", err)
	}

	return nil
}

// SendEmbed envoie un embed à Discord
func (dm *DiscordManager) SendEmbed(embed Embed) error {
	if dm.WebhookURL == "" {
		log.Println("Aucune URL de webhook Discord configurée.")
		return nil
	}

	message := DiscordEmbedMessage{
		Embeds: []Embed{embed},
	}
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = http.Post(dm.WebhookURL, "application/json", bytes.NewBuffer(payload))
	return err
}
