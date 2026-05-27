package queries

const (
	QuestionCreate = `
		INSERT INTO questions (author_id, category_id, specialist_id, title, body)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, status, view_count, created_at, updated_at`

	QuestionGetByID = `
		SELECT q.id, q.author_id, u.username,
		       q.category_id, q.specialist_id, s.username,
		       q.title, q.body, q.status, q.duplicate_of, q.view_count, q.created_at, q.updated_at,
		       (SELECT COUNT(*) FROM answers a WHERE a.question_id = q.id),
		       (SELECT COALESCE(bool_or(a.is_verified), false) FROM answers a WHERE a.question_id = q.id)
		FROM questions q
		JOIN users u ON u.id = q.author_id
		LEFT JOIN users s ON s.id = q.specialist_id
		WHERE q.id = $1`

	QuestionIncrementView = `
		UPDATE questions SET view_count = view_count + 1 WHERE id = $1`

	QuestionList = `
		SELECT q.id, q.author_id, u.username,
		       q.category_id, q.specialist_id, s.username,
		       q.title, q.body, q.status, q.duplicate_of, q.view_count, q.created_at, q.updated_at,
		       (SELECT COUNT(*) FROM answers a WHERE a.question_id = q.id),
		       (SELECT COALESCE(bool_or(a.is_verified), false) FROM answers a WHERE a.question_id = q.id)
		FROM questions q
		JOIN users u ON u.id = q.author_id
		LEFT JOIN users s ON s.id = q.specialist_id
		WHERE ($1::question_status IS NULL OR q.status = $1)
		  AND ($2::bigint IS NULL OR q.category_id = $2)
		  AND ($3::bigint IS NULL OR EXISTS (
		      SELECT 1 FROM question_tags qt WHERE qt.question_id = q.id AND qt.tag_id = $3
		  ))
		  AND ($4::bigint IS NULL OR q.specialist_id = $4)
		ORDER BY q.created_at DESC
		LIMIT $5 OFFSET $6`

	QuestionCount = `
		SELECT COUNT(*)
		FROM questions q
		WHERE ($1::question_status IS NULL OR q.status = $1)
		  AND ($2::bigint IS NULL OR q.category_id = $2)
		  AND ($3::bigint IS NULL OR EXISTS (
		      SELECT 1 FROM question_tags qt WHERE qt.question_id = q.id AND qt.tag_id = $3
		  ))
		  AND ($4::bigint IS NULL OR q.specialist_id = $4)`

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
		       q.category_id, q.specialist_id, s.username,
		       q.title, q.body, q.status, q.duplicate_of, q.view_count, q.created_at, q.updated_at,
		       (SELECT COUNT(*) FROM answers a WHERE a.question_id = q.id),
		       (SELECT COALESCE(bool_or(a.is_verified), false) FROM answers a WHERE a.question_id = q.id)
		FROM questions q
		JOIN users u ON u.id = q.author_id
		LEFT JOIN users s ON s.id = q.specialist_id
		WHERE q.author_id = $1
		ORDER BY q.created_at DESC
		LIMIT $2`

	QuestionAssign = `
		UPDATE questions SET specialist_id = $2, updated_at = NOW() WHERE id = $1`

	QuestionListForExport = `
		SELECT q.id, q.title, q.body, a.body
		FROM questions q
		LEFT JOIN answers a ON a.question_id = q.id AND a.is_verified = true
		ORDER BY q.id`

	QuestionListTitles = `SELECT title FROM questions`

	// QuestionListForPeriodExport — полный список вопросов за период для экспорта в CSV.
	// $1 = from (timestamptz, NULL = без ограничения), $2 = to (timestamptz, NULL = без ограничения).
	QuestionListForPeriodExport = `
		SELECT
			q.id,
			q.created_at,
			q.title,
			q.body,
			c.name,
			STRING_AGG(DISTINCT t.name, ', ' ORDER BY t.name),
			u.username,
			q.status,
			s.username,
			(SELECT COUNT(*) FROM answers a WHERE a.question_id = q.id),
			EXISTS(SELECT 1 FROM answers a WHERE a.question_id = q.id AND a.is_verified = true),
			(SELECT a.body FROM answers a WHERE a.question_id = q.id AND a.is_verified = true LIMIT 1)
		FROM questions q
		JOIN  users u ON u.id = q.author_id
		LEFT JOIN categories c  ON c.id  = q.category_id
		LEFT JOIN users      s  ON s.id  = q.specialist_id
		LEFT JOIN question_tags qt ON qt.question_id = q.id
		LEFT JOIN tags t           ON t.id = qt.tag_id
		WHERE ($1::timestamptz IS NULL OR q.created_at >= $1)
		  AND ($2::timestamptz IS NULL OR q.created_at <  $2)
		GROUP BY q.id, c.name, u.username, s.username
		ORDER BY q.created_at`

	QuestionMarkDuplicate = `
		UPDATE questions
		SET status = 'DUPLICATE', duplicate_of = $2, updated_at = NOW()
		WHERE id = $1 AND status = 'OPEN'`

	QuestionDelete = `DELETE FROM questions WHERE id = $1`

	QuestionHistoryGet = `
		SELECT qh.id, qh.question_id, qh.editor_id, u.username, qh.title, qh.body, qh.edited_at
		FROM question_history qh
		JOIN users u ON u.id = qh.editor_id
		WHERE qh.question_id = $1
		ORDER BY qh.edited_at DESC`
)
