package einograph

import (
	"context"
	"fmt"
	llm "libai/golang_ai_agent/agent"
	mytool "libai/golang_ai_agent/agent/my_tool"
	"log"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func GraphOfReact() {
	ctx := context.Background()

	chatModel := llm.CreateDeepseekModel()
	tools := []tool.BaseTool{
		mytool.CreateTimeTool(),
		mytool.CreateWeatherTool(),
		mytool.CreateLocationTool(),
	}
	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: tools,
	})
	if err != nil {
		log.Fatal(err)
	}
	toolInfos := make([]*schema.ToolInfo, 0, len(tools))
	for _, tool := range tools {
		if info, err := tool.Info(ctx); err == nil {
			toolInfos = append(toolInfos, info)
		}
	}
	chatModel.BindTools(toolInfos)

	addHistory := compose.InvokableLambda(func(ctx context.Context, input []*schema.Message) (output []*schema.Message, err error) {
		return
	},
	)

	graph := compose.NewGraph[[]*schema.Message, *schema.Message](
		compose.WithGenLocalState(func(ctx context.Context) (state *HistoryMessage) {
			return &HistoryMessage{
				History:       make([]*schema.Message, 0, 4),
				UserMessage:   "",
				SystemMessage: "",
			}
		}),
	)
	graph.AddChatModelNode("model", chatModel,
		compose.WithStatePreHandler(func(ctx context.Context, input []*schema.Message, state *HistoryMessage) ([]*schema.Message, error) {
			for _, msg := range input {
				switch msg.Role {
				case schema.System:
					state.SystemMessage = msg.Content
				case schema.User:
					state.UserMessage = msg.Content
				}
			}
			return input, nil
		}),
		compose.WithStatePostHandler(func(ctx context.Context, output *schema.Message, state *HistoryMessage) (*schema.Message, error) {
			state.History = append(state.History, output)
			return output, nil
		}),
	)
	graph.AddToolsNode("tools", toolsNode,
		compose.WithStatePostHandler(func(ctx context.Context, out []*schema.Message, state *HistoryMessage) ([]*schema.Message, error) {
			state.History = append(state.History, out...)
			return out, nil
		}),
	)
	graph.AddLambdaNode("history", addHistory,
		compose.WithStatePostHandler(func(ctx context.Context, out []*schema.Message, state *HistoryMessage) ([]*schema.Message, error) {
			result := []*schema.Message{{Role: schema.System, Content: state.SystemMessage}}
			result = append(result, state.History...)
			result = append(result, &schema.Message{Role: schema.User, Content: state.UserMessage})
			return result, nil
		},
		))
	graph.AddEdge(compose.START, "model")
	graph.AddBranch("model", compose.NewGraphBranch(func(ctx context.Context, in *schema.Message) (string, error) {
		if len(in.ToolCalls) > 0 {
			return "tools", nil
		} else {
			return compose.END, nil
		}
	}, map[string]bool{compose.END: true, "tools": true}))
	graph.AddEdge("tools", "history")
	graph.AddEdge("history", "model")

	runnable, err := graph.Compile(ctx)
	if err != nil {
		log.Fatal(err)
	}
	input := []*schema.Message{
		{Role: schema.User, Content: time.Now().Add(48*time.Hour).Format("2006年01月02日") + "适合开运动会吗？"},
	}
	msg, err := runnable.Invoke(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(msg.Content)
}
