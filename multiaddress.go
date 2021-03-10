package sr

type Multiaddress struct {
	Multiaddress string
	Dict
}

func (m Multiaddress) AsDict() Dict {
	return m.Dict.CopySet(String{m.Multiaddress}, String{m.Multiaddress})
}
