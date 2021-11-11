CREATE TABLE IF NOT EXISTS cats(
    _id UUID PRIMARY KEY,
    name VARCHAR(64),
    type VARCHAR(64)
);

CREATE TABLE IF NOT EXISTS users(
    username VARCHAR(64) PRIMARY KEY,
    password VARCHAR(64),
    is_admin BOOLEAN
);

INSERT INTO users (username, password, is_admin) VALUES ('admin', 'a3f03e82c7ba949052ce81a2f8f24a0d287e006103e28d15c38dac04f8eb56ab', true);