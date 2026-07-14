CREATE TABLE users (
     id INTEGER PRIMARY KEY,
     username TEXT NOT NULL,
     hash_password TEXT NOT NULL,
     email TEXT NOT NULL

);

CREATE TABLE links (
     id INTEGER PRIMARY KEY,
     short_id TEXT UNIQUE,
     orig_url TEXT NOT NULL,
     expiry DATETIME NOT NULL,
     user_id INTEGER NOT NULL,
     FOREIGN KEY (user_id) REFERENCES user(id)
);

CREATE TABLE linkdata (
     id INTEGER PRIMARY KEY,
     access_time DATETIM,E DEFAULT CURRENT_TIMESTAMP,
     country TEXT,
     link_id INTEGER NOT NULL,
     FOREIGN KEY (link_id) REFERENCES link(id)
);