-- +goose Up
-- +goose StatementBegin
CREATE TABLE comment_history (
    id         BIGSERIAL PRIMARY KEY,
    comment_id BIGINT      NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    editor_id  BIGINT      NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
    body       TEXT        NOT NULL,
    edited_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comment_history_comment_id ON comment_history(comment_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS comment_history;
-- +goose StatementEnd
