package queries

const (
	SessionCreate = `
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES ($1, $2, $3)`

	SessionGet = `
		SELECT id, user_id, created_at, expires_at, last_active
		FROM sessions
		WHERE id = $1`

	SessionDelete = `
		DELETE FROM sessions
		WHERE id = $1`

	SessionTouch = `
		UPDATE sessions
		SET last_active = $1
		WHERE id = $2`

	SessionDeleteExpired = `
		DELETE FROM sessions
		WHERE expires_at < NOW()`
)
