package zerocache

// 设计成只读，防止外部程序修改
type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

// 提供一个函数，返回b的深拷贝
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
