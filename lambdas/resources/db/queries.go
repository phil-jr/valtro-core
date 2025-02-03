package db

const (
	SelectResouceCost = `
		SELECT "COST", COUNT(*) AS occurrences
		FROM "COST"
	    JOIN "RESOURCE" ON "COST"."RESOURCE_ID" = "RESOURCE"."RESOURCE_ID"
		WHERE "COST"."RESOURCE_ID" = $1
	    AND "RESOURCE"."COMPANY_ID" = $2
		AND "START_TIMESTAMP" > $3
		AND "END_TIMESTAMP" < $4
		GROUP BY "COST";
	`
)
