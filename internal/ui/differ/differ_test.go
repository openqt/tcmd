package differ_test

import (
	"testing"

	"github.com/openqt/tcmd/internal/ui/differ"
)

func TestCompareText(t *testing.T) {
	out := differ.CompareText("hello", "hallo")
	if out == "" {
		t.Fatal("expected diff output")
	}
}
