package db

const (
	InsertUserQuery = `
		INSERT INTO "USER" (
			"USER_ID",
			"COMPANY_ID",
			"FIRST_NAME",
			"LAST_NAME",
			"PASSWORD",
			"EMAIL",
			"ADMIN",
			"CREATED_AT"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) RETURNING "USER_ID"
	`

	SelectUserQuery = `
		SELECT "USER_ID", "COMPANY_ID", "FIRST_NAME", "LAST_NAME", "EMAIL", "ADMIN", "CREATED_AT"
		FROM "USER"
		WHERE "USER_ID" = $1
	`

	UpdateUserQuery = `
		UPDATE "USER"
		SET "FIRST_NAME" = $2, "LAST_NAME" = $3, "EMAIL" = $4, "ADMIN" = $5
		WHERE "USER_ID" = $1
	`

	DeleteUserQuery = `
		DELETE FROM "USER" WHERE "USER_ID" = $1
	`

	RetrievePasswordHashQuery = `
		SELECT "USER_ID", "PASSWORD"
		FROM "USER"
		WHERE "EMAIL" = $1
	`
)
