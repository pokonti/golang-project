package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/zhanserikAmangeldi/chat-service/internal/adapters/handler"
	"github.com/zhanserikAmangeldi/chat-service/internal/adapters/repository"
	"github.com/zhanserikAmangeldi/chat-service/internal/adapters/websocket"
	"github.com/zhanserikAmangeldi/chat-service/internal/core/service"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// testing purpose now
	connStr := "user=postgres password=secret dbname=postgres sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	log.Println("Connected to Database")

	wsManager := websocket.NewClientManager()
	repo := repository.NewPostgresRepository(db)

	// Pass wsManager to the service
	chatService := service.NewChatService(repo, wsManager)

	wsHandler := handler.NewWSHandler(wsManager)
	http.HandleFunc("/ws", wsHandler.HandleConnection)

	type SendMessageRequest struct {
		SenderID    string `json:"sender_id"`
		RecipientID string `json:"recipient_id"`
		Content     string `json:"content"`
	}

	// HTTP Endpoint to simulate sending a message
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// JSON Request now expects INTs, not Strings
		type SendMessageRequest struct {
			SenderID    int64  `json:"sender_id"`
			RecipientID int64  `json:"recipient_id"`
			Content     string `json:"content"`
		}

		var req SendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		msg, err := chatService.SendMessage(r.Context(), req.SenderID, req.RecipientID, req.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)
	})

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("Server failed:", err)
	}
}
