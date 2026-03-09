-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE user_role AS ENUM ('guest', 'user', 'specialist', 'admin');

CREATE TABLE users (
    id          BIGSERIAL PRIMARY KEY,
    username    TEXT        NOT NULL UNIQUE,
    email       TEXT        NOT NULL UNIQUE,
    password_hash TEXT      NOT NULL,
    role        user_role   NOT NULL DEFAULT 'user',
    is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
    reputation  INT         NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sessions (
    id          TEXT        PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ NOT NULL,
    last_active TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user_id   ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

CREATE TABLE categories (
    id           BIGSERIAL PRIMARY KEY,
    name         TEXT    NOT NULL UNIQUE,
    description  TEXT    NOT NULL DEFAULT '',
    specialist_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tags (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT        NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TYPE question_status AS ENUM ('open', 'closed', 'duplicate');

CREATE TABLE questions (
    id           BIGSERIAL PRIMARY KEY,
    author_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id  BIGINT      REFERENCES categories(id) ON DELETE SET NULL,
    specialist_id BIGINT     REFERENCES users(id) ON DELETE SET NULL,
    title        TEXT        NOT NULL,
    body         TEXT        NOT NULL,
    status       question_status NOT NULL DEFAULT 'open',
    duplicate_of BIGINT      REFERENCES questions(id) ON DELETE SET NULL,
    view_count   INT         NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    search_vector TSVECTOR
);

CREATE INDEX idx_questions_author_id    ON questions(author_id);
CREATE INDEX idx_questions_category_id  ON questions(category_id);
CREATE INDEX idx_questions_status       ON questions(status);
CREATE INDEX idx_questions_created_at   ON questions(created_at DESC);
CREATE INDEX idx_questions_search       ON questions USING GIN(search_vector);

CREATE OR REPLACE FUNCTION questions_search_vector_update() RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('russian', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('russian', COALESCE(NEW.body,  '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_questions_search_vector
    BEFORE INSERT OR UPDATE OF title, body ON questions
    FOR EACH ROW EXECUTE FUNCTION questions_search_vector_update();

CREATE TABLE question_tags (
    question_id BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    tag_id      BIGINT NOT NULL REFERENCES tags(id)      ON DELETE CASCADE,
    PRIMARY KEY (question_id, tag_id)
);

CREATE INDEX idx_question_tags_tag_id ON question_tags(tag_id);

CREATE TABLE question_history (
    id          BIGSERIAL PRIMARY KEY,
    question_id BIGINT      NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    editor_id   BIGINT      NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
    title       TEXT        NOT NULL,
    body        TEXT        NOT NULL,
    edited_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_question_history_question_id ON question_history(question_id);

CREATE TABLE answers (
    id          BIGSERIAL PRIMARY KEY,
    question_id BIGINT      NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    author_id   BIGINT      NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
    body        TEXT        NOT NULL,
    is_verified BOOLEAN     NOT NULL DEFAULT FALSE,
    vote_score  INT         NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    search_vector TSVECTOR
);

CREATE INDEX idx_answers_question_id ON answers(question_id);
CREATE INDEX idx_answers_author_id   ON answers(author_id);
CREATE INDEX idx_answers_search      ON answers USING GIN(search_vector);

CREATE OR REPLACE FUNCTION answers_search_vector_update() RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('russian', COALESCE(NEW.body, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_answers_search_vector
    BEFORE INSERT OR UPDATE OF body ON answers
    FOR EACH ROW EXECUTE FUNCTION answers_search_vector_update();

CREATE TABLE answer_history (
    id        BIGSERIAL PRIMARY KEY,
    answer_id BIGINT      NOT NULL REFERENCES answers(id) ON DELETE CASCADE,
    editor_id BIGINT      NOT NULL REFERENCES users(id)   ON DELETE CASCADE,
    body      TEXT        NOT NULL,
    edited_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_answer_history_answer_id ON answer_history(answer_id);

CREATE TABLE comments (
    id        BIGSERIAL PRIMARY KEY,
    answer_id BIGINT      NOT NULL REFERENCES answers(id) ON DELETE CASCADE,
    author_id BIGINT      NOT NULL REFERENCES users(id)   ON DELETE CASCADE,
    body      TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comments_answer_id ON comments(answer_id);

CREATE TYPE vote_value AS ENUM ('up', 'down');

CREATE TABLE votes (
    id        BIGSERIAL PRIMARY KEY,
    answer_id BIGINT      NOT NULL REFERENCES answers(id) ON DELETE CASCADE,
    user_id   BIGINT      NOT NULL REFERENCES users(id)   ON DELETE CASCADE,
    value     vote_value  NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (answer_id, user_id)
);

CREATE INDEX idx_votes_answer_id ON votes(answer_id);

CREATE TYPE attachment_target AS ENUM ('question', 'answer', 'comment');

CREATE TABLE attachments (
    id          BIGSERIAL PRIMARY KEY,
    uploader_id BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_type attachment_target NOT NULL,
    target_id   BIGINT      NOT NULL,
    filename    TEXT        NOT NULL,
    object_key  TEXT        NOT NULL,
    mime_type   TEXT        NOT NULL DEFAULT 'application/octet-stream',
    size_bytes  BIGINT      NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_attachments_target ON attachments(target_type, target_id);

CREATE TYPE notification_type AS ENUM (
    'new_answer',
    'answer_verified',
    'question_assigned',
    'question_closed',
    'new_comment'
);

CREATE TABLE notifications (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type        notification_type NOT NULL,
    payload     JSONB       NOT NULL DEFAULT '{}',
    is_read     BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id  ON notifications(user_id);
CREATE INDEX idx_notifications_is_read  ON notifications(user_id, is_read);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS attachments;
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS answer_history;
DROP TABLE IF EXISTS answers;
DROP TABLE IF EXISTS question_history;
DROP TABLE IF EXISTS question_tags;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS notification_type;
DROP TYPE IF EXISTS attachment_target;
DROP TYPE IF EXISTS vote_value;
DROP TYPE IF EXISTS question_status;
DROP TYPE IF EXISTS user_role;
-- +goose StatementEnd