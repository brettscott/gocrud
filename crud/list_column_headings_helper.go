package crud

import (
	"bytes"
	"github.com/mergermarket/raymond"
)


// EachIndividualDescending combines all individuals liquidity events into one collection for rendering
func ListColumnHeadings(elements Elements, options *raymond.Options) string {
	var buffer bytes.Buffer
	for _, c := range elements {
		buffer.WriteString(options.FnWith(c))
	}
	return buffer.String()
}
