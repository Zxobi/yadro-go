CREATE TABLE IF NOT EXISTS keywords(
    word    TEXT,
    num     INTEGER,
    FOREIGN KEY (num) REFERENCES comics(num)
);