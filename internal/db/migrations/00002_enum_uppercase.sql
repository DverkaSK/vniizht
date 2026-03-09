-- +goose Up
-- +goose StatementBegin

ALTER TYPE user_role RENAME VALUE 'guest'      TO 'GUEST';
ALTER TYPE user_role RENAME VALUE 'user'       TO 'USER';
ALTER TYPE user_role RENAME VALUE 'specialist' TO 'SPECIALIST';
ALTER TYPE user_role RENAME VALUE 'admin'      TO 'ADMIN';

ALTER TYPE question_status RENAME VALUE 'open'      TO 'OPEN';
ALTER TYPE question_status RENAME VALUE 'closed'    TO 'CLOSED';
ALTER TYPE question_status RENAME VALUE 'duplicate' TO 'DUPLICATE';

ALTER TYPE vote_value RENAME VALUE 'up'   TO 'UP';
ALTER TYPE vote_value RENAME VALUE 'down' TO 'DOWN';

ALTER TYPE attachment_target RENAME VALUE 'question' TO 'QUESTION';
ALTER TYPE attachment_target RENAME VALUE 'answer'   TO 'ANSWER';
ALTER TYPE attachment_target RENAME VALUE 'comment'  TO 'COMMENT';

ALTER TYPE notification_type RENAME VALUE 'new_answer'         TO 'NEW_ANSWER';
ALTER TYPE notification_type RENAME VALUE 'answer_verified'    TO 'ANSWER_VERIFIED';
ALTER TYPE notification_type RENAME VALUE 'question_assigned'  TO 'QUESTION_ASSIGNED';
ALTER TYPE notification_type RENAME VALUE 'question_closed'    TO 'QUESTION_CLOSED';
ALTER TYPE notification_type RENAME VALUE 'new_comment'        TO 'NEW_COMMENT';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TYPE notification_type RENAME VALUE 'NEW_COMMENT'        TO 'new_comment';
ALTER TYPE notification_type RENAME VALUE 'QUESTION_CLOSED'    TO 'question_closed';
ALTER TYPE notification_type RENAME VALUE 'QUESTION_ASSIGNED'  TO 'question_assigned';
ALTER TYPE notification_type RENAME VALUE 'ANSWER_VERIFIED'    TO 'answer_verified';
ALTER TYPE notification_type RENAME VALUE 'NEW_ANSWER'         TO 'new_answer';

ALTER TYPE attachment_target RENAME VALUE 'COMMENT'  TO 'comment';
ALTER TYPE attachment_target RENAME VALUE 'ANSWER'   TO 'answer';
ALTER TYPE attachment_target RENAME VALUE 'QUESTION' TO 'question';

ALTER TYPE vote_value RENAME VALUE 'DOWN' TO 'down';
ALTER TYPE vote_value RENAME VALUE 'UP'   TO 'up';

ALTER TYPE question_status RENAME VALUE 'DUPLICATE' TO 'duplicate';
ALTER TYPE question_status RENAME VALUE 'CLOSED'    TO 'closed';
ALTER TYPE question_status RENAME VALUE 'OPEN'      TO 'open';

ALTER TYPE user_role RENAME VALUE 'ADMIN'      TO 'admin';
ALTER TYPE user_role RENAME VALUE 'SPECIALIST' TO 'specialist';
ALTER TYPE user_role RENAME VALUE 'USER'       TO 'user';
ALTER TYPE user_role RENAME VALUE 'GUEST'      TO 'guest';
-- +goose StatementEnd
