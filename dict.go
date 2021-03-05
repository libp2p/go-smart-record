package cr

// Dict is a set of uniquely-named child nodes.
type Dict struct {
	Children map[string]Node
}

func (r Record) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return XXX
}

func MergeDicts(x, y *Dict) Node {
	XXX
}
