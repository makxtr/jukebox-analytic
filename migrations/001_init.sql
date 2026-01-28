CREATE TABLE IF NOT EXISTS tracks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    artist TEXT NOT NULL,
    price REAL NOT NULL
);

CREATE TABLE IF NOT EXISTS playback_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    track_id INTEGER NOT NULL,
    played_at DATETIME NOT NULL,
    amount_paid REAL NOT NULL,
    FOREIGN KEY(track_id) REFERENCES tracks(id)
);

INSERT OR IGNORE INTO tracks (id, title, artist, price) VALUES
(1, 'Dirty Diana', 'Michael Jackson', 1.25),
(2, 'Comfortably Numb', 'Pink Floyd', 1.50),
(3, 'Space Oddity', 'David Bowie', 1.00);
