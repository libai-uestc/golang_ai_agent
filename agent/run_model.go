package llm

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func createMessagesDirect() []*schema.Message {
	messages := []*schema.Message{
		&schema.Message{Role: schema.System, Content: ""},
		&schema.Message{Role: schema.User, Content: "今天是" + time.Now().Format("2006年01月02日") + ",离元旦放假还有几天"},
	}
	return messages
}

func createMessagesByTemplate() []*schema.Message {
	ctx := context.Background()

	template := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("今天是{today},离元旦放假还有几天"),
	)
	messages, err := template.Format(
		ctx,
		map[string]any{"today": time.Now().Format("2006年01月02日")},
	)

	if err != nil {
		log.Fatal(err)
	}
	return messages
}

func RunChatModel(stream bool) {
	ctx := context.Background()
	messages := createMessagesByTemplate()
	chatModel := CreateDeepseekModel()

	if !stream {
		respMsg, err := chatModel.Generate(ctx, messages)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(respMsg.Content)
	} else {
		streamResult, err := chatModel.Stream(ctx, messages)
		if err != nil {
			log.Fatal(err)
		}
		defer streamResult.Close()

		for {
			message, err := streamResult.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(message.Content)
		}
		fmt.Println()
	}
}

func RunChatModelAgent() {
	ctx := context.Background()
	messages := createMessagesByTemplate()
	chatModel := CreateDeepseekModel()

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "time_agent",
		Description: "时间计算器",
		Instruction: "",
		Model:       chatModel,
	})
	if err != nil {
		log.Fatal(err)
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: agent,
	})
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
		lastMsg = msg
	}

	fmt.Println("【最终答案】")
	if lastMsg.Role == schema.Assistant && len(lastMsg.Content) > 0 {
		fmt.Println(lastMsg.Content)
	}
}
