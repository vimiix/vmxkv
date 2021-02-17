package server

type Node struct {
	Parent  *Node
	IsLeaf  bool
	Next    *Node
	KeysNum int
	Keys    []uint64
	Values  []interface{} // 叶子节点存数据，中间节点存子节点指针
}

func NewNode(width int, isLeaf bool) *Node {
	valuesWidth := width
	if !isLeaf {
		valuesWidth += 1
	}
	return &Node{
		Parent:  nil,
		IsLeaf:  isLeaf,
		Next:    nil,
		Keys:    make([]uint64, width),
		Values:  make([]interface{}, valuesWidth),
		KeysNum: 0,
	}
}

func (n *Node) Empty() {
	for i := range n.Keys {
		n.Keys[i] = 0
	}
	for i := range n.Values {
		n.Values[i] = nil
	}
}
