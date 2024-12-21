-- +goose Up
-- +goose StatementBegin
INSERT INTO Users (ID, Username, Email, PasswordHash) VALUES (1, 'admin', 'admin@localhost', 'password');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM Users WHERE ID = 1;
-- +goose StatementEnd
