package llm

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	cbutils "github.com/cloudwego/eino/utils/callbacks"
)

type LoggerCallbacks struct{}

func (l *LoggerCallbacks) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	log.Printf("[INPUT] name: %v, type: %v, component: %v, input: %v", info.Name, info.Type, info.Component, input)
	return ctx
}

func (l *LoggerCallbacks) OnEnd(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackOutput) context.Context {
	log.Printf("[OUTPUT] name: %v, type: %v, component: %v,input: %v", info.Name, info.Type, info.Component, input)
	return ctx
}

func (l *LoggerCallbacks) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	log.Printf("[ERROR] name: %v, type: %v, component: %v, output: %v", info.Name, info.Type, info.Component, err)
	return ctx
}

func (l *LoggerCallbacks) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo, input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	return ctx
}

func GetStartCallback() callbacks.Handler {
	return callbacks.NewHandlerBuilder().OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
		log.Printf("[INPUT] name: %v, type: %v,component: %v,input: %v", info.Name, info.Type, info.Component, input)
		return ctx
	}).Build()
}

func GetChatModelInputCallback() callbacks.Handler {
	return cbutils.NewHandlerHelper().ChatModel(&cbutils.ModelCallbackHandler{
		OnStart: func(ctx context.Context, runInfo *callbacks.RunInfo, input *model.CallbackInput) context.Context {
			fmt.Printf("\n[ChatModel Input] component: %s\n", runInfo.Name)
			for i, msg := range input.Messages {
				fmt.Printf("  Message %d [%s]: %s\n", i+1, msg.Role, msg.Content)
				if len(msg.ToolCalls) > 0 {
					fmt.Printf("    Tool Calls: %d\n", len(msg.ToolCalls))
					for j, tc := range msg.ToolCalls {
						fmt.Printf("      %d. %s: %s\n", j+1, tc.Function.Name, tc.Function.Arguments)
					}
				}
			}
			return ctx
		},
	}).Handler()
}

// func GetChatModelInputCallback() callbacks.Handler {
// 	// 只给ChatModel注入CallbackHandler
// 	return cbutils.NewHandlerHelper().ChatModel(&cbutils.ModelCallbackHandler{
// 		// 这里只定义了OnStart，当然还可以指定OnEnd、OnError等
// 		OnStart: func(ctx context.Context, info *callbacks.RunInfo, input *model.CallbackInput) context.Context {
// 			// 打印ChatModel的输入
// 			fmt.Printf("\n[ChatModel Input] component: %s\n", info.Name)
// 			for i, msg := range input.Messages {
// 				fmt.Printf("  Message %d [%s]: %s\n", i+1, msg.Role, msg.Content)
// 				if len(msg.ToolCalls) > 0 {
// 					fmt.Printf("    Tool Calls: %d\n", len(msg.ToolCalls))
// 					for j, tc := range msg.ToolCalls {
// 						fmt.Printf("      %d. %s: %s\n", j+1, tc.Function.Name, tc.Function.Arguments)
// 					}
// 				}
// 			}
// 			return ctx
// 		},
// 	}).Handler()
// }

func GetToolInputCallback() callbacks.Handler {
	return cbutils.NewHandlerHelper().Tool(&cbutils.ToolCallbackHandler{
		OnStart: func(ctx context.Context, info *callbacks.RunInfo, input *tool.CallbackInput) context.Context {
			fmt.Printf("\n[Tool Input] component: %s, args: %s\n", info.Name, input.ArgumentsInJSON)
			return ctx
		},
	}).Handler()
}
