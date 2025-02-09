CREATE TABLE IF NOT EXISTS user_nudges(
    ID INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    question text,
    answer text,
    `order` INT,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (ID),
    FOREIGN KEY (user_id) REFERENCES users(ID)
);