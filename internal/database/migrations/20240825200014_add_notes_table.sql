-- +goose Up
-- +goose StatementBegin
CREATE TABLE notes (
    id UUID DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    content TEXT NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE notes;
-- +goose StatementEnd
