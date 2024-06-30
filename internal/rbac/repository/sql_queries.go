package repository

const (
	fetchRolesList = `
		SELECT id, name, description, parent_role_id
		FROM roles
		ORDER BY COALESCE(NULLIF($1, ''), name) OFFSET $2 LIMIT $3
	`
	getTotal = `SELECT COUNT(id) FROM roles`
)
