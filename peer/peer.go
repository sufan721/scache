package peer

// 根据key来选择节点
type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}

// 请求远程节点
type PeerGetter interface {
	Get(key string) (string, error)
}
