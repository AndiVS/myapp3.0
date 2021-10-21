CREATE TABLE IF NOT EXISTS cats(
    _id UUID PRIMARY KEY,
    name VARCHAR(64),
    type VARCHAR(64)
);

CREATE TABLE IF NOT EXISTS users(
    username VARCHAR(64) PRIMARY KEY,
    password VARCHAR(64),
    is_admin BIT
);