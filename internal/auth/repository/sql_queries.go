package repository

const (
	createUserQuery = `INSERT INTO users (username, email, password_hash, created_at, updated_at, login_at)
						VALUES ($1, $2, $3, now(), now(), now()) RETURNING *`

	updateUserQuery = `UPDATE users
						SET username = COALESCE(NULLIF($1, ''), username),
						    email = COALESCE(NULLIF($2, ''), email),
						    updated_at = now()
						WHERE id = $13
						RETURNING *
						`

	deleteUserQuery = `DELETE FROM users WHERE id = $1`

	getUserQuery = `SELECT id, username, email, created_at, updated_at, login_at
					 FROM users
					 WHERE id = $1`
	getUserRoleQuery = `SELECT
							users.id AS "user.id",
							users.username AS "user.username",
							users.email AS "user.email",
							users.password_hash AS "user.password_hash",
							users.created_at AS "user.created_at",
							users.updated_at AS "user.updated_at",
							users.login_at AS "user.login_at",
							r.id AS "role.id",
							r.name AS "role.name",
							r.description AS "role.description",
							r.parent_role_id AS "role.parent_role_id"
						FROM public.users users
						JOIN public.user_roles ar ON ar.user_id = users.id
						JOIN public.roles r ON r.id = ar.role_id
						WHERE users.id = $1`

	getTotalCount = `SELECT COUNT(id) FROM users
						WHERE username ILIKE '%' || $1 || '%'`

	findUsers = `SELECT id, username, email,
	              created_at, updated_at, login_at
				  FROM users
				  WHERE username ILIKE '%' || $1 || '%'
				  ORDER BY username, last_name
				  OFFSET $2 LIMIT $3
				  `

	getTotal = `SELECT COUNT(id) FROM users`

	getUsers = `SELECT id, username, email, created_at, updated_at, login_at
				 FROM users
				 ORDER BY COALESCE(NULLIF($1, ''), username) OFFSET $2 LIMIT $3`

	findUserByEmail = `SELECT id, username, email, password_hash, created_at, updated_at, login_at
				 		FROM users
				 		WHERE email = $1`

	findByUsername = `SELECT
			users.id AS "user.id",
			users.username AS "user.username",
			users.email AS "user.email",
			users.password_hash AS "user.password_hash",
			users.created_at AS "user.created_at",
			users.updated_at AS "user.updated_at",
			users.login_at AS "user.login_at",
			r.id AS "role.id",
			r.name AS "role.name",
			r.description AS "role.description",
			r.parent_role_id AS "role.parent_role_id"
		FROM public.users users
		JOIN public.user_roles ar ON ar.user_id = users.id
		JOIN public.roles r ON r.id = ar.role_id
		WHERE users.username = $1`
)
