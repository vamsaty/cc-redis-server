package resp

/* NOTE: This is unused as of 10/12/2023 (dd/mm/yyyy) */

type DataType int

const (
	DataInvalid DataType = iota
	DataSimpleString
	DataSimpleError
	DataInteger
	DataBulkString
	DataArray
	DataNull
	DataBoolean
	DataDouble
	DataBigNumber
	DataBulkError
	DataVerbatimString
	DataMap
	DataSet
	DataPush
)

var DataTypeMap = map[string]DataType{
	"+": DataSimpleString,
	"-": DataSimpleError,
	":": DataInteger,
	"$": DataBulkString,
	"*": DataArray,
	"_": DataNull,
	"#": DataBoolean,
	",": DataDouble,
	"(": DataBigNumber,
	"!": DataBulkError,
	"=": DataVerbatimString,
	"%": DataMap,
	"~": DataSet,
	">": DataPush,
}

func GetDataType(symbol byte) (DataType, bool) {
	rType, found := DataTypeMap[string(symbol)]
	return rType, found
}

func GetDataTypeSymbol(rType DataType) string {
	for k, v := range DataTypeMap {
		if v == rType {
			return k
		}
	}
	return ""
}
