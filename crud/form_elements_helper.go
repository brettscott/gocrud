package crud

import (
	"bytes"
	"github.com/mergermarket/raymond"
)

//// TODO this could be done much much earlier to avoid having to send Elements & row around separately.
//type Payload struct {
//	*Element
//	Value interface{}
//}

func FormElements(evs []ElementValue, options *raymond.Options) string {
	var buffer bytes.Buffer

	//return fmt.Sprintf("evs: %+v", evs)

	for _, ev := range evs {
		buffer.WriteString(options.FnWith(ev))
	}
	return buffer.String()
}
