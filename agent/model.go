package llm

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/deepseek"
)

var (
	deepseekModel *deepseek.ChatModel
	arkModel      *ark.ChatModel
	dOnce         sync.Once
	aOnce         sync.Once
)

func CreateDeepseekModel() *deepseek.ChatModel {
	dOnce.Do(func() {
		ctx := context.Background()
		var err error

		deepseekModel, err = deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
			APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
			BaseURL: "https://api.deepseek.com/beta",
			Model:   "deepseek-chat",
		})
		if err != nil {
			log.Fatal(err)
		}
	})
	return deepseekModel
}

func CreateArkModel() *ark.ChatModel {
	aOnce.Do(func() {
		ctx := context.Background()
		var err error

		arkModel, err = ark.NewChatModel(ctx, &ark.ChatModelConfig{
			APIKey: os.Getenv("ARK_API_KEY"),
			Model:  "doubao-seed-2-1-pro-260915",
		})
		if err != nil {
			log.Fatal(err)
		}
	})
	return arkModel
}
