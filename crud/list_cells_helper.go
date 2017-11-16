package crud

import (
	"bytes"
	"github.com/aymerick/raymond"
)

func ListCells(elements Elements, row row, options *raymond.Options) string {
	var buffer bytes.Buffer

	for _, col := range elements {
		cell, ok := row[col.ID]
		if !ok {
			cell = ""
		}
		buffer.WriteString(options.FnWith(cell))
	}

	return buffer.String()
}
