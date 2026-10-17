package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "libai_mcp",
		Version: "1.0.0",
	}, nil)

	transport := &mcp.StreamableClientTransport{
		Endpoint:   "http://127.0.0.1:5678/mcp",
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	resp, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "time_tool",
		Arguments: map[string]any{"time_zone": ""},
	})
	if err != nil {
		log.Fatal("call time_cool failed", err)
	}
	if resp.IsError {
		log.Fatal("call time_tool error", err)
	}
	for _, c := range resp.Content {
		fmt.Println("当地时间 ", c.(*mcp.TextContent).Text)
	}

	resp, err = session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "time_tool",
		Arguments: map[string]any{"time_zone": "Europe/London"},
	})
	if err != nil {
		log.Fatal("call time_tool failed", err)
	}
	if resp.IsError {
		log.Fatal("call time_tool error", err)
	}
	for _, c := range resp.Content {
		fmt.Println("伦敦时间 ", c.(*mcp.TextContent).Text)
	}

	resp, err = session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "location_tool",
		Arguments: map[string]any{},
	})
	if err != nil {
		log.Fatal("call location_tool failed", err)
	}
	if resp.IsError {
		log.Fatal("call location_tool error")
	}
	for _, c := range resp.Content {
		fmt.Println(c.(*mcp.TextContent).Text)
	}
}
