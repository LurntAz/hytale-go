package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LurntAz/hytale-go/internal/modules/commands"
	"github.com/LurntAz/hytale-go/internal/modules/discord"
)

// StartServer démarre le serveur HTTP pour interagir avec l'outil.
func StartServer(commandManager *commands.CommandManager, discordManager *discord.DiscordManager) {
	http.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		command := r.URL.Query().Get("command")
		var output string
		var err error

		switch command {
		case "update":
			output, err = commandManager.UpdateServer()
		case "restart":
			err = commandManager.RestartServer()
			output = "Serveur redémarré avec succès."
		default:
			output = "Commande inconnue."
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Envoyer un message via Discord
		if discordManager.WebhookURL != "" {
			err = discordManager.SendMessage(fmt.Sprintf("Commande exécutée: %s, Sortie: %s", command, output))
			if err != nil {
				log.Printf("Erreur lors de l'envoi du message Discord: %v", err)
			}
		}

		fmt.Fprintf(w, "Sortie: %s", output)
	})

	log.Println("Serveur démarré sur :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
