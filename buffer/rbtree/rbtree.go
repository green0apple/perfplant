package rbtree

import "unsafe"

// TODO :: Inline?
// refer : https://github.com/HuKeping/rbtree/blob/master/rbtree.go#L39

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

func setNodeColorRed(node *Node)   { node.color = NODE_COLOR_RED }
func setNodeColorBlack(node *Node) { node.color = NODE_COLOR_BLACK }
func isNodeRed(node *Node) bool    { return node.color == NODE_COLOR_RED }
func isNodeBlack(node *Node) bool  { return node.color == NODE_COLOR_BLACK }

func min(node *Node) *Node {
	if node == nil {
		return nil
	}

	for node.left != nil {
		node = node.left
	}

	return node
}

func max(node *Node) *Node {
	if node == nil {
		return nil
	}

	for node.right != nil {
		node = node.right
	}

	return node
}

func NewNode(key uint32, val unsafe.Pointer) *Node {
	return &Node{
		key: key,
		val: val,
	}
}

type Tree struct {
	root *Node
}

func NewTree() *Tree {
	return &Tree{
		root: &Node{color: NODE_COLOR_BLACK},
	}
}

func (this *Tree) Insert(key uint32, val unsafe.Pointer) {
	node := NewNode(key, val)

	setNodeColorRed(node)
	this.insert(node)
}

func (this *Tree) rotateLeft(node *Node) {
	if node.right == nil {
		return
	}

	temp := node.right
	node.right = temp.left
	if temp.left != nil {
		temp.left.parent = node
	}
	temp.parent = node.parent

	if node.parent == nil {
		this.root = temp
	} else if node == node.parent.left {
		node.parent.left = temp
	} else {
		node.parent.right = temp
	}

	temp.left = node
	node.parent = temp
}

func (this *Tree) rotateRight(node *Node) {
	if node.left == nil {
		return
	}

	temp := node.left
	node.left = temp.right
	if temp.right != nil {
		temp.right.parent = node
	}
	temp.parent = node.parent

	if node.parent == nil {
		this.root = temp
	} else if node == node.parent.left {
		node.parent.left = temp
	} else {
		node.parent.right = temp
	}

	temp.right = node
	node.parent = temp
}

func (this *Tree) insert(node *Node) {
	var (
		root *Node = this.root
		temp *Node
	)

	for root != nil {
		temp = root
		if node.key < root.key {
			root = root.left
		} else if root.key < node.key {
			root = root.right
		} else {
			return // return root?
		}
	}

	node.parent = temp
	if temp == nil {
		this.root = node
	} else if node.key < temp.key {
		temp.left = node
	} else {
		temp.right = node
	}

	this.insertFixup(node)
	// return node?
}

func (this *Tree) insertFixup(node *Node) {
	if isNodeRed(node.parent) {
		if node.parent == node.parent.parent.left {
			temp := node.parent.parent.right
			if isNodeRed(temp) {
				setNodeColorBlack(node.parent)
				setNodeColorBlack(temp)
				setNodeColorRed(node.parent.parent)
				node = node.parent.parent
			} else {
				if node == node.parent.right {
					node = node.parent
					this.rotateLeft(node)
				}
			}

			setNodeColorBlack(node.parent)
			setNodeColorRed(node.parent.parent)
			this.rotateRight(node.parent.parent)
		}
	} else {
		temp := node.parent.parent.left
		if isNodeRed(temp) {
			setNodeColorBlack(node.parent)
			setNodeColorBlack(temp)
			setNodeColorRed(node.parent.parent)
			node = node.parent.parent
		} else {
			if node == node.parent.left {
				node = node.parent
				this.rotateRight(node)
			}

			setNodeColorBlack(node.parent)
			setNodeColorRed(node.parent.parent)
			this.rotateLeft(node.parent.parent)
		}
	}

	setNodeColorBlack(this.root)
}

func (this *Tree) lookup(key uint32) *Node {
	current := this.root

	for current != nil {
		if current.key < key {
			current = current.right
		} else if current.key > key {
			current = current.left
		} else {
			break
		}
	}

	return current
}

func (this *Tree) successor(node *Node) *Node {
	if node == nil {
		return nil
	}

	if node.right != nil {
		return min(node.right)
	}

	temp := node.parent
	if temp != nil && node == temp.right {
		node = temp
		temp = temp.parent
	}

	return temp
}

func (this *Tree) delete(key uint32) {
	node := this.lookup(key)
	if node == nil {
		return
	}

	var temp1, temp2 *Node
	if node.left == nil || node.right == nil {
		temp1 = node
	} else {
		temp1 = this.successor(node)
	}

	if temp1.left == nil {
		temp2 = temp1.right
	} else {
		temp2 = temp1.left
	}

	temp2.parent = temp1.parent

	if temp1.parent == nil {
		this.root = temp2
	} else if temp1 == temp1.parent.left {
		temp1.parent.left = temp2
	} else {
		temp1.parent.right = temp2
	}

	if temp1 != node {
		node.key = temp1.key
	}

	if isNodeBlack(temp1) {
		this.deleteFixup(temp2)
	}
}

func (this *Tree) deleteFixup(node *Node) {
	if node != this.root && isNodeBlack(node) {
		if node == node.parent.left {
			temp := node.parent.right
			if isNodeRed(temp) {
				setNodeColorRed(temp)
				setNodeColorRed(node.parent)
				this.rotateLeft(node.parent)
				temp = node.parent.right
			}

			if isNodeBlack(temp.left) && isNodeRed(temp.right) {
				setNodeColorRed(temp)
				node = node.parent
			} else {
				if isNodeBlack(temp.right) {
					setNodeColorBlack(node.left)
					setNodeColorRed(temp)
					this.rotateRight(temp)
					temp = node.parent.right
				}

				temp.color = node.parent.color
				setNodeColorBlack(node.parent)
				setNodeColorBlack(temp.right)
				this.rotateLeft(node.parent)

				// exit loop
				node = this.root
			}
		} else {
			temp := node.parent.left
			if isNodeRed(temp) {
				setNodeColorBlack(temp)
				setNodeColorRed(node)
				this.rotateRight(node.parent)
				temp = node.parent.left
			}

			if isNodeBlack(temp.left) && isNodeBlack(temp.right) {
				setNodeColorRed(temp)
				node = node.parent
			} else {
				if isNodeBlack(temp.left) {
					setNodeColorBlack(temp.right)
					setNodeColorRed(temp)
					this.rotateLeft(temp)
					temp = node.parent.left
				}

				temp.color = node.parent.color
				setNodeColorBlack(node.parent)
				setNodeColorBlack(temp.left)
				this.rotateRight(node.parent)

				node = this.root
			}
		}
	}

	setNodeColorBlack(node)
}
