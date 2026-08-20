package cache

// CloneBytes 拷贝字节切片。
func CloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	cp := make([]byte, len(b))
	copy(cp, b)
	return cp
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
