package cr

import (
	"encoding/xml"
)

interface Node {
	xml.Marshaler
}

type Update struct {
	Command string // e.g. edit, add
}
