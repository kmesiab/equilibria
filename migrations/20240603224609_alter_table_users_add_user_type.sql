-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN user_type_id BIGINT DEFAULT 1 AFTER account_status_id,
    ADD CONSTRAINT fk_user_type_id FOREIGN KEY (user_type_id) REFERENCES user_types (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN user_type_id,
    DROP CONSTRAINT fk_user_type_id;
-- +goose StatementEnd
