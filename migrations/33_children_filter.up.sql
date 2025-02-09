CREATE TABLE IF NOT EXISTS children_filter(
    ID INT AUTO_INCREMENT PRIMARY KEY,
    have_or_want_children VARCHAR(50) NOT NULL
);

INSERT INTO children_filter(have_or_want_children)
VALUES('Want someday'),
      ('Don''t Want'),
      ('Have and want more'),
      ('Have and don''t want more'),
      ('Not sure yet'),
      ('Have Kids'),
      ('Open to kids');