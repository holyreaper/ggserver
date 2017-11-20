package convert

//Int32ToBytes 转换
func Int32ToBytes(value int32) []byte {
	var data []byte
	data = make([]byte, 4)
	data[0] = byte(value >> 24)
	data[1] = byte(value >> 24)
	data[2] = byte(value >> 24)
	data[3] = byte(value)
	return data
}
