package db

const (
	SelectResouceTotalCost = `
		SELECT "COST", COUNT(*) AS occurrences
		FROM "COST"
	    JOIN "RESOURCE" ON "COST"."RESOURCE_ID" = "RESOURCE"."RESOURCE_ID"
		WHERE "COST"."RESOURCE_ID" = $1
	    AND "RESOURCE"."COMPANY_ID" = $2
		AND "START_TIMESTAMP" > $3
		AND "END_TIMESTAMP" < $4
		GROUP BY "COST";
	`

	SelectResouceCost = `
		SELECT *
		FROM "COST"
	    JOIN "RESOURCE" ON "COST"."RESOURCE_ID" = "RESOURCE"."RESOURCE_ID"
		WHERE "COST"."RESOURCE_ID" = $1
	    AND "RESOURCE"."COMPANY_ID" = $2
		AND "START_TIMESTAMP" > $3
		AND "END_TIMESTAMP" < $4
		GROUP BY "COST";
	`

	SelectResouceData = `
		SELECT
		  "METRIC_NAME",
		  "METRIC_VALUE",
		  "METRIC_UNIT",
		  "TIMESTAMP"
		FROM
		  "RESOURCE_DATA"
		JOIN
		  "RESOURCE" ON "RESOURCE_DATA"."RESOURCE_ID" = "RESOURCE"."RESOURCE_ID"
		WHERE
		  "RESOURCE_DATA"."RESOURCE_ID" = $1
		AND
		  "RESOURCE"."COMPANY_ID" = $2
		ORDER BY "TIMESTAMP" ASC;
	`
)
