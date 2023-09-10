package rbtree

import "unsafe"

// TODO :: Inline?
// refer : https://github.com/HuKeping/rbtree/blob/master/rbtree.go#L39

const (
	NODE_COLOR_BLACK = false
	NODE_COLOR_RED   = true
)

var nilNode = &Node{color: NODE_COLOR_BLACK}

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
	if node == nilNode {
		return nilNode
	}

	for node.left != nilNode {
		node = node.left
	}

	return node
}

func max(node *Node) *Node {
	if node == nilNode {
		return nilNode
	}

	for node.right != nilNode {
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
	root     *Node
	sentinel *Node
}

func (this *Tree) Insert(node *Node) {
	if node == nil {
		return
	}

	setNodeColorRed(node)
	this.insert(node)
}

func (this *Tree) rotateLeft(node *Node) {
	if node.right == nilNode {
		return
	}

	temp := node.right
	node.right = temp.left
	if temp.left != nilNode {
		temp.left.parent = node
	}
	temp.parent = node.parent

	if node.parent == nilNode {
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
	if node.left == nilNode {
		return
	}

	temp := node.left
	node.left = temp.right
	if temp.right != nilNode {
		temp.right.parent = node
	}
	temp.parent = node.parent

	if node.parent == nilNode {
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
	root := this.root
	temp := nilNode

	for root != nilNode {
		temp = root
		if node.key < root.key {
			root = root.left
		} else if root.key < node.key {
			root = root.right
		} else {
			break // return root?
		}
	}

	node.parent = temp
	if temp == nilNode {
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

	for current != nilNode {
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
	if node == nilNode {
		return nilNode
	}

	if node.right != nilNode {
		return min(node.right)
	}

	temp := node.parent
	if temp != nilNode && node == temp.right {
		node = temp
		temp = temp.parent
	}

	return temp
}

func (this *Tree) delete(key uint32) {
	node := this.lookup(key)
	if node == nilNode {
		return
	}

	var temp1, temp2 *Node
	if node.left == nilNode || node.right == nilNode {
		temp1 = node
	} else {
		temp1 = this.successor(node)
	}

	if temp1.left == nilNode {
		temp2 = temp1.right
	} else {
		temp2 = temp1.left
	}

	temp2.parent = temp1.parent

	if temp1.parent == nilNode {
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
