// TwinFlower — Phase 1: Minimal Rooting
//
// Usage:
//
//	go run . "北京天气"
//	go run . "hello翻译成日语"
//	go run . "列出当前目录"
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"twinflower/root/providers"
	"twinflower/runtime/engine"
	"twinflower/stem/tools/weather"
)

func main() {
	ctx := context.Background()

	// ── Model ───────────────────────────────────────────────────────────
	provider := providers.NewLocalProvider("http://127.0.0.1:8090", "qwen3.6-27B", "local-dev")

	// ── Engine ──────────────────────────────────────────────────────────
	e := engine.New(provider)

	// Register tools
	e.RegisterTool(weather.New())
	// More tools will be added as they're migrated from 66

	// ── Process input ────────────────────────────────────────────────────
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
