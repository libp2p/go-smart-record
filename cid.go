package sr

type Cid struct {
	Cid string
	Dict
}

func (c Cid) AsDict() Dict {
	return c.Dict.CopySet(String{c.Cid}, String{c.Cid})
}
