package main

import (
	"context"
	mytool "libai/golang_ai_agent/agent/my_tool"
	"log"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TimeInput struct {
	TimeZone string `json:"time_zone" jsonschema:"当地时区"`
}

type TimeOutput struct {
	Time string `json:"time" jsonschema:"当前时间"`
}

func GetCurrentTime(ctx context.Context, req *mcp.CallToolRequest, input TimeInput) (*mcp.CallToolResult, *TimeOutput, error) {
	if len(input.TimeZone) == 0 {
		input.TimeZone = "Asia/Shanghai"
	}
	loc, err := time.LoadLocation(input.TimeZone)
	if err != nil {
		return nil, nil, err
	}
	now := time.Now().In(loc)
	return nil, &TimeOutput{now.Format("2006-01-02 15:04:05")}, nil
}

func GetCurrentLocation(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, *mytool.Location, error) {
	location, err := mytool.GetMyLocation()
	return nil, location, err
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "libai_mcp",
		Version: "1.0.0",
	}, nil)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "time_tool",
		Description: "获取当前的时间",
	}, GetCurrentTime)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "location_tool",
		Description: "获取当前的地理位置，包含省、城市（含城市名称和城市编码）、经纬度等信息",
	}, GetCurrentLocation)

	handler := mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return server },
		&mcp.StreamableHTTPOptions{},
	)
	http.HandleFunc("/mcp", handler.ServeHTTP)

	log.Println("Starting MCP server on :5678")
	if err := http.ListenAndServe("127.0.0.1:5678", nil); err != nil {
		log.Fatal(err)
	}
}
