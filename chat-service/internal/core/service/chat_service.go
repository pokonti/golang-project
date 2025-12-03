package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zhanserikAmangeldi/chat-service/internal/adapters/websocket"

	"github.com/zhanserikAmangeldi/chat-service/internal/core/domain"
	"github.com/zhanserikAmangeldi/chat-service/internal/core/ports"
)

type ChatService struct {
	repo      ports.ChatRepository
	wsManager *websocket.ClientManager
}

func NewChatService(repo ports.ChatRepository, wsManager *websocket.ClientManager) *ChatService {
	return &ChatService{
		repo:      repo,
		wsManager: wsManager,
	}
}

// SendMessage handles the logic:
// 1. Check if conversation exists (if not, create it)
// 2. Save message
func (s *ChatService) SendMessage(ctx context.Context, senderID, recipientID int64, content string) (*domain.Message, error) {
	// 1. Check for existing conversation
	conv, err := s.repo.FindOneToOneConversation(ctx, senderID, recipientID)
	if err != nil {
		return nil, err
	}

	// 2. If no conversation, create one
	if conv == nil {
		newConv := &domain.Conversation{
			IsGroup:   false,
			CreatedAt: time.Now(),
		}
		// Repo will fill in the ID
		if err := s.repo.CreateConversation(ctx, newConv); err != nil {
			return nil, err
		}
		conv = newConv

		// Add participants
		p1 := &domain.Participant{ConversationID: conv.ID, UserID: senderID, JoinedAt: time.Now()}
		p2 := &domain.Participant{ConversationID: conv.ID, UserID: recipientID, JoinedAt: time.Now()}

		s.repo.AddParticipant(ctx, p1)
		s.repo.AddParticipant(ctx, p2)
	}

	// 3. Create the Message object
	msg := &domain.Message{
		ConversationID: conv.ID,
		SenderID:       senderID,
		Content:        content,
		CreatedAt:      time.Now(),
	}

	// 4. Save to DB
	if err := s.repo.SaveMessage(ctx, msg); err != nil {
		return nil, err
	}

	if conn, ok := s.wsManager.GetClient(recipientID); ok {
		// 2. Marshal message to JSON
		msgBytes, _ := json.Marshal(msg)

		// 3. Send to Recipient
		if err := conn.WriteMessage(1, msgBytes); err != nil {
			// If sending fails, maybe they just disconnected
			s.wsManager.RemoveClient(recipientID)
		}
	}

	return msg, nil
}

func (s *ChatService) GetHistory(ctx context.Context, conversationID int64) ([]domain.Message, error) {
	return s.repo.GetMessages(ctx, conversationID, 50, 0) // Limit 50 for now
}
