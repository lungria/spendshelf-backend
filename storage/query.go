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

// AsSQL formats query as sql code and params slice.
func (q Query) AsSQL() (string, []interface{}) {
	sqlBuilder := strings.Builder{}
	sqlParams := make([]interface{}, 0)

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

	return sqlBuilder.String(), sqlParams
}
