CREATE TABLE messages (
                          id BIGSERIAL PRIMARY KEY,
                          conversation_id BIGINT REFERENCES conversations(id) ON DELETE CASCADE,
                          sender_id BIGINT NOT NULL,
                          content TEXT NOT NULL,
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_messages_conv_id ON messages(conversation_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
