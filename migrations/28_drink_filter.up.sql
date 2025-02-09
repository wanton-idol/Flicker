CREATE TABLE IF NOT EXISTS drink_filter (
    ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    do_they_drink VARCHAR(50) NOT NULL
);

INSERT INTO drink_filter (do_they_drink)
VALUES ('Frequently'),
       ('Socially'),
       ('Rarely'),
       ('Never'),
       ('Sober');