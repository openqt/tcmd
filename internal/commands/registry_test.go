package commands_test

import (
	"testing"

	"github.com/openqt/tcmd/internal/commands"
)

func TestRegistryExecuteNotFound(t *testing.T) {
	t.Parallel()
	r := commands.NewRegistry()
	ctx := &commands.Context{
		SetStatus: func(s string) { /* noop */ },
	}
	if got := r.Execute(ctx, "cm_DoesNotExist", nil); got != commands.ResultNotFound {
		t.Fatalf("got %v want ResultNotFound", got)
	}
}

func TestRegistryNormalizeCommand(t *testing.T) {
	t.Parallel()
	r := commands.NewRegistry()
	called := false
	r.Register(commands.Def{
		Name: "Exit",
		Handler: func(_ *commands.Context, _ []string) commands.Result {
			called = true
			return commands.ResultSuccess
		},
	})
	if got := r.Execute(&commands.Context{}, "cm_Exit", nil); got != commands.ResultSuccess {
		t.Fatalf("got %v want success", got)
	}
	if !called {
		t.Fatal("handler not called")
	}
}
