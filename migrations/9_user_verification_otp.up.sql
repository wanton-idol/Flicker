CREATE TABLE IF NOT EXISTS user_verification_otp (
    ID INT NOT NULL AUTO_INCREMENT,
    phone_number VARCHAR(13) NOT NULL,
    otp int NOT NULL,
    expires_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (ID),
    INDEX idx_phone_number_otp (phone_number,otp)
);