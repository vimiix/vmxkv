package hotpool

type (
	// implement lru
	HotPool struct {
		cap        int
		buf        map[uint64]*Node
		head, tail *Node
	}

	Node struct {
		Key, Value uint64
		pre, next  *Node
	}
)

func NewHotPool(cap int) *HotPool {
	return &HotPool{
		cap: cap,
		buf: make(map[uint64]*Node, cap),
	}
}

func (p *HotPool) Get(key uint64) (val uint64, ok bool) {
	if node, ok := p.buf[key]; ok {
		p.refreshNode(node)
		return node.Value, ok
	}
	return
}

func (p *HotPool) Put(key, value uint64) {
	if node, ok := p.buf[key]; !ok {
		if len(p.buf) >= p.cap {
			p.removeNode(p.head)
		}
		node := &Node{
			Key:   key,
			Value: value,
		}
		p.addNode(node)
	} else {
		node.Value = value
		p.refreshNode(node)
	}
}

func (p *HotPool) Delete(key uint64) {
	if node, ok := p.buf[key]; ok {
		p.removeNode(node)
	}
}

func (p *HotPool) addNode(n *Node) {
	if p.tail != nil {
		p.tail.next = n
		n.pre = p.tail
		n.next = nil
	}
	p.tail = n
	if p.head == nil {
		p.head = n
		p.head.next = nil
	}
	p.buf[n.Key] = n
}

func (p *HotPool) removeNode(n *Node) {
	if n == p.tail {
		p.tail = p.tail.pre
	} else if n == p.head {
		p.head = p.head.next
	} else {
		n.pre.next = n.next
		n.next.pre = n.pre
	}
	delete(p.buf, n.Key)
}

func (p *HotPool) refreshNode(n *Node) {
	if p.tail == n {
		return
	}
	p.removeNode(n)
	p.addNode(n)
}
