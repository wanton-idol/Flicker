CREATE TABLE IF NOT EXISTS education_filter (
    ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    education VARCHAR(50) NOT NULL
);

INSERT INTO education_filter (education)
VALUES ('High school'),
       ('Vocational school'),
       ('In college'),
       ('Undergraduate degree'),
       ('In grade school'),
       ('Graduate degree');