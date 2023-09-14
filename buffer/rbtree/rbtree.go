package rbtree

import "unsafe"

type Node struct {
	left   *Node
	right  *Node
	parent *Node

	color bool

	key uint32
	val unsafe.Pointer
}

const (
	RED   = false
	BLACK = true
)

type Rbtree struct {
	NIL  *Node
	root *Node
}

func NewRbtree() *Rbtree { return new(Rbtree).Init() }

func (t *Rbtree) Init() *Rbtree {
	node := &Node{color: BLACK}
	return &Rbtree{
		NIL:  node,
		root: node,
	}
}

func (t *Rbtree) Insert(key uint32, val unsafe.Pointer) {
	t.insert(t.newDefaultNode(key, val))
}

func (t *Rbtree) Lookup(key uint32) unsafe.Pointer {
	n := t.search(t.newDefaultNode(key, nil))
	if n == nil {
		return nil
	}

	return n.val
}

func (t *Rbtree) Delete(key uint32) {
	t.delete(t.newDefaultNode(key, nil))
}

func (t *Rbtree) newDefaultNode(key uint32, val unsafe.Pointer) *Node {
	return &Node{
		left:   t.NIL,
		right:  t.NIL,
		parent: t.NIL,
		color:  RED,
		key:    key,
		val:    val,
	}
}

func (t *Rbtree) leftRotate(x *Node) {
	if x.right == t.NIL {
		return
	}

	y := x.right
	x.right = y.left
	if y.left != t.NIL {
		y.left.parent = x
	}
	y.parent = x.parent

	if x.parent == t.NIL {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.left = x
	x.parent = y
}

func (t *Rbtree) rightRotate(x *Node) {
	if x.left == t.NIL {
		return
	}

	y := x.left
	x.left = y.right
	if y.right != t.NIL {
		y.right.parent = x
	}
	y.parent = x.parent

	if x.parent == t.NIL {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.right = x
	x.parent = y
}

func (t *Rbtree) insert(z *Node) *Node {
	x := t.root
	y := t.NIL

	for x != t.NIL {
		y = x
		if z.key < x.key {
			x = x.left
		} else if x.key < z.key {
			x = x.right
		} else {
			return x
		}
	}

	z.parent = y
	if y == t.NIL {
		t.root = z
	} else if z.key < y.key {
		y.left = z
	} else {
		y.right = z
	}

	t.insertFixup(z)
	return z
}

func (t *Rbtree) insertFixup(z *Node) {
	for z.parent.color == RED {
		if z.parent == z.parent.parent.left {
			y := z.parent.parent.right
			if y.color == RED {
				z.parent.color = BLACK
				y.color = BLACK
				z.parent.parent.color = RED
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					z = z.parent
					t.leftRotate(z)
				}

				z.parent.color = BLACK
				z.parent.parent.color = RED
				t.rightRotate(z.parent.parent)
			}
		} else {
			y := z.parent.parent.left
			if y.color == RED {
				z.parent.color = BLACK
				y.color = BLACK
				z.parent.parent.color = RED
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					t.rightRotate(z)
				}
				z.parent.color = BLACK
				z.parent.parent.color = RED
				t.leftRotate(z.parent.parent)
			}
		}
	}
	t.root.color = BLACK
}

func (t *Rbtree) min(x *Node) *Node {
	if x == t.NIL {
		return t.NIL
	}

	for x.left != t.NIL {
		x = x.left
	}

	return x
}

func (t *Rbtree) max(x *Node) *Node {
	if x == t.NIL {
		return t.NIL
	}

	for x.right != t.NIL {
		x = x.right
	}

	return x
}

func (t *Rbtree) search(x *Node) *Node {
	p := t.root

	for p != t.NIL {
		if p.key < x.key {
			p = p.right
		} else if x.key < p.key {
			p = p.left
		} else {
			break
		}
	}

	return p
}

func (t *Rbtree) successor(x *Node) *Node {
	if x == t.NIL {
		return t.NIL
	}

	if x.right != t.NIL {
		return t.min(x.right)
	}

	y := x.parent
	for y != t.NIL && x == y.right {
		x = y
		y = y.parent
	}
	return y
}

func (t *Rbtree) delete(key *Node) *Node {
	z := t.search(key)

	if z == t.NIL {
		return t.NIL
	}
	ret := &Node{color: z.color}

	var y *Node
	var x *Node

	if z.left == t.NIL || z.right == t.NIL {
		y = z
	} else {
		y = t.successor(z)
	}

	if y.left != t.NIL {
		x = y.left
	} else {
		x = y.right
	}

	x.parent = y.parent

	if y.parent == t.NIL {
		t.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}

	if y != z {
		z.key = y.key
	}

	if y.color == BLACK {
		t.deleteFixup(x)
	}

	return ret
}

func (t *Rbtree) deleteFixup(x *Node) {
	for x != t.root && x.color == BLACK {
		if x == x.parent.left {
			w := x.parent.right
			if w.color == RED {
				w.color = BLACK
				x.parent.color = RED
				t.leftRotate(x.parent)
				w = x.parent.right
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				x = x.parent
			} else {
				if w.right.color == BLACK {
					w.left.color = BLACK
					w.color = RED
					t.rightRotate(w)
					w = x.parent.right
				}
				w.color = x.parent.color
				x.parent.color = BLACK
				w.right.color = BLACK
				t.leftRotate(x.parent)

				// this is to exit while loop
				x = t.root
			}
		} else {
			w := x.parent.left
			if w.color == RED {
				w.color = BLACK
				x.parent.color = RED
				t.rightRotate(x.parent)
				w = x.parent.left
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				x = x.parent
			} else {
				if w.left.color == BLACK {
					w.right.color = BLACK
					w.color = RED
					t.leftRotate(w)
					w = x.parent.left
				}
				w.color = x.parent.color
				x.parent.color = BLACK
				w.left.color = BLACK
				t.rightRotate(x.parent)
				x = t.root
			}
		}
	}
	x.color = BLACK
}
