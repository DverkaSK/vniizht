package queries

const (
	QuestionCreate = `
		INSERT INTO questions (author_id, category_id, title, body)
		VALUES ($1, $2, $3, $4)
		RETURNING id, status, view_count, created_at, updated_at`

	QuestionGetByID = `
		SELECT q.id, q.author_id, u.username,
		       q.category_id, q.specialist_id, q.title, q.body,
		       q.status, q.duplicate_of, q.view_count, q.created_at, q.updated_at,
		       (SELECT COUNT(*) FROM answers a WHERE a.question_id = q.id),
		       (SELECT COALESCE(bool_or(a.is_verified), false) FROM answers a WHERE a.question_id = q.id)
		FROM questions q
		JOIN users u ON u.id = q.author_id
		WHERE q.id = $1`

	QuestionIncrementView = `
		UPDATE questions SET view_count = view_count + 1 WHERE id = $1`

	QuestionList = `
		SELECT q.id, q.author_id, u.username,
		       q.category_id, q.specialist_id, q.title, q.body,
		       q.status, q.duplicate_of, q.view_count, q.created_at, q.updated_at,
		       (SELECT COUNT(*) FROM answers a WHERE a.question_id = q.id),
		       (SELECT COALESCE(bool_or(a.is_verified), false) FROM answers a WHERE a.question_id = q.id)
		FROM questions q
		JOIN users u ON u.id = q.author_id
		WHERE ($1::question_status IS NULL OR q.status = $1)
		  AND ($2::bigint IS NULL OR q.category_id = $2)
		ORDER BY q.created_at DESC
		LIMIT $3 OFFSET $4`

	QuestionCount = `
		SELECT COUNT(*)
		FROM questions
		WHERE ($1::question_status IS NULL OR status = $1)
		  AND ($2::bigint IS NULL OR category_id = $2)`

	QuestionUpdate = `
		UPDATE questions
		SET title = $1, body = $2, updated_at = NOW()
		WHERE id = $3`

	QuestionClose = `
		UPDATE questions
		SET status = 'CLOSED', updated_at = NOW()
		WHERE id = $1 AND status = 'OPEN'`

	QuestionTagsGet = `
		SELECT t.id, t.name, t.created_at
		FROM tags t
		JOIN question_tags qt ON qt.tag_id = t.id
		WHERE qt.question_id = $1
		ORDER BY t.name`

	QuestionTagInsert = `
		INSERT INTO question_tags (question_id, tag_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING`

	QuestionTagsDelete = `
		DELETE FROM question_tags WHERE question_id = $1`

	QuestionHistoryInsert = `
		INSERT INTO question_history (question_id, editor_id, title, body)
		VALUES ($1, $2, $3, $4)`

	QuestionCountByAuthor = `
		SELECT COUNT(*) FROM questions WHERE author_id = $1`

	QuestionListRecentByAuthor = `
		SELECT q.id, q.author_id, u.username,
		       q.category_id, q.specialist_id, q.title, q.body,
		       q.status, q.duplicate_of, q.view_count, q.created_at, q.updated_at,
		       (SELECT COUNT(*) FROM answers a WHERE a.question_id = q.id),
		       (SELECT COALESCE(bool_or(a.is_verified), false) FROM answers a WHERE a.question_id = q.id)
		FROM questions q
		JOIN users u ON u.id = q.author_id
		WHERE q.author_id = $1
		ORDER BY q.created_at DESC
		LIMIT $2`

	QuestionListForExport = `
		SELECT q.id, q.title, q.body, a.body
		FROM questions q
		LEFT JOIN answers a ON a.question_id = q.id AND a.is_verified = true
		ORDER BY q.id`

	QuestionListTitles = `SELECT title FROM questions`
)
