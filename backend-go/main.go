package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
)

var db *sql.DB

type Message struct {
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

func handleSendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var message Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Insert message into mysql daabase
	timestamp := time.Now()
	_, err := db.Exec("INSERT INTO messages (text, timestamp) VALUES (?, ?)", message.Text, timestamp)
	if err != nil {
		http.Error(w, "Failed to insert message into database", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Message received: %s", message.Text)))
}

func handleGetMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	log.Println("Received a request to get messages in the backend")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported by the backend", http.StatusMethodNotAllowed)
		return
	}
	query := `SELECT text, timestamp 
	FROM messages
	ORDER BY timestamp DESC;`
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch messages in backend", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.Text, &msg.Timestamp); err != nil {
			http.Error(w, "Error scanning message in backend", http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func connectToMySQL() {
	var err error
	for retries := 0; retries < 5; retries++ {
		db, err = sql.Open("mysql", "root:password@tcp(mysql:3306)/messages_db")
		if err != nil {
			log.Printf("Error connecting to the database: %v. Retrying...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if err := db.Ping(); err != nil {
			log.Printf("Failed to ping database: %v. Retrying...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("Connected to the MySQL database")
		go func() {
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			<-sigs
			log.Println("Shutting down gracefully...")
			db.Close() // Close DB connection before exiting
		}()
		break
	}
	if err != nil {
		log.Fatalf("Failed to connect to the database after multiple attempts: %v", err)
	}
}

func main() {
	// Enable CORS with a custom policy
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	})
	connectToMySQL()
	// Ensure the table exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id INT AUTO_INCREMENT PRIMARY KEY,
		text VARCHAR(255),
		timestamp DATETIME
	);`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/send-message", handleSendMessage)
	http.HandleFunc("/get-message", handleGetMessage)
	handler := c.Handler(http.DefaultServeMux)
	log.Fatal(http.ListenAndServe(":8080", handler))
	log.Println("Server started at :8080, connected to MySQL")
}
