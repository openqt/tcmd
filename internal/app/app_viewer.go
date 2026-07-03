package app

import (
	"github.com/openqt/tcmd/internal/commands/builtin"
	"github.com/openqt/tcmd/internal/ui/viewer"
)

func (m *Model) SetViewerMode(path string, mode viewer.Mode) {
	content, err := viewer.Render(path, mode, 8000)
	if err != nil {
		m.status = err.Error()
		return
	}
	m.viewer = content
	m.viewerMode = mode
	m.status = "Viewer"
}

func (m *Model) OpenEditor(path string) {
	content, err := builtin.RenderEditor(path)
	if err != nil {
		m.status = err.Error()
		return
	}
	m.editorPath = path
	m.editor = content
	m.status = "Editor: " + path
}

func (m *Model) OpenDiffer(left, right string) {
	content, err := builtin.RenderDiffer(left, right)
	if err != nil {
		m.status = err.Error()
		return
	}
	m.differText = content
	m.status = "Differ"
}
