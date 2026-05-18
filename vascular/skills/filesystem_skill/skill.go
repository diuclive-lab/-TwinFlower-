// Package filesystem_skill orchestrates filesystem operations with procedural intelligence.
// Not just tool wrappers — it handles path inference, ambiguity, and result shaping.
package filesystem_skill

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

)

// FileIntent describes what the user wants to do with the filesystem.
type FileIntent int

const (
	IntentList   FileIntent = iota // list directory contents
	IntentRead                     // read a specific file
	IntentSearch                   // search for files
	IntentLargest                  // find the largest file
	IntentStatus                   // check file/dir status
	IntentAmbiguous                // can't determine intent
)

func (i FileIntent) String() string {
	switch i {
	case IntentList:
		return "list"
	case IntentRead:
		return "read"
	case IntentSearch:
		return "search"
	case IntentLargest:
		return "largest"
	case IntentStatus:
		return "status"
	case IntentAmbiguous:
		return "ambiguous"
	default:
		return "unknown"
	}
}

// ToolRunner is the interface a tool must satisfy (matches engine.Tool).
type ToolRunner interface {
	Name() string
	Run(ctx context.Context, args map[string]any) (string, error)
}

// Skill is the filesystem orchestration skill.
type Skill struct {
	tools     map[string]ToolRunner
}

// New creates a filesystem skill with references to the engine's tools.
func New(tools map[string]ToolRunner) *Skill {
	return &Skill{tools: tools}
}

// Handle processes a filesystem request through the full skill pipeline.
func (s *Skill) Handle(ctx context.Context, input string) (string, error) {
	// 1. Parse intent from natural language
	intent := s.parseIntent(input)

	// 2. Extract path
	path := s.extractPath(input)
	if path == "" {
		path = "." // default to current directory
	}

	// 3. Handle ambiguity
	if intent == IntentAmbiguous {
		return "", fmt.Errorf("ambiguous_request: %s", input)
	}

	// 4. Existence check (for specific file/dir paths)
	if path != "." {
		info, err := os.Stat(path)
		if err != nil {
			// Path doesn't exist — try prefix matching
			return s.handleMissingPath(ctx, path)
		}
		if intent == IntentStatus || intent == IntentAmbiguous {
			if info.IsDir() {
				intent = IntentList
			} else {
				intent = IntentRead
			}
		}
	}

	// 5. Execute
	switch intent {
	case IntentList:
		return s.executeList(ctx, path)
	case IntentRead:
		return s.executeRead(ctx, path)
	case IntentSearch:
		return s.executeSearch(ctx, path, input)
	case IntentLargest:
		return s.executeLargest(ctx, path)
	default:
		return s.executeList(ctx, path)
	}
}

// ── Intent Parsing ────────────────────────────────────────────────────────

func (s *Skill) parseIntent(input string) FileIntent {
	lower := strings.ToLower(strings.TrimSpace(input))
	// Check for "最大的" / "largest" / "最大"
	if strings.Contains(lower, "最大的") || strings.Contains(lower, " largest") ||
		strings.Contains(lower, "最大") || strings.Contains(lower, "大文件") ||
		strings.Contains(lower, "biggest") || strings.Contains(lower, "largest") {
		return IntentLargest
	}
	// Check for "找" / "搜索" / "search"
	if strings.Contains(lower, "找") || strings.Contains(lower, "搜索") ||
		strings.Contains(lower, "search") || strings.Contains(lower, "查找") {
		return IntentSearch
	}
	// Check for "读" / "看" / "打开" / "read"
	if strings.Contains(lower, "读") || strings.Contains(lower, "打开") ||
		strings.Contains(lower, "read ") || strings.Contains(lower, "cat ") {
		return IntentRead
	}
	// Check for "列出" / "列表" / "list" / "目录" / "文件夹"
	if strings.Contains(lower, "列") || strings.Contains(lower, "list") ||
		strings.Contains(lower, "目录") || strings.Contains(lower, "文件夹") ||
		strings.Contains(lower, "dir") || strings.Contains(lower, "ls") {
		return IntentList
	}
	return IntentList // default to list
}

// ── Path Extraction ───────────────────────────────────────────────────────

func (s *Skill) extractPath(input string) string {
	// Known directory markers in Chinese and English
	markers := []string{"目录", "文件夹", "路径", "文件", "directory", "folder", "path", "file"}

	words := strings.Fields(input)
	for i, w := range words {
		for _, m := range markers {
			if strings.Contains(w, m) && i+1 < len(words) {
				candidate := words[i+1]
				// Check if it's a valid path
				if !isPath(candidate) && i+2 < len(words) {
					candidate = words[i+1] + words[i+2]
				}
				if isPath(candidate) || fileExists(candidate) {
					return candidate
				}
			}
		}
	}

	// Check for absolute paths or paths with /
	for _, w := range words {
		if isPath(w) && fileExists(w) {
			return w
		}
	}

	return ""
}

func isPath(s string) bool {
	return strings.Contains(s, "/") || strings.Contains(s, "\\") || strings.HasPrefix(s, ".")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ── Missing Path Recovery ─────────────────────────────────────────────────

func (s *Skill) handleMissingPath(ctx context.Context, path string) (string, error) {
	// Try prefix match: search for files starting with the given name
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	if dir == "." {
		dir = "."
	}

	var results []string
	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.HasPrefix(strings.ToLower(info.Name()), strings.ToLower(base)) {
			results = append(results, p)
		}
		return nil
	})

	if len(results) == 0 {
		return fmt.Sprintf("路径 %q 不存在，且未找到匹配的文件", path), nil
	}
	if len(results) == 1 {
		return s.executeRead(ctx, results[0])
	}
	// Multiple matches
	msg := fmt.Sprintf("路径 %q 不存在，但找到多个匹配：\n", path)
	for _, r := range results {
		msg += "  " + r + "\n"
	}
	return msg, nil
}

// ── Tool Execution ────────────────────────────────────────────────────────

func (s *Skill) executeList(ctx context.Context, path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("路径不存在: %s", path)
	}
	if !info.IsDir() {
		// It's a file, not a directory — read it
		return s.executeRead(ctx, path)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("读取目录失败: %w", err)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "📁 %s （共 %d 项）\n\n", path, len(entries))
	for _, e := range entries {
		fi, _ := e.Info()
		name := e.Name()
		if e.IsDir() {
			name += "/"
		}
		if fi != nil {
			size := fi.Size()
			if size > 1024*1024 {
				fmt.Fprintf(&b, "  %s  (%.1f MB)\n", name, float64(size)/(1024*1024))
			} else if size > 1024 {
				fmt.Fprintf(&b, "  %s  (%.1f KB)\n", name, float64(size)/1024)
			} else {
				fmt.Fprintf(&b, "  %s\n", name)
			}
		} else {
			fmt.Fprintf(&b, "  %s\n", name)
		}
	}
	return b.String(), nil
}

func (s *Skill) executeRead(ctx context.Context, path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}
	content := string(data)
	if len(content) > 2000 {
		content = content[:2000] + "\n...（文件过长，仅显示前 2000 字符）"
	}
	return fmt.Sprintf("📄 %s\n\n%s", path, content), nil
}

func (s *Skill) executeSearch(ctx context.Context, path string, input string) (string, error) {
	pattern := ""
	lower := strings.ToLower(input)
	for _, keyword := range []string{"找", "搜索", "search", "查找", "find"} {
		idx := strings.Index(lower, keyword)
		if idx >= 0 {
			after := input[idx+len([]rune(keyword)):]
			// Take the first meaningful word after the keyword
			words := strings.Fields(after)
			if len(words) > 0 {
				pattern = strings.TrimSpace(words[0])
				break
			}
		}
	}
	if pattern == "" {
		pattern = "*.go" // default search pattern
	}

	var results []string
	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(pattern)) {
			results = append(results, p)
		}
		return nil
	})

	if len(results) == 0 {
		return fmt.Sprintf("未找到匹配 %q 的文件", pattern), nil
	}
	msg := fmt.Sprintf("🔍 找到 %d 个匹配 %q 的文件：\n\n", len(results), pattern)
	for _, r := range results {
		msg += "  " + r + "\n"
	}
	return msg, nil
}

func (s *Skill) executeLargest(ctx context.Context, path string) (string, error) {
	var files []struct {
		path string
		size int64
	}

	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		files = append(files, struct {
			path string
			size int64
		}{p, info.Size()})
		return nil
	})

	if len(files) == 0 {
		return fmt.Sprintf("目录 %s 中没有文件", path), nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].size > files[j].size
	})

	// Top 5 largest
	n := 5
	if len(files) < n {
		n = len(files)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "📊 %s 中最大的 %d 个文件：\n\n", path, n)
	for i := 0; i < n; i++ {
		f := files[i]
		var sizeStr string
		if f.size > 1024*1024*1024 {
			sizeStr = fmt.Sprintf("%.2f GB", float64(f.size)/(1024*1024*1024))
		} else if f.size > 1024*1024 {
			sizeStr = fmt.Sprintf("%.1f MB", float64(f.size)/(1024*1024))
		} else if f.size > 1024 {
			sizeStr = fmt.Sprintf("%.1f KB", float64(f.size)/1024)
		} else {
			sizeStr = fmt.Sprintf("%d B", f.size)
		}
		fmt.Fprintf(&b, "  %d. %s  (%s)\n", i+1, f.path, sizeStr)
	}
	return b.String(), nil
}
