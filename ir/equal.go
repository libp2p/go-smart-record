package ir

func IsEqual(x, y Node) bool {
	switch x1 := x.(type) {
	case String:
		switch y1 := y.(type) {
		case String:
			return IsEqualString(x1, y1)
		}
	case Int64:
		switch y1 := y.(type) {
		case Int64:
			return IsEqualInt64(x1, y1)
		}
	case Number:
		switch y1 := y.(type) {
		case Number:
			return IsEqualNumber(x1, y1)
		}
	case Blob:
		switch y1 := y.(type) {
		case Blob:
			return IsEqualBlob(x1, y1)
		}
	case Dict:
		switch y1 := y.(type) {
		case Dict:
			return IsEqualDict(x1, y1)
		}
	case Smart:
		switch y1 := y.(type) {
		case Smart:
			return IsEqualSmart(x1, y1)
		}
	}
	return false
}
