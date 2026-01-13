package discord

import (
	"fmt"
	"strings"
)

// FormatLogForDiscord formate un log pour l'envoyer à Discord
func FormatLogForDiscord(log string) string {
	if strings.Contains(log, "Hytale Server Booted!") {
		return "**🚀 Serveur Hytale démarré !**"
	} else if strings.Contains(log, "[ServerManager|P] Listening on /") {
		port := extractPort(log)
		return fmt.Sprintf("**🌐 Serveur en écoute sur le port :** `%s`", port)
	} else if strings.Contains(log, "[World|default] Player ") {
		playerName := extractPlayerName(log)
		return fmt.Sprintf("**👤 Nouveau joueur connecté :** `%s`", playerName)
	} else if strings.Contains(log, "[PlayerSystems] Removing player") {
		playerName := extractPlayerName(log)
		return fmt.Sprintf("**👋 Joueur déconnecté :** `%s`", playerName)
	} else if strings.Contains(log, "ERROR]") {
		errorMessage := extractErrorMessage(log)
		return fmt.Sprintf("**❌ Erreur critique :** ```%s```", errorMessage)
	} else if strings.Contains(log, "Shutting down server") {
		return "**🛑 Arrêt du serveur Hytale...**"
	} else if strings.Contains(log, "Update available") {
		return "**🔄 Mise à jour disponible !**"
	}
	return fmt.Sprintf("```%s```", log)
}

// formatLogAsEmbed formate un log en embed Discord
func FormatLogAsEmbed(log string) Embed {
	if strings.Contains(log, "Hytale Server Booted!") {
		return Embed{
			Title:       "🚀 Serveur Hytale démarré !",
			Description: "Le serveur Hytale a démarré avec succès.",
			Color:       0x00ff00,
			Footer:      EmbedFooter{Text: "Hytale Manager • Ziroh"},
		}
	} else if strings.Contains(log, "[ServerManager|P] Listening on /") {
		port := extractPort(log)
		return Embed{
			Title:       "🌐 Serveur en écoute",
			Description: fmt.Sprintf("Le serveur écoute sur le port **%s**.", port),
			Color:       0x00ffff,
			Footer:      EmbedFooter{Text: "Hytale Manager • Ziroh"},
		}
	} else if strings.Contains(log, "[World|default] Player ") {
		playerName := extractPlayerName(log)
		return Embed{
			Title:       "👤 Nouveau joueur connecté",
			Description: fmt.Sprintf("**%s** a rejoint le serveur.", playerName),
			Color:       0x00ff00, // Vert
			Footer: EmbedFooter{
				Text: "Hytale Manager • Ziroh",
			},
			Thumbnail: EmbedImage{
				URL: "https://example.com/avatar.png", // URL d'un avatar ou icône
			},
		}
	} else if strings.Contains(log, "[PlayerSystems] Removing player") {
		playerName := extractPlayerName(log)
		return Embed{
			Title:       "👋 Joueur déconnecté",
			Description: fmt.Sprintf("**%s** a quitté le serveur.", playerName),
			Color:       0xff8c00, // Orange
			Footer: EmbedFooter{
				Text: "Hytale Manager • Ziroh",
			},
		}
	}
	// ... (autres cas)
	return Embed{
		Description: fmt.Sprintf("```%s```", log),
		Color:       0x808080, // Gris
		Footer:      EmbedFooter{Text: "Hytale Manager • Ziroh"},
	}
}

// extractPort extrait le port du log "[ServerManager|P] Listening on /0.0.0.0:5520"
func extractPort(log string) string {
	parts := strings.Split(log, "/0.0.0.0:")
	if len(parts) > 1 {
		port := strings.Split(parts[1], " ")[0]
		return port
	}
	return "inconnu"
}

// extractPlayerName extrait le nom du joueur du log "[World|default] Adding player 'Lurnt'..."
func extractPlayerName(log string) string {
	parts := strings.Split(log, "'")
	if len(parts) > 1 {
		return parts[1]
	}
	return "inconnu"
}

// extractErrorMessage extrait le message d'erreur du log "[ERROR] Message d'erreur..."
func extractErrorMessage(log string) string {
	parts := strings.Split(log, "] ")
	if len(parts) > 1 {
		return parts[1]
	}
	return "Message d'erreur non identifié"
}
