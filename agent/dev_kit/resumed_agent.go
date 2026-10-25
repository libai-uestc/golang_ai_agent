package devkit

import (
	"bufio"
	"context"
	"fmt"
	llm "libai/golang_ai_agent/agent"
	checkpointstore "libai/golang_ai_agent/agent/checkpoint_store"
	mytool "libai/golang_ai_agent/agent/my_tool"
	"log"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/rs/xid"
)

func InterruptResume() {
	ctx := context.Background()

	model := llm.CreateDeepseekModel()
	tools := []tool.BaseTool{
		mytool.CreateTimeTool(),
		mytool.CreateWeatherTool(),
		mytool.NewAskForClarificationTool(),
	}
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Model: model,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
		Name:        "activity_arrangement",
		Description: "活动计划大师",
		Instruction: "请根据未来的天气安排活动",
	})
	if err != nil {
		log.Fatal(err)
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		CheckPointStore: checkpointstore.NewInMemoryStore(),
	})
	messages := []adk.Message{
		&schema.Message{Role: schema.User, Content: time.Now().Add(48*time.Hour).Format("2006年01月02日") + "适合开运动会吗?"},
	}
	checkPointId := xid.New().String()
	iter := runner.Run(ctx, messages, adk.WithCheckPointID(checkPointId))
	var lastMsg adk.Message
LB:
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

		if len(msg.ToolCalls) > 0 {
			for _, toolCall := range msg.ToolCalls {
				if toolCall.Function.Name == mytool.AskForClarificationToolName {
					var req mytool.AskForClarificationInput
					if err := sonic.UnmarshalString(toolCall.Function.Arguments, &req); err == nil {
						fmt.Println(req.Question)
						scanner := bufio.NewScanner(os.Stdin)
						fmt.Println("请在此输入答案：")
						scanner.Scan()
						fmt.Println()
						input := scanner.Text()

						iter, err = runner.Resume(ctx, checkPointId, adk.WithToolOptions([]tool.Option{mytool.WithNewInput(input)}))
						if err != nil {
							log.Fatal(err)
						}
						continue LB
					} else {
						log.Fatal(err)
					}
				}
			}
		}

		fmt.Println(msg.Role, msg.Content)
		lastMsg = msg
	}

	fmt.Println("【最终答案】")
	if lastMsg.Role == schema.Assistant && len(lastMsg.Content) > 0 {
		fmt.Println(lastMsg.Content)
		messages = append(messages, lastMsg)
	}
}
