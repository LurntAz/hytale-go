# **Hytale Go (Hytale Manager)**

**Hytale Go** is a Go-based tool for managing a Hytale server. It allows you to start, monitor, update, and restart the Hytale server while sending logs and notifications via Discord.

---

## **Features**

- **Version Management**: Checks the current and latest available versions of the Hytale server.
- **Server Startup**: Launches the Hytale server in console mode and displays real-time logs.
- **Real-Time Interaction**: Sends commands directly to the server via the terminal.
- **Log Monitoring**: Captures server logs and sends them to a Discord webhook.
- **Update & Restart**: Updates and restarts the server via HTTP commands.

---

## **Prerequisites**

- **Go 1.16+**: To compile and run the tool.
- **Java 25**: To run the Hytale server.
- **Hytale Downloader**: To download Hytale server files.
- **Discord Webhook**: To receive log notifications.

---

## **Installation**

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-username/hytale-manager.git
   cd hytale-manager
   ```

2. **Build the project**:
   ```bash
   go build -o hytale-manager
   ```

---

## **Usage**

### **Start the Tool**

```bash
./hytale-manager \
  --credentials-path /path/to/credentials.json \
  --download-path /path/to/download \
  --patchline release \
  --server-path /path/to/server \
  --webhook-url YOUR_DISCORD_WEBHOOK_URL
```

### **HTTP Commands**

- **Update the server**:
  ```bash
  curl "http://localhost:8080/execute?command=update"
  ```

- **Restart the server**:
  ```bash
  curl "http://localhost:8080/execute?command=restart"
  ```

### **Real-Time Interaction**

- Once the server is running, you can enter commands directly in the terminal to interact with the Hytale server.

---

## **Project Structure**

```
/hytale-go
  ├── cmd/
  │   └── main.go
  ├── internal/module/
  │   ├── commands/
  │   │   └── command_manager.go
  │   ├── discord/
  │   │   └── discord_manager.go
  │   └── server/
  │       └── server.go
  └── README.md
```

---

## **Contributing**

Contributions are welcome! Feel free to open an **issue** or **pull request** to suggest improvements.

---

## **License**

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.

---

## **Authors**

- **Lurnt Az** – Lead Developer at **Ziroh**.

---

## **Acknowledgments**

- Thanks to the **Hytale community** for their support and documentation.
- Thanks to **Mistral AI** and its assistant **Le Chat** for technical support and assistance in developing this tool.

---

### **Instructions for Use**

1. **Customize** paths and URLs according to your setup.
2. **Add a license** if needed.
3. **Include additional instructions** if you have specific features to document.
