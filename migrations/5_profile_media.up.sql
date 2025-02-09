CREATE TABLE IF NOT EXISTS profile_media (
    ID INT NOT NULL AUTO_INCREMENT,
    user_id INT,
    user_profile_id INT,
    url VARCHAR(255) NOT NULL,
    order_id int,
    deleted_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(ID),
    FOREIGN KEY (user_profile_id) REFERENCES user_profile(ID)
);