package crud

import (
	"bytes"
	"github.com/aymerick/raymond"
)

func ListRows(elements Elements, rows []row, options *raymond.Options) string {
	var buffer bytes.Buffer

	//orderedRow := []row{}
	//for _, col := range elements {
	//	cell, ok := row[col.ID]
	//	if !ok {
	//		continue
	//	}
	//
	//}

	for _, row := range rows {
		buffer.WriteString(options.FnWith(row))
	}

	return buffer.String()
}
