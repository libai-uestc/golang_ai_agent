package llm_test

import (
	"fmt"
	llm "libai/golang_ai_agent/agent"
	"strings"
	"testing"
)

func TestRunChatModel(t *testing.T) {
	llm.RunChatModel(false)
	fmt.Println(strings.Repeat("-", 100))
	llm.RunChatModel(true)
}

// go test -v ./agent -run=^TestRunChatModel$ -count=1

func TestRunChatModelAgent(t *testing.T) {
	llm.RunChatModelAgent()
}

// go test -v ./agent -run=^TestRunChatModelAgent$ -count=1
