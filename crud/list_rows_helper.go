package crud

import (
	"bytes"
	"github.com/mergermarket/raymond"
)

func ListRows(elements Elements, rows []row, options *raymond.Options) string {
	var buffer bytes.Buffer

	var primaryKey string
	for _, element := range elements {
		if element.PrimaryKey == true {
			primaryKey = element.ID
		}
	}

	for _, row := range rows {
		ID, ok := row[primaryKey]
		if !ok {
			ID = 0
		}
		ctx := map[string]interface{}{
			"ID":  ID,
			"Row": row,
		}

		buffer.WriteString(options.FnWith(ctx))
	}

	return buffer.String()
}
