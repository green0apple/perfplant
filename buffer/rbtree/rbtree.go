package rbtree

import "unsafe"

// TODO :: Inline?

const (
	NODE_COLOR_BLACK = false
	NODE_COLOR_RED   = true
)

type Node struct {
	key    uint32
	val    unsafe.Pointer
	left   *Node
	right  *Node
	parent *Node
	color  bool
}

func NewNode(key uint32, val unsafe.Pointer) *Node {
	return &Node{
		key: key,
		val: val,
	}
}

func (this *Node) setRed()       { this.color = NODE_COLOR_RED }
func (this *Node) setBlock()     { this.color = NODE_COLOR_BLACK }
func (this *Node) isRed() bool   { return this.color }
func (this *Node) isBlack() bool { return !this.color }

type Tree struct {
	root     *Node
	sentinel *Node
}

func (this *Tree) Insert(node *Node) {
	if this.root == nil {
		this.root = node
		node.setBlock()
		return
	}

	/*
		#define likely(x)       __builtin_expect((x),1)
		#define unlikely(x)     __builtin_expect((x),0)
	*/
	if node.parent.isRed() {

	}

}
