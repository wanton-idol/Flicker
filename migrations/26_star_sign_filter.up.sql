CREATE TABLE IF NOT EXISTS star_sign_filter (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    star_sign VARCHAR(50) NOT NULL
);

INSERT INTO star_sign_filter (ID, star_sign)
VALUES (1, 'Aries'),
       (2, 'Taurus'),
       (3, 'Gemini'),
       (4, 'Cancer'),
       (5, 'Leo'),
       (6, 'Virgo'),
       (7, 'Libra'),
       (8, 'Scorpio'),
       (9, 'Sagittarius'),
       (10, 'Capricorn'),
       (11, 'Aquarius'),
       (12, 'Pisces');