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
)
