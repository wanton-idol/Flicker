BEGIN ;
CREATE TABLE IF NOT EXISTS interests (
    ID INT NOT NULL AUTO_INCREMENT,
    interest_type VARCHAR(50),
    deleted_at TIMESTAMP,
    PRIMARY KEY (ID)
);

INSERT INTO interests (ID, interest_type)
VALUES (1, 'Self care'),
       (2, 'Sports'),
       (3, 'Creativity'),
       (4, 'Going out'),
       (5, 'Film and TV'),
       (6, 'Staying in'),
       (7, 'Reading'),
       (8, 'Music'),
       (9, 'Food and Drink'),
       (10, 'Travelling'),
       (11, 'Pets'),
       (12, 'Values and Traits'),
       (13, 'Pluto values and allyship');

COMMIT ;