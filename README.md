# hytale-go (Hytale Manager ?)

**Hytale Go** est un outil écrit en Go pour gérer un serveur Hytale. Il permet de démarrer, surveiller, mettre à jour et redémarrer le serveur Hytale, tout en envoyant les logs et notifications via Discord.

---

## Fonctionnalités

- **Gestion des versions** : Vérifie la version actuelle et la dernière version disponible du serveur Hytale.
- **Démarrage du serveur** : Démarre le serveur Hytale en mode console et affiche les logs en temps réel.
- **Interaction en temps réel** : Permet d'envoyer des commandes directement au serveur via le terminal.
- **Surveillance des logs** : Capture les logs du serveur et les envoie à un webhook Discord.
- **Mise à jour et redémarrage** : Permet de mettre à jour et redémarrer le serveur via des commandes HTTP.

---

## Prérequis

- **Go 1.16+** : Pour compiler et exécuter l'outil.
- **Java 25** : Pour exécuter le serveur Hytale.
- **Hytale Downloader** : Pour télécharger les fichiers du serveur Hytale.
- **Webhook Discord** : Pour recevoir les notifications des logs.

---

## Installation

1. **Cloner le dépôt** :
   ```bash
   git clone https://github.com/votre-utilisateur/hytale-manager.git
   cd hytale-manager
   ```

2. **Construire le projet** :
   ```bash
   go build -o hytale-manager
   ```

---

## Utilisation

### Démarrer l'outil

```bash
./hytale-manager \
  --credentials-path /chemin/vers/credentials.json \
  --download-path /chemin/vers/téléchargement \
  --patchline release \
  --server-path /chemin/vers/serveur \
  --webhook-url VOTRE_URL_WEBHOOK_DISCORD
```

### Commandes HTTP

- **Mettre à jour le serveur** :
  ```bash
  curl "http://localhost:8080/execute?command=update"
  ```

- **Redémarrer le serveur** :
  ```bash
  curl "http://localhost:8080/execute?command=restart"
  ```

### Interaction en temps réel

- Une fois le serveur démarré, vous pouvez entrer des commandes directement dans le terminal pour interagir avec le serveur Hytale.

---

## Structure du Projet

```
/hytale-manager
  ├── main.go
  ├── commands/
  │   └── command_manager.go
  ├── discord/
  │   └── discord_manager.go
  ├── server/
  │   └── server.go
  └── README.md
```

---

## Contribution

Les contributions sont les bienvenues ! N'hésitez pas à ouvrir une *issue* ou une *pull request* pour proposer des améliorations.

---

## Licence

Ce projet est sous licence MIT. Voir le fichier [LICENSE](LICENSE) pour plus de détails.

---

## Auteurs

- **Lurnt Az** - Développeur principal chez **Ziroh**.

---

## Remerciements

- Merci à la communauté Hytale pour leur support et leur documentation.
- Merci à **Mistral AI** et son assistant **Le Chat** pour leur soutien technique et leur aide précieuse dans le développement de cet outil.

---

### Instructions pour l'utiliser

1. **Personnalisez** les chemins et les URLs selon votre configuration.
2. **Ajoutez une licence** si nécessaire.
3. **Ajoutez des instructions supplémentaires** si vous avez des fonctionnalités spécifiques à documenter.
