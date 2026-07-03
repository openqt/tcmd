package commands

import (
	"fmt"
	"strings"
)

// Handler executes a Double Commander internal command.
type Handler func(ctx *Context, params []string) Result

// Def describes a registered cm_* command.
type Def struct {
	Name     string
	Category string
	Handler  Handler
}

// Registry stores cm_* command handlers.
type Registry struct {
	commands map[string]Def
}

// NewRegistry returns an empty command registry.
func NewRegistry() *Registry {
	return &Registry{commands: make(map[string]Def)}
}

// Register adds or replaces a command definition.
func (r *Registry) Register(def Def) {
	name := normalizeCommand(def.Name)
	r.commands[name] = Def{
		Name:     name,
		Category: def.Category,
		Handler:  def.Handler,
	}
}

// Execute runs a command by name.
func (r *Registry) Execute(ctx *Context, command string, params []string) Result {
	name := normalizeCommand(command)
	def, ok := r.commands[name]
	if !ok {
		if ctx != nil && ctx.SetStatus != nil {
			ctx.SetStatus(fmt.Sprintf("[dc-tui] %s: not implemented", name))
		}
		return ResultNotFound
	}
	if def.Handler == nil {
		if ctx != nil && ctx.SetStatus != nil {
			ctx.SetStatus(fmt.Sprintf("[dc-tui] %s: handler missing", name))
		}
		return ResultDisabled
	}
	return def.Handler(ctx, params)
}

// Has reports whether a command is registered.
func (r *Registry) Has(command string) bool {
	_, ok := r.commands[normalizeCommand(command)]
	return ok
}

// Names returns registered command names sorted lexicographically.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.commands))
	for name := range r.commands {
		names = append(names, name)
	}
	return names
}

func normalizeCommand(command string) string {
	command = strings.TrimSpace(command)
	if command == "" {
		return ""
	}
	if !strings.HasPrefix(command, "cm_") {
		command = "cm_" + command
	}
	return command
}
