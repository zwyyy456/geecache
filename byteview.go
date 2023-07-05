package geecache

// a byteview holds an immutable view of bytes

type ByteView struct {
	b []byte
}

// Len returns the view's length
func (v ByteView) Len() int {
	return len(v.b)
}
