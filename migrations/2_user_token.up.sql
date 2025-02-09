CREATE TABLE IF NOT EXISTS user_token (
ID INT NOT NULL AUTO_INCREMENT,
user_id INT NOT NULL,
token varchar(300) NOT NULL,
is_active bool default true,
expires_at TIMESTAMP,
deleted_at TIMESTAMP,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY (ID),
FOREIGN KEY (user_id) REFERENCES users(ID),
INDEX user_token_user_id_index (user_id,token,is_active),
INDEX token_is_active (token,is_active)
);



