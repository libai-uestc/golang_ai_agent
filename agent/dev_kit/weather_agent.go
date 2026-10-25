package devkit

import (
	"context"
	"fmt"
	llm "libai/golang_ai_agent/agent"
	mytool "libai/golang_ai_agent/agent/my_tool"
	"log"
	"strings"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func WeatherAgent() {
	ctx := context.Background()
	model := llm.CreateDeepseekModel()

	tools := []tool.BaseTool{mytool.CreateTimeTool(), mytool.CreateLocationTool(), mytool.CreateWeatherTool()}

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Model: model,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
		Name:        "activity_arrangement",
		Description: "活动规划大师",
		Instruction: "请根据未来的天气安排活动",
	})
	if err != nil {
		log.Fatal(err)
	}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: agent,
	})
	messages := []adk.Message{
		&schema.Message{Role: schema.User, Content: time.Now().Format("2006年01月02日") + "适合开运动会吗？"},
	}
	// iter := runner.Query(ctx, time.Now().Format("2006年01月02日")+"适合开运动会吗？")
	iter := runner.Run(ctx, messages)

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
		fmt.Println(msg.Role, msg.Content)
		lastMsg = msg
	}
	fmt.Println("【最终答案】")
	if lastMsg.Role == schema.Assistant && len(lastMsg.Content) > 0 {
		fmt.Println(lastMsg.Content)
		messages = append(messages, lastMsg)
	}
	fmt.Println(strings.Repeat("-", 100))

	messages = append(messages, &schema.Message{Role: schema.User, Content: "运动会的前天晚上会结冰吗？"})
	iter = runner.Run(ctx, messages)

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
		fmt.Println(msg.Role, msg.Content)
	}
	fmt.Println()
	if lastMsg.Role == schema.Assistant && len(lastMsg.Content) > 0 {
		fmt.Println(lastMsg.Content)
		messages = append(messages, lastMsg)
	}
}
