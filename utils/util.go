package utils

var orderType = map[string]bool{
	"desc": true,
	"asc":  true,
}

func ValidOrderType(oType string) bool {
	_, ok := orderType[oType]
	return ok
}
