package queries

const (
	SearchQuestions = `
		SELECT q.id, q.title, q.body, q.status, q.author_id, q.category_id, q.created_at,
		       ts_rank(q.search_vector, plainto_tsquery('russian', $1)) AS rank
		FROM questions q
		WHERE q.search_vector @@ plainto_tsquery('russian', $1)
		  AND ($2::BIGINT IS NULL OR q.category_id = $2)
		  AND ($3::TEXT IS NULL OR q.status::TEXT = $3)
		  AND ($4::BIGINT IS NULL OR EXISTS (
		      SELECT 1 FROM question_tags qt WHERE qt.question_id = q.id AND qt.tag_id = $4
		  ))
		ORDER BY rank DESC
		LIMIT 20`

	SearchAnswers = `
		SELECT a.id, a.question_id, q.title, a.body, q.status, a.author_id, q.category_id, a.created_at,
		       ts_rank(a.search_vector, plainto_tsquery('russian', $1)) AS rank
		FROM answers a
		JOIN questions q ON q.id = a.question_id
		WHERE a.search_vector @@ plainto_tsquery('russian', $1)
		  AND ($2::BIGINT IS NULL OR q.category_id = $2)
		  AND ($3::TEXT IS NULL OR q.status::TEXT = $3)
		  AND ($4::BIGINT IS NULL OR EXISTS (
		      SELECT 1 FROM question_tags qt WHERE qt.question_id = q.id AND qt.tag_id = $4
		  ))
		ORDER BY rank DESC
		LIMIT 20`
)
