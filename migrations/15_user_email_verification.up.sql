CREATE TABLE IF NOT EXISTS email_verification (
    ID INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    verification_code VARCHAR(255) UNIQUE NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (ID),
    FOREIGN KEY (user_id) REFERENCES users(ID)
);