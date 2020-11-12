package doc

import "encoding/json"

// MarshalJSON ...
func (n *Heading) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type  string `json:"node.type"`
		Level int    `json:"level"`
		Text  string `json:"text"`
	}{
		"heading", n.Level, n.Text,
	})
}

// MarshalJSON ...
func (n *Paragraph) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type string `json:"node.type"`
		Text string `json:"text"`
	}{
		"paragraph", n.Text,
	})
}

// MarshalJSON ...
func (n *List) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type  string  `json:"node.type"`
		Items []*Item `json:"items"`
	}{
		"list", n.Items,
	})
}

// MarshalJSON ...
func (n *Item) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type  string `json:"node.type"`
		Depth int    `json:"depth"`
		Text  string `json:"text"`
	}{
		"list.item", n.Depth, n.Text,
	})
}
