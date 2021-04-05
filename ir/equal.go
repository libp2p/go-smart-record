package ir

func IsEqual(x, y Node) bool {
	switch x1 := x.(type) {
	case Bool:
		switch y1 := y.(type) {
		case Bool:
			return IsEqualBool(x1, y1)
		}
	case String:
		switch y1 := y.(type) {
		case String:
			return IsEqualString(x1, y1)
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
	case Set:
		switch y1 := y.(type) {
		case Set:
			return IsEqualSet(x1, y1)
		}
	case Smart:
		switch y1 := y.(type) {
		case Smart:
			return IsEqualSmart(x1, y1)
		}
	}
	return false
}

func IsEqualType(x, y Node) bool {
	switch x.(type) {
	case Bool:
		switch y.(type) {
		case Bool:
			return true
		}
	case String:
		switch y.(type) {
		case String:
			return true
		}
	case Number:
		switch y.(type) {
		case Number:
			return true
		}
	case Blob:
		switch y.(type) {
		case Blob:
			return true
		}
	case Dict:
		switch y.(type) {
		case Dict:
			return true
		}
	case Set:
		switch y.(type) {
		case Set:
			return true
		}
	case Smart:
		switch y.(type) {
		case Smart:
			return true
		}
	}
	return false
}
