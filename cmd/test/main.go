package main

import (
"fmt"
)

const M = 4

type node struct {
	keys     []int
	children []*node
}

func newNode() *node { return &node{keys: make([]int, 0), children: nil} }

func (n *node) isLeaf() bool { return n.children == nil }

func (n *node) findKeyIdx(k int) int {
	l, r := 0, len(n.keys)
	for l < r {
		m := (l + r) / 2
		if n.keys[m] < k {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

type BTree struct{ root *node }

func NewBTree() *BTree { return &BTree{root: newNode()} }

func (t *BTree) Insert(k int) { // insert key
	if t.root == nil {
		t.root = newNode()
	}
	if len(t.root.keys) == M-1 { // root full -> split
		s := &node{
			keys:     []int{t.root.keys[M/2]},
			children: []*node{t.root, t.splitChild(t.root)},
		}
		t.root = s
	}
	t.insertNonFull(t.root, k)
}

func (t *BTree) Search(k int) bool { return t.searchRec(t.root, k) }

func (t *BTree) Delete(k int) { t.deleteRec(t.root, k) }

func (t *BTree) insertNonFull(n *node, k int) {
	i := n.findKeyIdx(k)
	if n.isLeaf() {
		n.keys = append(n.keys, 0)
		copy(n.keys[i+1:], n.keys[i:])
		n.keys[i] = k
		return
	}
	child := n.children[i]
	if len(child.keys) == M-1 { // child full -> split
		newChild := t.splitChild(child)
		t.insertNonFull(child, k) // after split child may still be full but will go into correct side
	} else {
		t.insertNonFull(child, k)
	}
}

func (t *BTree) splitChild(n *node) *node { // splits node n and returns right sibling
	mid := M / 2
	right := &node{
		keys:     append([]int(nil), n.keys[mid+1:]...), // copy keys after mid
		children: nil,
	}
	if !n.isLeaf() {
		right.children = append([]*node(nil), n.children[mid+1:]...)
		n.children = n.children[:mid+1]
	}
	n.keys = n.keys[:mid] // keep left part
	return right
}

func (t *BTree) searchRec(n *node, k int) bool {
	if n == nil { return false }
	i := n.findKeyIdx(k)
	if i < len(n.keys) && n.keys[i] == k { return true }
	if n.isLeaf() { return false }
	return t.searchRec(n.children[i], k)
}

func (t *BTree) deleteRec(n *node, k int) {
	// सरल: यदि key leaf में है तो हटाएँ; अन्यथा predecessor से बदलें
	i := n.findKeyIdx(k)
	if i < len(n.keys) && n.keys[i] == k { // key found in this node
		if n.isLeaf() {
			n.keys = append(n.keys[:i], n.keys[i+1:]...)
			return
		}
		// replace with predecessor (max of left subtree)
		pred := t.maxKey(n.children[i])
		n.keys[i] = pred
		t.deleteRec(n.children[i], pred)
	} else if !n.isLeaf() { // recurse into child
		t.deleteRec(n.children[i], k)
	}
}

func (t *BTree) maxKey(n *node) int {
	for !n.isLeaf() {
		n = n.children[len(n.children)-1]
	}
	return n.keys[len(n.keys)-1]
}

// ---------- उपयोग ----------
func main() {
	b := NewBTree()
	keys := []int{10, 20, 5, 6, 12, 30, 7, 17}

	for _, k := range keys {
		b.Insert(k)
	}
	fmt.Println("सभी keys डाल दिए गए।")

	testKeys := []int{6, 15, 17}
	for _, k := range testKeys {
		if b.Search(k) {
			fmt.Printf("%d पाया गया!\n", k)
		} else {
			fmt.Printf("%d नहीं मिला।\n", k)
		}
	}

	b.Delete(6)
	fmt.Println("6 हटाया गया।")
	if !b.Search(6) {
		fmt.Println("पुष्टि: 6 अब मौजूद नहीं है।")
	}
}
