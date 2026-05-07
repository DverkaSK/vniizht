package queries

const (
	AnswerCreate = `
		INSERT INTO answers (question_id, author_id, body)
		VALUES ($1, $2, $3)
		RETURNING id, is_verified, vote_score, created_at, updated_at`

	AnswerGetByID = `
		SELECT id, question_id, author_id, body, is_verified, vote_score, created_at, updated_at
		FROM answers
		WHERE id = $1`

	AnswerListByQuestion = `
		SELECT a.id, a.question_id, a.author_id, u.username, a.body, a.is_verified, a.vote_score, v.value, a.created_at, a.updated_at
		FROM answers a
		JOIN users u ON u.id = a.author_id
		LEFT JOIN votes v ON v.answer_id = a.id AND v.user_id = $2
		WHERE a.question_id = $1
		ORDER BY a.is_verified DESC, a.vote_score DESC, a.created_at ASC`

	AnswerUpdate = `
		UPDATE answers
		SET body = $1, updated_at = NOW()
		WHERE id = $2`

	AnswerDelete = `
		DELETE FROM answers
		WHERE id = $1`

	AnswerHistoryInsert = `
		INSERT INTO answer_history (answer_id, editor_id, body)
		VALUES ($1, $2, $3)`

	AnswerCountByAuthor = `
		SELECT COUNT(*) FROM answers WHERE author_id = $1`

	AnswerListRecentByAuthor = `
		SELECT a.id, a.question_id, a.author_id, u.username, a.body, a.is_verified, a.vote_score, a.created_at, a.updated_at
		FROM answers a
		JOIN users u ON u.id = a.author_id
		WHERE a.author_id = $1
		ORDER BY a.created_at DESC
		LIMIT $2`

	AnswerSetVerified = `
		UPDATE answers
		SET is_verified = CASE WHEN id = $2 THEN TRUE ELSE FALSE END,
		    updated_at = NOW()
		WHERE question_id = $1
		  AND (id = $2 OR is_verified = TRUE)`

	AnswerClearVerified = `
		UPDATE answers
		SET is_verified = FALSE, updated_at = NOW()
		WHERE id = $1`

	AnswerVoteInsert = `
		INSERT INTO votes (answer_id, user_id, value)
		VALUES ($1, $2, $3::vote_value)`

	AnswerVoteUpdate = `
		UPDATE votes
		SET value = $3::vote_value
		WHERE answer_id = $1 AND user_id = $2`

	AnswerVoteDelete = `
		DELETE FROM votes
		WHERE answer_id = $1 AND user_id = $2`

	AnswerVoteGet = `
		SELECT value
		FROM votes
		WHERE answer_id = $1 AND user_id = $2`

	AnswerHistoryGet = `
		SELECT ah.id, ah.answer_id, ah.editor_id, u.username, ah.body, ah.edited_at
		FROM answer_history ah
		JOIN users u ON u.id = ah.editor_id
		WHERE ah.answer_id = $1
		ORDER BY ah.edited_at DESC`

	AnswerVoteRecount = `
		UPDATE answers
		SET vote_score = (
			SELECT COALESCE(SUM(CASE value WHEN 'UP' THEN 1 ELSE -1 END), 0)
			FROM votes
			WHERE answer_id = $1
		)
		WHERE id = $1`
)
