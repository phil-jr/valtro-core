package db

const (
	InsertUserQuery = `
		INSERT INTO "USER" (
			user_id,
			company_id,
			first_name,
			last_name,
			password,
			email,
			admin,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) RETURNING user_id
	`

	// Add more queries as needed
	SelectUserQuery = `
		SELECT user_id, company_id, first_name, last_name, email, admin, created_at
		FROM "USER"
		WHERE user_id = $1
	`

	UpdateUserQuery = `
		UPDATE "USER"
		SET first_name = $2, last_name = $3, email = $4, admin = $5
		WHERE user_id = $1
	`

	DeleteUserQuery = `
		DELETE "USER" WHERE user_id = $1
	`

	RetrievePasswordHashQuery = `
		SELECT user_id, password
		FROM "USER"
		WHERE email = $1
	`
)
