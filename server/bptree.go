package server

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
	ErrKeyExist = errors.New("key already exists")
)

type BPTree struct {
	Version int64
	Root    *Node
	Degree  int
}

func NewBPTree(degree int) *BPTree {
	return &BPTree{
		Version: time.Now().Unix(),
		Root:    nil,
		Degree:  degree,
	}
}

type kv struct {
	K, V uint64
}

func NewFromFile(r io.Reader) (bpt *BPTree, err error) {
	var kvs []*kv
	decoder := gob.NewDecoder(r)
	if err = decoder.Decode(&kvs); err != nil {
		return
	}
	version := kvs[0].V
	degree := kvs[1].V
	bpt = &BPTree{
		Version: int64(version),
		Degree:  int(degree),
	}
	for _, item := range kvs[2:] {
		bpt.Insert(item.K, item.V)
	}
	return
}

func (t *BPTree) updateVersion() {
	t.Version = time.Now().Unix()
}

func (t *BPTree) PrintAll() {
	n := t.Root
	for !n.IsLeaf {
		n = n.Values[0].(*Node)
	}
	for n != nil {
		for i := 0; i < n.KeysNum; i++ {
			fmt.Printf("<%d:%d>\n", n.Keys[i], n.Values[i].(uint64))
		}
		n = n.Next
	}
}

func (t *BPTree) Dump(w io.Writer) (err error) {
	n := t.Root
	for !n.IsLeaf {
		n = n.Values[0].(*Node)
	}
	var kvs = []*kv{
		{K: 0, V: uint64(t.Version)},
		{K: 0, V: uint64(t.Degree)},
	}
	for n != nil {
		for i := 0; i < n.KeysNum; i++ {
			kvs = append(kvs, &kv{K: n.Keys[i], V: n.Values[i].(uint64)})
		}
		n = n.Next
	}
	encoder := gob.NewEncoder(w)
	err = encoder.Encode(kvs)
	return
}

func (t *BPTree) Height() int {
	h := 1
	n := t.Root
	for !n.IsLeaf {
		n = n.Values[0].(*Node)
		h++
	}
	return h
}

// Upsert 插入
func (t *BPTree) Insert(key, value uint64) (err error) {
	if _, err = t.Find(key); err == nil {
		return ErrKeyExist
	}
	if t.Root == nil {
		t.newTree(key, value)
		t.updateVersion()
		return
	}

	n := t.findLeafNode(key)
	if n.KeysNum < t.Degree-1 {
		t.insertIntoLeaf(n, key, value)
	} else {
		t.splitLeafAndInsert(n, key, value)
	}
	t.updateVersion()
	return
}

// Find 查询
func (t *BPTree) Find(key uint64) (value uint64, err error) {
	n := t.findLeafNode(key)
	if n == nil {
		err = ErrNotFound
		return
	}
	for i, k := range n.Keys {
		if key == k {
			value = n.Values[i].(uint64)
			return
		}
	}
	err = ErrNotFound
	return
}

func (t *BPTree) RangeFind(startKey, endKey uint64, rowHandler func(k, v uint64) bool) (err error) {
	count, keys, values := t.rangeFind(startKey, endKey)
	if count == 0 {
		err = ErrNotFound
		return
	}
	for i := 0; i < count; i++ {
		rowHandler(keys[i], values[i])
	}
	return
}

func (t *BPTree) Delete(key uint64) (err error) {
	var value uint64
	value, err = t.Find(key)
	if err != nil {
		return
	}
	leaf := t.findLeafNode(key)
	if leaf != nil {
		t.deleteKey(leaf, key, value)
	}
	t.updateVersion()
	return
}

func (t *BPTree) newTree(key, value uint64) {
	n := NewNode(t.Degree, true)
	n.Keys[0] = key
	n.Values[0] = value
	n.KeysNum++
	t.Root = n
	return
}

func (t *BPTree) findLeafNode(key uint64) (n *Node) {
	n = t.Root
	if n == nil {
		return
	}

	for !n.IsLeaf {
		var idx int
		for idx < n.KeysNum {
			if key >= n.Keys[idx] {
				idx++
			} else {
				break
			}
		}
		n = n.Values[idx].(*Node)
	}
	return
}

func (t *BPTree) insertIntoLeaf(n *Node, key, value uint64) {
	var position int
	for position < n.KeysNum && n.Keys[position] < key {
		position++
	}
	for i := n.KeysNum; i > position; i-- {
		n.Keys[i] = n.Keys[i-1]
		n.Values[i] = n.Values[i-1]
	}
	n.Keys[position] = key
	n.Values[position] = value
	n.KeysNum++
	return
}

func (t *BPTree) splitLeafAndInsert(oldLeaf *Node, key, value uint64) {
	newLeaf := NewNode(t.Degree, true)
	split := t.Degree / 2
	var position int
	for position < oldLeaf.KeysNum && oldLeaf.Keys[position] < key {
		position++
	}
	tmpKeys := make([]uint64, oldLeaf.KeysNum+1)
	tmpValues := make([]interface{}, oldLeaf.KeysNum+1)
	var tmpIdx int
	for i := 0; i < oldLeaf.KeysNum; i++ {
		if i == position {
			tmpIdx++
		}
		tmpKeys[tmpIdx] = oldLeaf.Keys[i]
		tmpValues[tmpIdx] = oldLeaf.Values[i]
		tmpIdx++
	}
	tmpKeys[position] = key
	tmpValues[position] = value

	// clear old leaf
	oldLeaf.KeysNum = 0
	oldLeaf.Empty()

	for i, k := range tmpKeys {
		if i < split {
			oldLeaf.Keys[oldLeaf.KeysNum] = k
			oldLeaf.Values[oldLeaf.KeysNum] = tmpValues[i]
			oldLeaf.KeysNum++
		} else {
			newLeaf.Keys[newLeaf.KeysNum] = k
			newLeaf.Values[newLeaf.KeysNum] = tmpValues[i]
			newLeaf.KeysNum++
		}
	}

	newLeaf.Next = oldLeaf.Next
	oldLeaf.Next = newLeaf
	newLeaf.Parent = oldLeaf.Parent
	newKey := newLeaf.Keys[0]
	t.insertIntoParent(oldLeaf, newKey, newLeaf)
	return
}

func (t *BPTree) insertIntoParent(left *Node, key uint64, right *Node) {
	parent := left.Parent
	if parent == nil {
		n := NewNode(t.Degree, false)
		n.Keys[0] = key
		n.Values[0] = left
		n.Values[1] = right
		n.KeysNum++
		left.Parent = n
		right.Parent = n
		t.Root = n
		return
	}
	var leftIdx int
	for leftIdx <= parent.KeysNum && parent.Values[leftIdx] != left {
		leftIdx++
	}
	if parent.KeysNum < t.Degree-1 {
		t.insertIntoNode(parent, leftIdx, key, right)
	} else {
		t.splitNodeAndInsert(parent, leftIdx, key, right)
	}
}

func (t *BPTree) insertIntoNode(n *Node, leftIdx int, key uint64, right *Node) {
	for i := n.KeysNum; i > leftIdx; i-- {
		n.Values[i+1] = n.Values[i]
		n.Keys[i] = n.Keys[i-1]
	}
	n.Values[leftIdx+1] = right
	n.Keys[leftIdx] = key
	n.KeysNum++
}

func (t *BPTree) splitNodeAndInsert(n *Node, leftIdx int, key uint64, right *Node) {
	var (
		tmpKeys   = make([]uint64, n.KeysNum+1)
		tmpValues = make([]interface{}, n.KeysNum+2)
		tmpIdx    int
	)
	for i := 0; i < n.KeysNum; i++ {
		if i == leftIdx {
			tmpIdx++
		}
		tmpKeys[tmpIdx] = n.Keys[i]
		tmpIdx++
	}
	tmpIdx = 0
	for i := 0; i < n.KeysNum+1; i++ {
		if i == leftIdx+1 {
			tmpIdx++
		}
		tmpValues[tmpIdx] = n.Values[i]
		tmpIdx++
	}
	tmpKeys[leftIdx] = key
	tmpValues[leftIdx+1] = right

	n.KeysNum = 0
	n.Empty()
	split := t.Degree / 2
	newNode := NewNode(t.Degree, false)
	for i := 0; i < split; i++ {
		n.Keys[i] = tmpKeys[i]
		n.Values[i] = tmpValues[i]
		n.KeysNum++
	}
	n.Values[split] = tmpValues[split]
	liftKey := tmpKeys[split]
	for i := split + 1; i < len(tmpKeys); i++ {
		newNode.Keys[newNode.KeysNum] = tmpKeys[i]
		newNode.Values[newNode.KeysNum] = tmpValues[i]
		newNode.KeysNum++
	}
	newNode.Values[newNode.KeysNum] = tmpValues[len(tmpValues)-1]
	newNode.Parent = n.Parent
	for i := 0; i <= newNode.KeysNum; i++ {
		child := newNode.Values[i].(*Node)
		child.Parent = newNode
	}
	t.insertIntoParent(n, liftKey, newNode)
}

func (t *BPTree) rangeFind(startKey, endKey uint64) (count int, keys, values []uint64) {
	leaf := t.findLeafNode(startKey)
	if leaf == nil {
		return
	}
	var idx int
	for ; idx < leaf.KeysNum && leaf.Keys[idx] < startKey; idx++ {
		if idx == leaf.KeysNum {
			return
		}
	}
	for leaf != nil {
		for ; idx < leaf.KeysNum && leaf.Keys[idx] <= endKey; idx++ {
			keys = append(keys, leaf.Keys[idx])
			values = append(values, leaf.Values[idx].(uint64))
			count++
		}
		leaf = leaf.Next
		idx = 0
	}
	return
}

func (t *BPTree) deleteKey(n *Node, key, value interface{}) {
	t.removeKeyFromNode(n, key, value)
	fmt.Printf("remove key[%+v] from node:%+v\n", key, n.Keys)
	if n == t.Root {
		t.adjustRoot()
		return
	}
	var minKeys int
	if n.IsLeaf {
		minKeys = (t.Degree - 1) / 2
	} else {
		minKeys = t.Degree/2 - 1
	}

	if n.KeysNum >= minKeys {
		return
	}

	// 合并或重建子树
	var preIdx int
	for i := 0; i < n.Parent.KeysNum+1; i++ { // 遍历父节点的所有子节点，找到改节点的左边节点
		if reflect.DeepEqual(n.Parent.Values[i], n) {
			preIdx = i - 1
		}
	}
	var (
		liftKeyIdx, capacity int
		preNode              *Node
	)
	if preIdx == -1 {
		liftKeyIdx = 0
	} else {
		liftKeyIdx = preIdx
	}
	liftKey := n.Parent.Keys[liftKeyIdx]
	if preIdx == -1 {
		preNode = n.Parent.Values[1].(*Node)
	} else {
		preNode = n.Parent.Values[preIdx].(*Node)
	}
	if n.IsLeaf {
		capacity = t.Degree
	} else {
		capacity = t.Degree - 1
	}
	if preNode.KeysNum+n.KeysNum < capacity {
		t.mergeNodes(n, preNode, preIdx, liftKey)
	} else {
		t.rebuildNodes(n, preNode, preIdx, liftKeyIdx, liftKey)
	}
}

func (t *BPTree) removeKeyFromNode(n *Node, key, value interface{}) {
	var i int
	for n.Keys[i] != key {
		i += 1
	}
	for i += 1; i <= n.KeysNum; i++ {
		n.Keys[i-1] = n.Keys[i]
		n.Keys[i] = 0
	}

	var valuesNum int
	if n.IsLeaf {
		valuesNum = n.KeysNum
	} else {
		valuesNum = n.KeysNum + 1
	}

	i = 0
	for n.Values[i] != value {
		i += 1
	}
	for i += 1; i <= valuesNum; i++ {
		n.Values[i-1] = n.Values[i]
		n.Values[i] = nil
	}
	n.KeysNum--
	return
}

func (t *BPTree) adjustRoot() {
	if t.Root.KeysNum > 0 {
		return
	}

	if !t.Root.IsLeaf {
		root := t.Root.Values[0].(*Node)
		root.Parent = nil
		t.Root = root
	} else {
		t.Root = nil
	}
}

func (t *BPTree) mergeNodes(n, preNode *Node, preIdx int, liftKey uint64) {
	if preIdx == -1 {
		n, preNode = preNode, n
	}
	if !n.IsLeaf {
		preNode.Keys[preNode.KeysNum] = liftKey
		preNode.KeysNum++

		i := preNode.KeysNum
		j := 0
		for j = 0; j < n.KeysNum; j++ {
			preNode.Keys[i] = n.Keys[j]
			preNode.Values[i] = n.Values[j]
			preNode.KeysNum++
			n.KeysNum--
			i++
		}
		preNode.Values[i] = n.Values[j]
		for _, v := range preNode.Values {
			v.(*Node).Parent = preNode
		}

	} else {
		i := preNode.KeysNum
		for j := 0; j < n.KeysNum; j++ {
			preNode.Keys[i] = n.Keys[j]
			preNode.Values[i] = n.Values[j]
			i++
			preNode.KeysNum++
		}
	}
	t.deleteKey(n.Parent, liftKey, n)
}

func (t *BPTree) rebuildNodes(n, preNode *Node, preIdx, liftKeyIdx int, liftKey uint64) {
	var (
		i int
	)
	if preIdx == -1 {
		if !n.IsLeaf {
			n.Values[n.KeysNum+1] = n.Values[n.KeysNum]
		}
		for i = n.KeysNum; i > 0; i-- {
			n.Keys[i] = n.Keys[i-1]
			n.Values[i] = n.Values[i-1]
		}
		if !n.IsLeaf {
			n.Values[0] = preNode.Values[preNode.KeysNum]
			n.Values[0].(*Node).Parent = n
			preNode.Values[preNode.KeysNum] = nil
			n.Keys[0] = liftKey
			n.Parent.Keys[liftKeyIdx] = preNode.Keys[preNode.KeysNum-1]
		} else {
			n.Values[0] = preNode.Values[preNode.KeysNum-1]
			preNode.Values[preNode.KeysNum-1] = nil
			n.Keys[0] = preNode.Keys[preNode.KeysNum-1]
			n.Parent.Keys[liftKeyIdx] = n.Keys[0]
		}
	} else {
		if n.IsLeaf {
			n.Keys[n.KeysNum] = preNode.Keys[0]
			n.Values[n.KeysNum] = preNode.Values[0]
			n.Parent.Keys[liftKeyIdx] = preNode.Keys[1]
		} else {
			n.Keys[n.KeysNum] = liftKey
			n.Values[n.KeysNum+1] = preNode.Values[0]
			n.Values[n.KeysNum+1].(*Node).Parent = n
			n.Parent.Keys[liftKeyIdx] = preNode.Keys[0]
		}
		for i = 0; i < preNode.KeysNum-1; i++ {
			preNode.Keys[i] = preNode.Keys[i+1]
			preNode.Values[i] = preNode.Values[i+1]
		}
		if !n.IsLeaf {
			preNode.Values[i] = preNode.Values[i+1]
		}
	}
	n.KeysNum++
	preNode.KeysNum--
}
