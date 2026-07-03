package panel

// Notebook manages tabs for one panel side.
type Notebook struct {
	Side   Side
	Tabs   []*Panel
	Active int
}

// NewNotebook creates a notebook with one tab.
func NewNotebook(side Side, path string) *Notebook {
	return &Notebook{
		Side:   side,
		Tabs:   []*Panel{New(side, path)},
		Active: 0,
	}
}

// Current returns the active tab panel.
func (n *Notebook) Current() *Panel {
	if n == nil || len(n.Tabs) == 0 {
		return New(n.Side, "/")
	}
	if n.Active < 0 || n.Active >= len(n.Tabs) {
		n.Active = 0
	}
	return n.Tabs[n.Active]
}

// NewTab opens a new tab at path.
func (n *Notebook) NewTab(path string) {
	n.Tabs = append(n.Tabs, New(n.Side, path))
	n.Active = len(n.Tabs) - 1
}

// CloseTab closes the active tab.
func (n *Notebook) CloseTab() bool {
	if len(n.Tabs) <= 1 {
		return false
	}
	n.Tabs = append(n.Tabs[:n.Active], n.Tabs[n.Active+1:]...)
	if n.Active >= len(n.Tabs) {
		n.Active = len(n.Tabs) - 1
	}
	return true
}

// NextTab advances active tab.
func (n *Notebook) NextTab() {
	if len(n.Tabs) == 0 {
		return
	}
	n.Active = (n.Active + 1) % len(n.Tabs)
}

// PrevTab moves to previous tab.
func (n *Notebook) PrevTab() {
	if len(n.Tabs) == 0 {
		return
	}
	n.Active--
	if n.Active < 0 {
		n.Active = len(n.Tabs) - 1
	}
}

// ActivateTab switches to tab index (0-based).
func (n *Notebook) ActivateTab(index int) {
	if index < 0 || index >= len(n.Tabs) {
		return
	}
	n.Active = index
}

// TabTitles returns tab labels.
func (n *Notebook) TabTitles() []string {
	titles := make([]string, len(n.Tabs))
	for i, t := range n.Tabs {
		titles[i] = t.Title()
	}
	return titles
}
