package einograph_test

import (
	einograph "libai/golang_ai_agent/agent/eino_graph"
	"testing"
)

func TestComposeGraph(t *testing.T) {
	einograph.ComposeGraph("hello, how are you")
	einograph.ComposeGraph("hello, how are you!")
}

// go test -v ./agent/eino_graph -run=^TestComposeGraph$ -count=1

func TestGraphWithCallBack(t *testing.T) {
	einograph.GraphWithCallBack("hello, how are you")
	einograph.GraphWithCallBack("hello, how are you!")
	// 	einograph.GrapWithCallBack("hello, how are you")
	// 	einograph.GraphWithCallBack("hello, how are you!")
}

// go test -v ./agent/eino_graph -run=^TestGraphWithCallBack$ -count=1

func TestGraphWithState(t *testing.T) {
	einograph.GraphWithState("hello, how are you")
	einograph.GraphWithState("hello, how are you!")
}
