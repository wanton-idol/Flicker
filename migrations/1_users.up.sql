CREATE TABLE IF NOT EXISTS users (
    ID INT NOT NULL AUTO_INCREMENT,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    email text NOT NULL,
    code VARCHAR(5) NOT NULL,
    mobile VARCHAR(10) NOT NULL,
    password VARCHAR(100),
    is_active boolean default true,
    sign_up_type VARCHAR(50) NOT NULL,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (ID),
    KEY idx_mobile (mobile)
);