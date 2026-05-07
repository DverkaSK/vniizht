package queries

const (
	UserGetByUsername = `
		SELECT id, username, email, password_hash, role, is_active, reputation, created_at, updated_at
		FROM users
		WHERE username = $1`

	UserGetByID = `
		SELECT id, username, email, password_hash, role, is_active, reputation, created_at, updated_at
		FROM users
		WHERE id = $1`

	UserCreate = `
		INSERT INTO users (username, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	UserUpdateRole = `
		UPDATE users
		SET role = $1, updated_at = NOW()
		WHERE id = $2`

	UserList = `
		SELECT id, username, email, password_hash, role, is_active, reputation, created_at, updated_at
		FROM users
		ORDER BY id`

	UserSetActive = `
		UPDATE users
		SET is_active = $1, updated_at = NOW()
		WHERE id = $2`

	// Пересчитывает репутацию автора ответа по его answer_id (сумма всех голосов на все его ответы)
	UserReputationRecountByAnswer = `
		UPDATE users
		SET reputation = (
		    SELECT COALESCE(SUM(CASE v.value WHEN 'UP' THEN 1 ELSE -1 END), 0)
		    FROM votes v
		    JOIN answers a ON a.id = v.answer_id
		    WHERE a.author_id = (SELECT author_id FROM answers WHERE id = $1)
		)
		WHERE id = (SELECT author_id FROM answers WHERE id = $1)`
)
