CREATE TABLE IF NOT EXISTS religion_filter(
    ID INT AUTO_INCREMENT PRIMARY KEY,
    religion VARCHAR(50) NOT NULL
);

INSERT INTO religion_filter(religion)
VALUES ('Agnostic'),
       ('Atheist'),
       ('Buddhist'),
       ('Catholic'),
       ('Christian'),
       ('Hindu'),
       ('Jain'),
       ('Jewish'),
       ('Mormon'),
       ('Latter-day Saint'),
       ('Muslim'),
       ('Zoroastrian'),
       ('Sikh'),
       ('Spiritual'),
       ('Other');