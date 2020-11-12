package doc

import "strings"

// MinimarkNode ...
type MinimarkNode interface {
	Content() string
}

// Heading ...
type Heading struct {
	Level int
	Text  string
}

// Content ...
func (n *Heading) Content() string {
	return n.Text
}

// Paragraph ...
type Paragraph struct {
	Text string
}

// Content ...
func (n *Paragraph) Content() string {
	return n.Text
}

// List ...
type List struct {
	Items []MinimarkNode
}

// Content ...
func (n *List) Content() string {
	items := []string{}
	for _, item := range n.Items {
		items = append(items, item.Content())
	}
	return strings.Join(items, "\n")
}

// Item ...
type Item struct {
	Depth int
	Text  string
}

// Content ...
func (n *Item) Content() string {
	return strings.Repeat("    ", n.Depth) + " - " + n.Text
}
