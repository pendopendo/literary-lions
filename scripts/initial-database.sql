-- siia paneme koodid vahepeal
-- (26.09.24)

CREATE TABLE category (
	id INTEGER PRIMARY KEY, --et saaks olla ainult üks unikaalne, see see P KEY
   	name TEXT NOT NULL --ei tohi olla tühi (not null)
);

CREATE TABLE comment (
    id INTEGER PRIMARY KEY, 
    user INTEGER NOT NULL, 
    post INTEGER NOT NULL, 
    text TEXT NOT NULL, 
    created TEXT NOT NULL,
    FOREIGN KEY(user) REFERENCES user(id), --selleks et  kontrollida, et see user on päriselt olemas
    FOREIGN KEY(post) REFERENCES post(id)
);

CREATE TABLE post (
    id INTEGER PRIMARY KEY, 
    user INTEGER NOT NULL, 
    title TEXT NOT NULL, 
    text TEXT NOT NULL, 
    category INTEGER NOT NULL, 
    created TEXT NOT NULL,
    FOREIGN KEY(user) REFERENCES user(id),
    FOREIGN KEY(category) REFERENCES category(id)
);


CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- Primary key with auto-increment
    name TEXT NOT NULL, -- Using TEXT instead of VARCHAR
    email TEXT NOT NULL, -- Using TEXT instead of VARCHAR
    hashed_password CHAR(60) NOT NULL, -- CHAR is fine as it is
    created DATETIME NOT NULL
);

-- Adding a UNIQUE constraint to the email column. ja username ka!
CREATE UNIQUE INDEX users_uc_email ON users (email);
CREATE UNIQUE INDEX users_uc_name ON users (name);


CREATE TABLE reactions (
    id INTEGER PRIMARY KEY, 
    reaction STRING NOT NULL, 
    comment INTEGER, 
    post INTEGER, 
    user INTEGER NOT NULL,
    FOREIGN KEY(comment) REFERENCES comment(id),
    FOREIGN KEY(post) REFERENCES post(id),
    FOREIGN KEY(user) REFERENCES user(id)
);

CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,  -- The token will be treated as text with a fixed length of 43
    data BLOB NOT NULL,          -- Binary data storage
    expiry DATETIME NOT NULL     -- Use DATETIME instead of TIMESTAMP(6)
);

-- Create an index on the expiry column
CREATE INDEX sessions_expiry_idx ON sessions (expiry);
