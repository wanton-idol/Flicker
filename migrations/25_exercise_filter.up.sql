CREATE TABLE IF NOT EXISTS exercise_filter (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    do_they_exercise VARCHAR(50) NOT NULL
);

INSERT INTO exercise_filter (ID, do_they_exercise)
VALUES (1, 'Active'),
       (2, 'Sometimes'),
       (3, 'Almost Never');