CREATE TABLE IF NOT EXISTS user_match (
    ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    match_id INT NOT NULL,
    match_type INT NOT NULL,
    chat_id INT,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
FOREIGN KEY (user_id) REFERENCES users(ID),
FOREIGN KEY (match_id) REFERENCES users(ID),
INDEX user_id_match_id (user_id,match_id),
INDEX match_id_user_id (match_id,user_id)
);





