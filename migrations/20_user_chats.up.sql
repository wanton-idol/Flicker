CREATE TABLE IF NOT EXISTS user_chats (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    chat_id VARCHAR(30) NOT NULL,
    message TEXT,
    media_url text,
    is_read BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users(ID),
    FOREIGN KEY (receiver_id) REFERENCES users(ID)
);