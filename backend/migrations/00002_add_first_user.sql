-- +goose Up
-- +goose StatementBegin
INSERT INTO "User" (ID, Username, Email, PasswordHash) VALUES (1, 'admin', 'admin@localhost', 'password');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "User" WHERE ID = 1;
-- +goose StatementEnd
