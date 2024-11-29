

-- Create the `auth` table
CREATE TABLE auth (
    username VARCHAR(100) PRIMARY KEY,
    password VARCHAR(25) NOT NULL
);

-- Create the `tasks` table
CREATE TABLE tasks (
    id INT NOT NULL,
    description TEXT,
    username VARCHAR(100) NOT NULL,
    FOREIGN KEY (username) REFERENCES auth(username)
);

-- Create the `session` table
CREATE TABLE session (
    username VARCHAR(100) NOT NULL,
    session_id VARCHAR(200) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (username) REFERENCES auth(username)
);


