CREATE TABLE IF NOT EXISTS category (
    id INTEGER PRIMARY KEY, 
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS comment (
    id INTEGER PRIMARY KEY, 
    user INTEGER NOT NULL, 
    post INTEGER NOT NULL, 
    text TEXT NOT NULL, 
    created TEXT NOT NULL,
    FOREIGN KEY(user) REFERENCES users(id),
    FOREIGN KEY(post) REFERENCES post(id)
);

CREATE TABLE IF NOT EXISTS post (
    id INTEGER PRIMARY KEY, 
    user INTEGER NOT NULL, 
    title TEXT NOT NULL, 
    text TEXT NOT NULL, 
    category INTEGER NOT NULL, 
    created TEXT NOT NULL,
    FOREIGN KEY(user) REFERENCES users(id),
    FOREIGN KEY(category) REFERENCES category(id)
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    name TEXT NOT NULL, 
    email TEXT NOT NULL, 
    hashed_password CHAR(60) NOT NULL, 
    created DATETIME NOT NULL
);

-- Adding a UNIQUE constraint to the email and username columns
CREATE UNIQUE INDEX IF NOT EXISTS users_uc_email ON users (email);
CREATE UNIQUE INDEX IF NOT EXISTS users_uc_name ON users (name);

CREATE TABLE IF NOT EXISTS reactions (
    id INTEGER PRIMARY KEY, 
    reaction STRING NOT NULL, 
    comment INTEGER, 
    post INTEGER, 
    user INTEGER NOT NULL,
    FOREIGN KEY(comment) REFERENCES comment(id),
    FOREIGN KEY(post) REFERENCES post(id),
    FOREIGN KEY(user) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS sessions (
    token CHAR(43) PRIMARY KEY, 
    data BLOB NOT NULL, 
    expiry DATETIME NOT NULL
);

-- Create an index on the expiry column
CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);

-- Insert predefined categories into the category table
INSERT INTO category (id, name) VALUES 
(1, 'Fiction'),
(2, 'Non-Fiction'),
(3, 'Science Fiction'),
(4, 'Fantasy')
ON CONFLICT(id) DO NOTHING;