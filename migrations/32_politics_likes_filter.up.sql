CREATE TABLE IF NOT EXISTS politics_likes_filter(
    ID INT AUTO_INCREMENT PRIMARY KEY,
    politics_likes VARCHAR(50) NOT NULL
);

INSERT INTO politics_likes_filter(politics_likes)
VALUES('Apolitical'),
      ('Moderate'),
      ('Left'),
      ('Right'),
      ('Communist'),
      ('Socialist');