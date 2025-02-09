CREATE TABLE IF NOT EXISTS looking_for_filter(
    ID INT AUTO_INCREMENT PRIMARY KEY,
    looking_for VARCHAR(50) NOT NULL
);

INSERT INTO looking_for_filter(looking_for)
VALUES ('Relationship'),
       ('Something Casual'),
       ('Don''t Know Yet'),
       ('Marriage');