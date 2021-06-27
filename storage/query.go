package storage

import (
	"fmt"
	"strings"
	"time"
)

// Query fields for UpdateTransactionCommand. Zero value fields are treated as absent.
type Query struct {
	ID            string
	CategoryID    int32
	LastUpdatedAt time.Time
}

// appendToSQL formats query to SQL and adds it to existing sqlBuilder.
// Returns updated sqlParams slice with all added parameters for query.
func (q Query) appendToSQL(sqlBuilder *strings.Builder, sqlParams []interface{}) []interface{} {
	sqlBuilder.WriteString("WHERE ")

	if q.ID != "" {
		sqlParams = append(sqlParams, q.ID)
		sqlBuilder.WriteString(fmt.Sprintf(`"ID" = $%v `, len(sqlParams)))
	}

	if q.CategoryID != 0 {
		if len(sqlParams) > 0 {
			sqlBuilder.WriteString("AND ")
		}

		sqlParams = append(sqlParams, q.CategoryID)
		sqlBuilder.WriteString(fmt.Sprintf(`"categoryID" = $%v `, len(sqlParams)))
	}

	zeroTimeValue := time.Time{}
	if q.LastUpdatedAt != zeroTimeValue {
		if len(sqlParams) > 0 {
			sqlBuilder.WriteString("AND ")
		}

		sqlParams = append(sqlParams, q.LastUpdatedAt)
		sqlBuilder.WriteString(fmt.Sprintf(`"lastUpdatedAt" = $%v `, len(sqlParams)))
	}

	return sqlParams
}
