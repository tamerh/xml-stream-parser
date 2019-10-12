package xmlparser

import (
	"github.com/tamerh/xpath"
)

// CreateXPathNavigator creates a new xpath.NodeNavigator for the specified html.Node.
func (x *XMLParser) CreateXPathNavigator(top *XMLElement) *XmlNodeNavigator {
	return &XmlNodeNavigator{curr: top, root: top, attr: -1}
}

// Compile the given xpath expression
func (x *XMLParser) CompileXpath(expr string) (*xpath.Expr, error) {

	exp, err := xpath.Compile(expr)
	if err != nil {
		return nil, err
	}
	return exp, nil

}

// CreateXPathNavigator creates a new xpath.NodeNavigator for the specified html.Node.
func createXPathNavigator(top *XMLElement) *XmlNodeNavigator {
	return &XmlNodeNavigator{curr: top, root: top, attr: -1}
}

type XmlNodeNavigator struct {
	root, curr *XMLElement
	attr       int
}

// Find searches the Node that matches by the specified XPath expr.
func find(top *XMLElement, expr string) ([]*XMLElement, error) {
	exp, err := xpath.Compile(expr)
	if err != nil {
		return []*XMLElement{}, err
	}
	t := exp.Select(createXPathNavigator(top))
	var elems []*XMLElement
	for t.MoveNext() {
		elems = append(elems, t.Current().(*XmlNodeNavigator).curr)
	}
	return elems, nil
}

// FindOne searches the Node that matches by the specified XPath expr,
// and returns first element of matched.
func findOne(top *XMLElement, expr string) (*XMLElement, error) {
	exp, err := xpath.Compile(expr)
	if err != nil {
		return nil, err
	}
	t := exp.Select(createXPathNavigator(top))
	var elem *XMLElement
	if t.MoveNext() {
		elem = t.Current().(*XmlNodeNavigator).curr //getCurrentNode(t)
	}
	return elem, nil
}

func (x *XmlNodeNavigator) Current() *XMLElement {
	return x.curr
}

func (x *XmlNodeNavigator) NodeType() xpath.NodeType {

	if x.curr == x.root {
		return xpath.RootNode
	}
	if x.attr != -1 {
		return xpath.AttributeNode
	}
	return xpath.ElementNode
}

func (x *XmlNodeNavigator) LocalName() string {
	if x.attr != -1 {
		return x.curr.attrs[x.attr].name
	}

	return x.curr.localName

}

func (x *XmlNodeNavigator) Prefix() string {

	return x.curr.prefix

}

func (x *XmlNodeNavigator) Value() string {

	if x.attr != -1 {
		return x.curr.attrs[x.attr].value
	}
	return x.curr.InnerText

}

func (x *XmlNodeNavigator) Copy() xpath.NodeNavigator {
	n := *x
	return &n
}

func (x *XmlNodeNavigator) MoveToRoot() {
	x.curr = x.root
}

func (x *XmlNodeNavigator) MoveToParent() bool {
	if x.attr != -1 {
		x.attr = -1
		return true
	} else if node := x.curr.parent; node != nil {
		x.curr = node
		return true
	}
	return false
}

func (x *XmlNodeNavigator) MoveToNextAttribute() bool {
	if x.attr >= len(x.curr.attrs)-1 {
		return false
	}
	x.attr++
	return true
}

func (x *XmlNodeNavigator) MoveToChild() bool {
	if node := x.curr.FirstChild(); node != nil {
		x.curr = node
		return true
	}
	return false
}

func (x *XmlNodeNavigator) MoveToFirst() bool {
	if x.curr.parent != nil {
		node := x.curr.parent.FirstChild()
		if node != nil {
			x.curr = node
			return true
		}
	}
	return false
}

func (x *XmlNodeNavigator) MoveToPrevious() bool {
	node := x.curr.PrevSibling()
	if node != nil {
		x.curr = node
		return true
	}
	return false
}

func (x *XmlNodeNavigator) MoveToNext() bool {
	node := x.curr.NextSibling()
	if node != nil {
		x.curr = node
		return true
	}
	return false
}

func (x *XmlNodeNavigator) String() string {
	return x.Value()
}

func (x *XmlNodeNavigator) MoveTo(other xpath.NodeNavigator) bool {
	node, ok := other.(*XmlNodeNavigator)
	if !ok || node.root != x.root {
		return false
	}

	x.curr = node.curr
	x.attr = node.attr
	return true
}
