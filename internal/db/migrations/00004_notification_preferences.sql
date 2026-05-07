-- +goose Up
-- +goose StatementBegin
CREATE TABLE notification_preferences (
    user_id       BIGINT            NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type          notification_type NOT NULL,
    email_enabled BOOLEAN           NOT NULL DEFAULT TRUE,
    PRIMARY KEY (user_id, type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notification_preferences;
-- +goose StatementEnd
