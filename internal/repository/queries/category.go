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
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at`

	CategoryUpdate = `
		UPDATE categories
		SET name = $1, description = $2
		WHERE id = $3`

	CategoryDelete = `
		DELETE FROM categories
		WHERE id = $1`
)
