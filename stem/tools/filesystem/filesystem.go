// Package filesystem provides file operations. Read-only in Phase 1.
package filesystem

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ListTool lists files in a directory.
type ListTool struct{}

func NewList() *ListTool { return &ListTool{} }
func (t *ListTool) Name() string { return "filesystem_list" }
func (t *ListTool) Run(ctx context.Context, args map[string]any) (string, error) {
	path := getPath(args)
	if path == "" {
		path = "."
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("list: %w", err)
	}
	var b strings.Builder
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() {
			name += "/"
		}
		b.WriteString(name)
		b.WriteString("\n")
	}
	return b.String(), nil
}

// ReadTool reads a file's contents.
type ReadTool struct{}

func NewRead() *ReadTool { return &ReadTool{} }
func (t *ReadTool) Name() string { return "filesystem_read" }
func (t *ReadTool) Run(ctx context.Context, args map[string]any) (string, error) {
	path := getPath(args)
	if path == "" {
		return "", fmt.Errorf("read: path is required")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read: %w", err)
	}
	return string(data), nil
}

// SearchTool finds files by name pattern.
type SearchTool struct{}

func NewSearch() *SearchTool { return &SearchTool{} }
func (t *SearchTool) Name() string { return "filesystem_search" }
func (t *SearchTool) Run(ctx context.Context, args map[string]any) (string, error) {
	pattern := ""
	if v, ok := args["pattern"]; ok {
		pattern, _ = v.(string)
	}
	if pattern == "" {
		// Try alternative field names
		if v, ok := args["query"]; ok {
			pattern, _ = v.(string)
		}
	}
	if pattern == "" {
		return "", fmt.Errorf("search: pattern is required")
	}

	root := "."
	if v, ok := args["path"]; ok {
		root, _ = v.(string)
	}

	var results []string
	filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.Contains(info.Name(), pattern) {
			results = append(results, p)
		}
		return nil
	})
	if len(results) == 0 {
		return "No files found matching \"" + pattern + "\"", nil
	}
	return strings.Join(results, "\n"), nil
}

func getPath(args map[string]any) string {
	for _, key := range []string{"path", "dir", "directory", "文件夹", "目录"} {
		if v, ok := args[key]; ok {
			s, _ := v.(string)
			if s != "" {
				return s
			}
		}
	}
	return ""
}
