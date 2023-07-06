package geecache

// PeePicker is the interface that must be implemented to locate
// the peer that owns a specify key.

type PeerPicker interface { // 选择结点
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer
type PeerGetter interface { // 从节点对应 group 中查找值
	Get(group string, key string) ([]byte, error)
}
