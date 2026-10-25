package devkit

import (
	"context"
	"fmt"
	llm "libai/golang_ai_agent/agent"
	mytool "libai/golang_ai_agent/agent/my_tool"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func TripAgent() {
	ctx := context.Background()
	// callbacks.AppendGlobalHandlers(new(llm.LoggerCallbacks))        // 所有component的CallBack全开
	// callbacks.AppendGlobalHandlers(llm.GetStartCallback()) // 所有component只开OnStartCallBack
	// callbacks.AppendGlobalHandlers(llm.GetEndCallback())   // 所有component只开OnEndCallBack
	callbacks.AppendGlobalHandlers(llm.GetChatModelInputCallback())
	callbacks.AppendGlobalHandlers(llm.GetToolInputCallback())
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
		Name:        "trip_plan",
		Description: "出行规划",
		Instruction: "请帮我完成旅行规划，包括什么时间去什么景点，并列出详细的高铁车次和出发时间",
	})
	if err != nil {
		log.Fatal(err)
	}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: agent,
	})

	iter := runner.Query(ctx, "我想去上海旅游3天(从明天开始)，请帮我做一份旅游攻略")
	var lastMsg adk.Message
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatal(event.Err)
		}
		msg, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			log.Fatal(err)
		}
		lastMsg = msg
	}

	fmt.Println("【最终答案】")
	if lastMsg.Role == schema.Assistant && len(lastMsg.Content) > 0 {
		fmt.Println(lastMsg.Content)
	}
}
