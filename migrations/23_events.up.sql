CREATE TABLE IF NOT EXISTS events (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    event_time TIMESTAMP NOT NULL,
    type VARCHAR(30) NOT NULL NOT NULL ,
    description TEXT,
    attendees text,
    expires_at TIMESTAMP,
    address1 text,
    address2 text,
    city varchar(100),
    state varchar(50),
    country varchar(50),
    pincode int ,
    latitude float,
    longitude float,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_profile(ID)
);