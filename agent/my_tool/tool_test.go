package mytool_test

import (
	"fmt"
	mytool "libai/golang_ai_agent/agent/my_tool"
	"testing"
)

func TestUseTimeTool(t *testing.T) {
	fmt.Print("纽约时间 ")
	mytool.UseTimeTool("America/New_York")
	fmt.Print("伦敦时间 ")
	mytool.UseTimeTool("Europe/London")
	fmt.Print("上海时间 ")
	mytool.UseTimeTool("")
}

// go test -v ./agent/my_tool -run=^TestUseTimeTool$ -count=1
func TestUseLocationTool(t *testing.T) {
	mytool.UseLocationTool()
}

// go test -v ./agent/my_tool -run=^TestUseLocationTool$ -count=1
func TestUseWeatherTool(t *testing.T) {
	mytool.UseWeatherTool(0)
	mytool.UseWeatherTool(1)
	mytool.UseWeatherTool(4)
	mytool.UseWeatherTool(5) // 数day非法
}

// go test -v ./agent/my_tool -run=^TestUseWeatherTool$ -count=1

func TestMcpTool(t *testing.T) {
	mytool.InitMcpClient()
	// mytool.List12306Tool()
	mytool.Use12306Tool()
}

// go test -v ./agent/my_tool -run=^TestMcpTool$ -count=1
