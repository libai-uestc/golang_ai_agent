package mytool_test

import (
	"fmt"
	mytool "libai/golang_ai_agent/agent/my_tool"
	"testing"
)

func TestUserTimeTool(t *testing.T) {
	fmt.Print("纽约时间 ")
	mytool.UseTimeTool("America/New_York")
	fmt.Print("伦敦时间 ")
	mytool.UseTimeTool("Europe/London")
	fmt.Print("上海时间 ")
	mytool.UseTimeTool("")
}

func TestUseLocationTool(t *testing.T) {
	mytool.UseLocationTool()
}

func TestUseWeatherTool(t *testing.T) {
	mytool.UseWeatherTool(0)
	mytool.UseWeatherTool(1)
	mytool.UseWeatherTool(4)
	mytool.UseWeatherTool(5) // 数day非法
}
