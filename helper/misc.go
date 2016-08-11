package helper

// InUint 查询V值 是否在L列表中
func InUint(v uint, l []uint) bool {
	for _, b := range l {
		if b == v {
			return true
		}
	}
	return false
}
