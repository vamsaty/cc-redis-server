package server

type RespType int

const (
	RespInvalid RespType = iota
	RespSimpleString
	RespSimpleError
	RespInteger
	RespBulkString
	RespArray
	RespNull
	RespBoolean
	RespDouble
	RespBigNumber
	RespBulkError
	RespVerbatimString
	RespMap
	RespSet
	RespPush
)

var RespTypeMap = map[string]RespType{
	"+": RespSimpleString,
	"-": RespSimpleError,
	":": RespInteger,
	"$": RespBulkString,
	"*": RespArray,
	"_": RespNull,
	"#": RespBoolean,
	",": RespDouble,
	"(": RespBigNumber,
	"!": RespBulkError,
	"=": RespVerbatimString,
	"%": RespMap,
	"~": RespSet,
	">": RespPush,
}

func GetRespType(symbol byte) (RespType, bool) {
	rType, found := RespTypeMap[string(symbol)]
	return rType, found
}

func GetRespTypeSymbol(rType RespType) string {
	for k, v := range RespTypeMap {
		if v == rType {
			return k
		}
	}
	return ""
}
