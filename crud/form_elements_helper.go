package crud

import "bytes"
import "github.com/mergermarket/raymond"

// TODO this could be done much much earlier to avoid having to send Elements & row around separately.
type Payload struct {
	*Element
	Value interface{}
}

func FormElements(elements Elements, row row, options *raymond.Options) string {
	var buffer bytes.Buffer
	for _, e := range elements {
		payload := Payload{}
		if _, ok := row[e.ID]; ok {
			payload = Payload{
				&e,
				row[e.ID],
			}
		}
		buffer.WriteString(options.FnWith(payload))
	}
	return buffer.String()
}
