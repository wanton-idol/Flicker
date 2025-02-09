CREATE TABLE IF NOT EXISTS user_search_profile (
    ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    user_profile_id INT NOT NULL,
    gender varchar(100) DEFAULT NULL,
    min_age INT DEFAULT 18,
    max_age INT DEFAULT 50,
    distance INT DEFAULT NULL,
    language JSON Default NULL,
    deleted_at TIME DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_profile_id) REFERENCES user_profile(ID),
    FOREIGN KEY (user_id) REFERENCES users(ID)
);