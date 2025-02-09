CREATE TABLE IF NOT EXISTS user_nudge_profile(
    ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_profile_id INT NOT NULL,
    question text NOT NULL,
    answer text NOT NULL,
    `order` INT NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_profile_id) REFERENCES user_profile(ID)
);