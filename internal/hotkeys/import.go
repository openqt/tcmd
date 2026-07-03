package hotkeys

import (
	"os"

	"gopkg.in/yaml.v3"
)

// File represents importable shortcut bindings.
type File struct {
	Context  string    `yaml:"context"`
	Bindings []Binding `yaml:"bindings"`
}

// LoadYAML imports bindings from a YAML file.
func LoadYAML(path string) ([]File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var files []File
	if err := yaml.Unmarshal(data, &files); err != nil {
		var single File
		if err2 := yaml.Unmarshal(data, &single); err2 != nil {
			return nil, err
		}
		return []File{single}, nil
	}
	return files, nil
}

// SaveYAML exports bindings to YAML.
func SaveYAML(path string, files []File) error {
	data, err := yaml.Marshal(files)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// MergeBindings overlays imported bindings onto defaults.
func MergeBindings(defaults []Binding, imported []Binding) []Binding {
	out := append([]Binding(nil), defaults...)
	seen := map[string]bool{}
	for _, b := range out {
		seen[b.Key] = true
	}
	for _, b := range imported {
		if !seen[b.Key] {
			out = append(out, b)
			seen[b.Key] = true
		}
	}
	return out
}
