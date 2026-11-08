package main

import (
	"context"
	llm "libai/golang_ai_agent/agent"
	mytool "libai/golang_ai_agent/agent/my_tool"
	"log"
	"sync"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

var (
	runner     *adk.Runner
	runnerOnce sync.Once
)

func GetRunner() *adk.Runner {
	runnerOnce.Do(func() {
		ctx := context.Background()
		callbacks.AppendGlobalHandlers(llm.GetChatModelInputCallback())
		model := llm.CreateDeepseekModel()

		tools := []tool.BaseTool{
			mytool.CreateTimeTool(),
			mytool.CreateLocationTool(),
			mytool.CreateWeatherTool(),
		}
		agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
			Model: model,
			ToolsConfig: adk.ToolsConfig{
				ToolsNodeConfig: compose.ToolsNodeConfig{
					Tools: tools,
				},
			},
			Name:        "my_agent",
			Description: "我的万能AI助手",
			Instruction: "",
		})
		if err != nil {
			log.Fatal(err)
		}
		runner = adk.NewRunner(ctx, adk.RunnerConfig{
			Agent:           agent,
			EnableStreaming: true,
		})

	})
	return runner
}
