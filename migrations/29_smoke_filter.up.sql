CREATE TABLE IF NOT EXISTS smoke_filter (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    do_they_smoke VARCHAR(50) NOT NULL
);

INSERT INTO smoke_filter (do_they_smoke) VALUES ('Socially'), ('Never'), ('Regularly');