package queries

const (
	CategoryList = `
		SELECT id, name, description, specialist_id, created_at
		FROM categories
		ORDER BY name`

	CategoryGetByID = `
		SELECT id, name, description, specialist_id, created_at
		FROM categories
		WHERE id = $1`

	CategoryCreate = `
		INSERT INTO categories (name, description, specialist_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	CategoryUpdate = `
		UPDATE categories
		SET name = $1, description = $2, specialist_id = $3
		WHERE id = $4`

	CategoryDelete = `
		DELETE FROM categories
		WHERE id = $1`
)
