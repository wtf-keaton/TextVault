-- +goose Up
CREATE TABLE "User" (
    ID BIGINT PRIMARY KEY,
    Username VARCHAR(50) NOT NULL,
    Email VARCHAR(100) UNIQUE NOT NULL,
    PasswordHash VARCHAR(255) NOT NULL,

    UNIQUE (Username, Email)
);

CREATE INDEX idx_user_username ON "User" (Username);
CREATE INDEX idx_user_email ON "User" (Email);

CREATE TABLE "Paste" (
    ID BIGINT PRIMARY KEY,
    Title VARCHAR(255) NOT NULL,
    Hash VARCHAR(64) UNIQUE NOT NULL,
    AuthorID BIGINT NOT NULL,
    Content TEXT NOT NULL,
    FOREIGN KEY (AuthorID) REFERENCES "User"(ID)
);

CREATE INDEX idx_paste_title ON "Paste" (Title);
CREATE INDEX idx_paste_author_id ON "Paste" (AuthorID);

-- +goose Down
DROP INDEX IF EXISTS idx_paste_author_id;
DROP INDEX IF EXISTS idx_paste_title;
DROP TABLE IF EXISTS "Paste";

DROP INDEX IF EXISTS idx_user_email;
DROP INDEX IF EXISTS idx_user_username;
DROP TABLE IF EXISTS "User";