// TwinFlower — Phase 1: Minimal Rooting
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"twinflower/root/cognition"
	"twinflower/root/cognition/preferences"
	"twinflower/root/providers"
	"twinflower/runtime/engine"
	fs_skill "twinflower/vascular/skills/filesystem_skill"
	search_skill "twinflower/vascular/skills/search_skill"
	"twinflower/stem/tools/filesystem"
	"twinflower/stem/tools/search"
	"twinflower/stem/tools/translate"
	"twinflower/stem/tools/weather"
)

func main() {
	ctx := context.Background()

	provider := providers.NewLocalProvider("http://127.0.0.1:8090", "qwen3.6-27B", "local-dev")

	cp := cognition.QwenDense()

	prefs := preferences.NewStore("runtime/preferences.jsonl")

	e := engine.New(provider, cp)
	e.SetPreferences(prefs)

	e.RegisterTool(weather.New())
	e.RegisterTool(filesystem.NewList())
	e.RegisterTool(filesystem.NewRead())
	e.RegisterTool(filesystem.NewSearch())
	e.RegisterTool(search.New())
	e.RegisterTool(translate.New())

	// Register procedural filesystem skill
	fsTools := map[string]fs_skill.ToolRunner{
		"filesystem_list":   filesystem.NewList(),
		"filesystem_read":   filesystem.NewRead(),
		"filesystem_search": filesystem.NewSearch(),
	}
	e.SetFilesystemSkill(fs_skill.New(fsTools))
	e.SetSearchSkill(search_skill.New(search.New()))

	input := strings.Join(os.Args[1:], " ")
	if input == "" {
		input = "你好"
	}

	response, err := e.Handle(ctx, input)
	if err != nil {
		log.Fatalf("handle: %v", err)
	}

	fmt.Println(response)
}
