package hotkeys

import tea "github.com/charmbracelet/bubbletea"

// Dispatcher resolves key chords for a shortcut context.
type Dispatcher struct {
	bindings map[Context][]Binding
}

// NewDispatcher builds a dispatcher with the provided bindings per context.
func NewDispatcher(bindings map[Context][]Binding) *Dispatcher {
	if bindings == nil {
		bindings = map[Context][]Binding{
			ContextMain: DefaultMainBindings(),
		}
	}
	return &Dispatcher{bindings: bindings}
}

// Lookup returns the command bound to a key in the given context.
func (d *Dispatcher) Lookup(ctx Context, key string) (Match, bool) {
	for _, binding := range d.bindings[ctx] {
		if binding.Key == key {
			return Match{
				Command: binding.Command,
				Params:  append([]string(nil), binding.Params...),
			}, true
		}
	}
	return Match{}, false
}

// LookupMsg resolves a Bubble Tea key message.
func (d *Dispatcher) LookupMsg(ctx Context, msg tea.KeyMsg) (Match, bool) {
	return d.Lookup(ctx, NormalizeKey(msg))
}

// Bindings returns bindings for a context.
func (d *Dispatcher) Bindings(ctx Context) []Binding {
	out := make([]Binding, len(d.bindings[ctx]))
	copy(out, d.bindings[ctx])
	return out
}
