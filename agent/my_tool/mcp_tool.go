package mytool

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	eino_mcp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// import (
// 	"context"
// 	"fmt"
// 	"sync"
// 	"time"

// 	"github.com/bytedance/sonic"
// 	eino_mcp "github.com/cloudwego/eino-ext/components/tool/mcp"
// 	"github.com/cloudwego/eino/components/tool"
// 	"github.com/mark3labs/mcp-go/client"
// 	"github.com/mark3labs/mcp-go/mcp"
// )

var (
	trainCli     *client.Client
	trainCliOnce sync.Once
)

const (
	TrainMcpName = "12306-mcp"
	TrainMcpUrl  = "https://mcp.api-inference.modelscope.net/d1003e8db1e94a/mcp"
)

func CreateMcpClient(url, name string) *client.Client {
	cli, err := client.NewStreamableHttpClient(url)
	if err != nil {
		log.Fatal(err)
	}
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    name,
		Version: "1.0.0",
	}
	_, err = cli.Initialize(context.Background(), initRequest)
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

func InitMcpClient() {
	trainCli = CreateMcpClient(TrainMcpUrl, TrainMcpName)
}

func CreateMcp12306Tool() []tool.BaseTool {
	ctx := context.Background()
	tools, err := eino_mcp.GetTools(ctx, &eino_mcp.Config{
		Cli: trainCli,
		// ToolNameList: []string{"tool1", "tool2"},
		ToolNameList: []string{},
	})
	if err != nil {
		log.Fatal(err)
	}
	return tools
}

func List12306Tool() {
	ctx := context.Background()
	tools := CreateMcp12306Tool()
	for _, tool := range tools {
		info, _ := tool.Info(ctx)
		fmt.Println(info.Name, info.Desc)
	}
}

type GetTicketsRequest struct {
	Date        string `json:"date"`
	FromStation string `json:"fromStation"`
	ToStation   string `json:"toStation"`
}

func Use12306Tool() {
	toolName := "get-tickets"
	ctx := context.Background()
	tools, err := eino_mcp.GetTools(ctx, &eino_mcp.Config{
		Cli:          trainCli,
		ToolNameList: []string{toolName},
	})
	if err != nil {
		log.Fatal(err)
	}
	ticketTool := tools[0].(tool.InvokableTool)
	request := GetTicketsRequest{
		Date:        time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		FromStation: "BJP",
		ToStation:   "SHH",
	}
	jsonStr, _ := sonic.MarshalString(request)
	resp, err := ticketTool.InvokableRun(ctx, jsonStr)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
	// ticketTool.InvokableRun(ctx)
}
