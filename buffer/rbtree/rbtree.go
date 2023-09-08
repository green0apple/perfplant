package rbtree

import "unsafe"

const (
	NODE_COLOR_BLACK = false
	NODE_COLOR_RED   = true
)

type Node struct {
	key    uint
	val    unsafe.Pointer
	left   *Node
	right  *Node
	parent *Node
	color  bool
}

type Tree struct {
	Insert func(root, node, sentinel *Node)

	root     *Node
	sentinel *Node
}

func isRed(node *Node) bool   { return node.color }
func isBlack(node *Node) bool { return !node.color }

// inline
func (t *Tree) min(node, sentinel *Node) *Node {
	for node.left != sentinel {
		node = node.left
	}

	return node
}

func (t *Tree) insert(new *Node) {
	var (
		root     **Node
		temp     *Node
		sentinel *Node
	)

	root = &t.root
	sentinel = t.sentinel

	if *root == sentinel {
		new.parent = nil
		new.left = sentinel
		new.right = sentinel
		new.color = NODE_COLOR_BLACK
		*root = new
		return
	}

	t.Insert(*root, new, sentinel)

	for new != *root && isRed(new.parent) {
		if new.parent == new.parent.parent.left {
			temp = new.parent.parent.right

			if isRed(temp) {
				new.parent.color = NODE_COLOR_BLACK
				temp.color = NODE_COLOR_BLACK
				new.parent.parent.color = NODE_COLOR_RED
				new = new.parent.parent
			} else {
				if new == new.parent.right {
					new = new.parent
					t.rotateLeft(root, sentinel, new)
				}

				new.parent.color = NODE_COLOR_BLACK
				new.parent.parent.color = NODE_COLOR_RED
				t.rotateRight(root, sentinel, new.parent.parent)
			}
		} else {
			temp = new.parent.parent.left

			if isRed(temp) {
				new.parent.color = NODE_COLOR_BLACK
				temp.color = NODE_COLOR_BLACK
				new.parent.parent.color = NODE_COLOR_RED
				new = new.parent.parent
			} else {
				if new == new.parent.left {
					new = new.parent
					t.rotateRight(root, sentinel, new)
				}

				new.parent.color = NODE_COLOR_BLACK
				new.parent.parent.color = NODE_COLOR_RED
				t.rotateLeft(root, sentinel, new.parent.parent)
			}
		}
	}

	(*root).color = NODE_COLOR_BLACK
}

func (t *Tree) insertValue(temp, node, sentinel *Node) {
	var p **Node

	for {
		if node.key < temp.key {
			p = &temp.left
		} else {
			p = &temp.right
		}

		if *p == sentinel {
			break
		}

		temp = *p
	}

	*p = node
	node.parent = temp
	node.left = sentinel
	node.right = sentinel
	node.color = NODE_COLOR_RED
}

func (t *Tree) delete(tree, node *Node) {
	var (
		red                      bool
		root                     **Node
		sentinel, subst, temp, w *Node
	)

	root = &t.root
	sentinel = t.sentinel

	if node.left == sentinel {
		temp = node.right
		subst = node
	} else if node.right == sentinel {
		temp = node.left
		subst = node
	} else {
		subst = t.min(node.right, sentinel)
		temp = subst.right
	}

	if subst == *root {
		*root = temp
		temp.color = NODE_COLOR_BLACK

		node.left = nil
		node.right = nil
		node.parent = nil
		node.key = 0

		return
	}

	red = isRed(subst)

	if subst == subst.parent.left {
		subst.parent.left = temp
	} else {
		subst.parent.right = temp
	}

	if subst == node {
		temp.parent = subst.parent
	} else {
		if subst.parent == node {
			temp.parent = subst
		} else {
			temp.parent = subst.parent
		}

		subst.left = node.left
		subst.right = node.right
		subst.parent = node.parent
		subst.color = node.color

		if node == *root {
			*root = subst
		} else {
			if node == node.parent.left {
				node.parent.left = subst
			} else {
				node.parent.right = subst
			}
		}

		if subst.left != sentinel {
			subst.left.parent = subst
		}

		if subst.right != sentinel {
			subst.right.parent = subst
		}
	}

	node.left = nil
	node.right = nil
	node.parent = nil
	node.key = 0

	if red {
		return
	}

	for temp != *root && isBlack(temp) {

		if temp == temp.parent.left {
			w = temp.parent.right

			if isRed(w) {
				w.color = NODE_COLOR_BLACK
				temp.parent.color = NODE_COLOR_RED
				t.rotateLeft(root, sentinel, temp.parent)
				w = temp.parent.right
			}

			if isBlack(w.left) && isBlack(w.right) {
				w.color = NODE_COLOR_RED
				temp = temp.parent
			} else {
				if isBlack(w.right) {
					w.left.color = NODE_COLOR_BLACK
					w.color = NODE_COLOR_RED
					t.rotateRight(root, sentinel, w)
					w = temp.parent.right
				}

				w.color = temp.parent.color
				temp.parent.color = NODE_COLOR_BLACK
				w.right.color = NODE_COLOR_BLACK
				t.rotateLeft(root, sentinel, temp.parent)
				temp = *root
			}

		} else {
			w = temp.parent.left
			if isRed(w) {
				w.color = NODE_COLOR_BLACK
				temp.parent.color = NODE_COLOR_RED
				t.rotateRight(root, sentinel, temp.parent)
				w = temp.parent.left
			}

			if isBlack(w.left) && isBlack(w.right) {
				w.color = NODE_COLOR_BLACK
				temp = temp.parent

			} else {
				if isBlack(w.left) {
					w.right.color = NODE_COLOR_BLACK
					w.color = NODE_COLOR_RED
					t.rotateLeft(root, sentinel, w)
					w = temp.parent.left
				}

				w.color = temp.parent.color
				temp.parent.color = NODE_COLOR_BLACK
				w.left.color = NODE_COLOR_BLACK
				t.rotateRight(root, sentinel, temp.parent)
				temp = *root
			}
		}
	}

	temp.color = NODE_COLOR_BLACK
}

// inline
func (t *Tree) rotateLeft(root **Node, sentinel, node *Node) {
	var temp *Node

	temp = node.right
	node.right = temp.left

	if temp.left != sentinel {
		temp.left.parent = node
	}

	temp.parent = node.parent

	if node == *root {
		*root = temp
	} else if node == node.parent.left {
		node.parent.left = temp
	} else {
		node.parent.right = temp
	}

	temp.left = node
	node.parent = temp
}

// inline
func (t *Tree) rotateRight(root **Node, sentinel, node *Node) {
	var temp *Node

	temp = node.left
	node.left = temp.right

	if temp.right != sentinel {
		temp.right.parent = node
	}

	temp.parent = node.parent

	if node == *root {
		*root = temp

	} else if node == node.parent.right {
		node.parent.right = temp

	} else {
		node.parent.left = temp
	}

	temp.right = node
	node.parent = temp
}

func (t *Tree) next(node *Node) *Node {

	var root, sentinel, parent *Node

	sentinel = t.sentinel

	if node.right != sentinel {
		return t.min(node.right, sentinel)
	}

	root = t.root

	for {
		parent = node.parent

		if node == root {
			return nil
		}

		if node == parent.left {
			return parent
		}

		node = parent
	}
}
