package devkit_test

import (
	devkit "libai/golang_ai_agent/agent/dev_kit"
	mytool "libai/golang_ai_agent/agent/my_tool"
	"testing"
)

func TestWeatherAgent(t *testing.T) {
	devkit.WeatherAgent()
}

// go test -v .\agent\dev_kit\ -run=^TestWeatherAgent$ -count=1

func TestLoopAgent(t *testing.T) {
	// 待优化的代码示例
	code := `
func processData(data []int) []int {
    result := []int{}
    for i := 0; i < len(data); i++ {
        for j := 0; j < len(data); j++ {
            if data[i] > data[j] {
                result = append(result, data[i])
                break
            }
        }
    }
    return result
}`
	devkit.LoopAgent(code)
}

// go test -v ./agent/dev_kit -run=^TestLoopAgent$ -count=1

func TestParallelAgent(t *testing.T) {
	devkit.ParalleAgent()
}

// go test -v ./agent/dev_kit -run=^TestParallelAgent$ -count=1

func TestSearchAgent(t *testing.T) {
	defer mytool.CloseBrowser()
	devkit.SearchAgent()
}

// go test -v ./agent/dev_kit -run=^TestSearchAgent$ -count=1
