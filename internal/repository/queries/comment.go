package queries

const (
	CommentCreate = `
		INSERT INTO comments (answer_id, author_id, body)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	CommentGetByID = `
		SELECT id, answer_id, author_id, body, created_at, updated_at
		FROM comments
		WHERE id = $1`

	CommentListByAnswer = `
		SELECT c.id, c.answer_id, c.author_id, u.username, c.body, c.created_at, c.updated_at
		FROM comments c
		JOIN users u ON u.id = c.author_id
		WHERE c.answer_id = $1
		ORDER BY c.created_at ASC`

	CommentUpdate = `
		UPDATE comments
		SET body = $1, updated_at = NOW()
		WHERE id = $2`

	CommentDelete = `
		DELETE FROM comments
		WHERE id = $1`
)
