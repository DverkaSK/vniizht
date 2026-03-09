package queries

const (
	TagList = `
		SELECT id, name, created_at
		FROM tags
		ORDER BY name`

	TagGetByID = `
		SELECT id, name, created_at
		FROM tags
		WHERE id = $1`

	TagCreate = `
		INSERT INTO tags (name)
		VALUES ($1)
		RETURNING id, created_at`

	TagUpdate = `
		UPDATE tags
		SET name = $1
		WHERE id = $2`

	TagDelete = `
		DELETE FROM tags
		WHERE id = $1`
)
