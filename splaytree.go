package splaytree

import (
	"bytes"
	"fmt"
)

const (
	Left = iota
	Right
	ChildNum
)

const Indent = 2

/*
  self < other, return negative value
  self > other, return positive value
  self == other,return zero value
*/
type Comparable interface {
	Compare(other Comparable) int
}

type PositiveInfinity struct {
}

func (pi PositiveInfinity) Compare(other Comparable) int {
	return 1
}

var pi = PositiveInfinity{}

type node struct {
	links [ChildNum]*node
	key   Comparable
	//val   interface{}
}

type SplayTree struct {
	root *node
}

func New() *SplayTree {
	return &SplayTree{}
}

func rotate(root **node, dir int) {
	var (
		x *node
		y *node
		b *node
	)
	/*
	       y           x
	      / \         / \
	     x   c  <=>  a   y
	    / \             / \
	   a   b           b   c
	*/
	y = *root
	if y == nil {
		return
	}
	x = y.links[dir]
	if x == nil {
		return
	}
	b = x.links[dir^1]
	*root = x
	x.links[dir^1] = y
	y.links[dir] = b
}

func (n *node) printHelper(depth int, buf *bytes.Buffer) {
	if n == nil {
		return
	}

	n.links[Left].printHelper(depth+1, buf)

	for i := 0; i < depth*Indent; i++ {
		buf.WriteRune(' ')
	}
	buf.WriteString(fmt.Sprintf("%v [%p]\n", n.key, n))

	n.links[Right].printHelper(depth+1, buf)
}

func (t *SplayTree) String() string {
	var buf bytes.Buffer
	t.root.printHelper(0, &buf)
	return buf.String()
}

func link(hook []**node, dir int, node *node) {
	*hook[dir] = node
	node.links[dir^1] = nil
	hook[dir] = &node.links[dir^1]
}

//自顶向下伸展
func splay(root **node, target Comparable) {
	var (
		w           *node
		child       *node
		grandchild  *node
		top         [ChildNum]*node  //记录伸展过程中，切分出来的子树L,R 的根
		hook        [ChildNum]**node //记录L,R 的挂载点
		d           int
		dChild      int
		dGrandchild int
		buf         bytes.Buffer
	)
	w = *root
	if w == nil {
		return
	}
	for d := 0; d < ChildNum; d++ {
		top[d] = nil
		hook[d] = &top[d]
	}
	for {
		ret := w.key.Compare(target)
		if ret == 0 {
			break
		}
		if ret < 0 {
			dChild = Right
		} else {
			dChild = Left
		}
		child = w.links[dChild]
		if child == nil {
			break
		}
		ret2 := child.key.Compare(target)
		if ret2 < 0 {
			dGrandchild = Right
		} else {
			dGrandchild = Left
		}
		grandchild = child.links[dGrandchild]

		top[0].printHelper(0, &buf)
		fmt.Printf("\ntree top[0]:\n%s\n", buf.String())
		buf.Reset()
		fmt.Println("-----------")
		w.printHelper(0, &buf)
		fmt.Printf("\ntree w:\n%s\n", buf.String())
		buf.Reset()
		fmt.Println("-----------")
		top[1].printHelper(0, &buf)
		fmt.Printf("\ntree top[1]:\n%s\n", buf.String())
		buf.Reset()
		fmt.Println("===============")

		if grandchild == nil || child.key.Compare(target) == 0 {
			/* zig case */
			link(hook[:], dChild^1, w)
			w = child
			break
		} else if dChild == dGrandchild {
			/* zig-zig case */
			rotate(&w, dChild)
			link(hook[:], dChild^1, child)
			w = grandchild
		} else {
			/* zig-zag case */
			link(hook[:], dChild^1, w)
			link(hook[:], dChild, child)
			w = grandchild
		}
	}

	/*
		完成伸展之后,重新组合L,w,R,
	*/
	for d = 0; d < ChildNum; d++ {
		*hook[d] = w.links[d]
		w.links[d] = top[d]
	}

	/*更新w为新的树根*/
	*root = w
}

func (t *SplayTree) Exist(target Comparable) bool {
	splay(&t.root, target)
	return t.root != nil && t.root.key.Compare(target) == 0
}

func (t *SplayTree) Insert(data Comparable) {
	insert(&t.root, data)
}

func insert(root **node, data Comparable) {
	var (
		e *node
		w *node
		d int
	)
	splay(root, data)
	w = *root

	/*已存在data*/
	if w != nil && w.key.Compare(data) == 0 {
		return
	}
	e = &node{key: data}

	if w != nil {
		if (*root).key.Compare(data) > 0 {
			d = Right
		} else {
			d = Left
		}
		/*
			                   e                         e
					 /   \                     /   \
				        nil   w      =>           a     w
					     / \                       /  \
					     a  b                    nil    b
					              or
				          e                            e
					/   \                        /   \
					w    nil     =>             w     b
				       / \                         / \
				      a   b                       a  nil

		*/
		/* 分割树，新节点变成树根 */
		/* 围绕目标作伸展之后, 节点w 变成目标节点的前驱或后继 */
		e.links[d] = w
		e.links[d^1] = w.links[d^1]
		w.links[d^1] = nil
	}
	*root = e
}

func (t *SplayTree) Delete(data Comparable) {
	delete(&t.root, data)
}

func delete(root **node, target Comparable) {
	var (
		left  *node
		right *node
	)
	splay(root, target)
	//伸展之后，找到了待删除的目标
	if *root != nil && (*root).key.Compare(target) == 0 {
		left = (*root).links[Left]
		right = (*root).links[Right]
		if left == nil {
			*root = right
		} else {
			/* 伸展左子树的最大元素节点到顶部*/
			splay(&left, pi)
			/*左子树的右链接挂载右子树*/
			left.links[Right] = right
			/*更新左子树为新的根节点*/
			*root = left
		}
	}
}
