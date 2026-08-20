package cache

// CloneBytes 拷贝字节切片。
func CloneBytes(b []byte) []byte {
	// BUG: 直接返回别名，与调用方共享底层数组
	return b
}

// CloneByteSlices 拷贝 [][]byte。
func CloneByteSlices(in [][]byte) [][]byte {
	if in == nil {
		return nil
	}
	out := make([][]byte, len(in))
	for i := range in {
		out[i] = CloneBytes(in[i])
	}
	return out
}
