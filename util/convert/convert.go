package convert

//Int32ToBytes 转换
func Int32ToBytes(value int32) []byte {
	var data []byte
	data = make([]byte, 4)
	data[0] = byte(value >> 24)
	data[1] = byte(value >> 16)
	data[2] = byte(value >> 8)
	data[3] = byte(value)
	return data
}

//BytesToInt32 转换
func BytesToInt32(data []byte) (value int32) {
	value = int32(data[0])
	value = value << 8
	value += int32(data[1])
	value = value << 8
	value += int32(data[2])
	value = value << 8
	value += int32(data[3])
	return value
}
