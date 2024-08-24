package library

var Data = map[string]string{}

func SetData(key string, value string) {
	Data[key] = value
}

func GetData(key string) (string, bool) {
	value, ok := Data[key]
	return value, ok
}

func DeleteData(key string) {
	delete(Data, key)
}
