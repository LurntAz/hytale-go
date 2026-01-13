package domains

// Liste des motifs de logs intéressants à envoyer à Discord
var InterestingLogs = []string{
	"Hytale Server Booted!",            // Démarrage du serveur
	"[ServerManager|P] Listening on /", // Port d'écoute
	"[World|default] Player ",          // Connexion d'un joueur
	"[PlayerSystems] Removing player",  // Déconnexion d'un joueur
	"ERROR]",                           // Erreurs critiques
	"Shutting down server",             // Arrêt du serveur
	"Update available",                 // Mise à jour disponible
}
