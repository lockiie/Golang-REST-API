package functions

//BoolToByte converter um boleano em um byte
func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

//BoolStrToByte converter uma string boleana em byte
func BoolStrToByte(s string) byte {
	switch s {
	case "true":
		return 1
	case "false":
		return 0
	default:
		return 10
	}
}
