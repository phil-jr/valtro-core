package db

const (
	SelectResouceCost = `
		SELECT "COST", COUNT(*) AS occurrences
		FROM "COST"
		WHERE "RESOURCE_ID" = $1
		GROUP BY "COST";
	`
)
