package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Message struct {
	Role             string `json:"role"`
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content,omitempty"`
}

type RequestBody struct {
	Model    string     `json:"model"`
	Messages []*Message `json:"messages"`
	Stream   bool       `json:"stream"`
}

type Choice struct {
	Index   int      `json:"index"`
	Message *Message `json:"message"`
}

type StreamChoice struct {
	Index int      `json:"index"`
	Delta *Message `json:"delta"`
}

type ResponseBody struct {
	Choice []*Choice `json:"choices"`
}

type StreamResponseBody struct {
	Choices []*StreamChoice `json:"choices"`
}

type ResponseRecord struct {
	Timestamp        string `json:"timestamp"`
	Model            string `json:"model"`
	ReasoningContent string `json:"reasoning_content"`
	Content          string `json:"content"`
}

const (
	DeepseekUrl = "https://api.deepseek.com/chat/completions"
)

func ChatWithDeepseek(rb *RequestBody) {
	bs, err := json.Marshal(rb)
	if err != nil {
		log.Printf("json序列化失败: %s", err)
		return
	}

	request, err := http.NewRequest(http.MethodPost, DeepseekUrl, bytes.NewReader(bs))
	if err != nil {
		log.Printf("创建请求失败: %s", err)
		return
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+os.Getenv("DEEPSEEK_API_KEY"))
	client := http.Client{
		Timeout: 60 * time.Minute,
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Printf("http请求失败: %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, resp.Body)
		return
	}

	// bs, _ = io.ReadAll(resp.Body)
	// fmt.Println(string(bs))

	// 打印非流式的
	// var respBody ResponseBody
	// err = json.Unmarshal(bs, &respBody)
	// if err == nil {
	//  fmt.Println(respBody.Choice[0].Message.Role, respBody.Choice[0].Message.Content)
	// } else {
	//  fmt.Println(string(bs))
	// }

	// 打印流式的返回结果
	// 缓存思考链与正文内容
	var reasoningBuilder strings.Builder
	var contentBuilder strings.Builder

	reader := bufio.NewReader(resp.Body)

	// buffer := make([]byte, 2048)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("读取数据流失败: %s", err)
			break
		}
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "data:") {
			continue
		}

		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "[DONE]" {
			break
		}

		var chunk StreamResponseBody
		if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
			log.Printf("json反序列化失败: %s. raw: %s", err, payload)
			continue
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			delta := chunk.Choices[0].Delta
			if delta.ReasoningContent != "" {
				reasoningBuilder.WriteString(delta.ReasoningContent)
			}
			if delta.Content != "" {
				contentBuilder.WriteString(delta.Content)
				fmt.Print(delta.Content)
			}
		}

		// if err != nil {
		// 	log.Printf("read response body failed: %s", err)
		// 	break
		// }
		// fmt.Println(string(buffer[:n]))

		// for _, segment := range bytes.Split(buffer[:n], []byte("data: ")) {
		// 	segment = bytes.TrimSpace(segment)
		// 	if len(segment) == 0 {
		// 		continue
		// 	}
		// 	if string(segment) == "[DONE]" {
		// 		break
		// 	}
		// 	var response StreamResponseBody
		// 	err := json.Unmarshal(segment, &response)
		// 	if err == nil {
		// 		fmt.Print(response.Choice[0].Message.Content)
		// 	} else {
		// 		log.Printf("json反序列化失败: %s. [%s]", err, string(segment))
		// 		break
		// 	}
		// }
	}
	fmt.Println()

	saveToFile(rb.Model, reasoningBuilder.String(), contentBuilder.String())

}

func saveToFile(model, reasoningContent, content string) {
	record := ResponseRecord{
		Timestamp:        time.Now().Format("2006-01-02 15:04:05"),
		Model:            model,
		ReasoningContent: reasoningContent,
		Content:          content,
	}

	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		log.Printf("序列化存档数据失败: %s", err)
		return
	}

	fileName := fmt.Sprintf("response_%s.json", time.Now().Format("20060102_150405"))
	if err := os.WriteFile(fileName, data, 0644); err != nil {
		log.Printf("保存文件失败: %s", err)
		return
	}

	log.Printf("完整数据（含思考链）已保存至: %s\n", fileName)
}

func main() {
	_ = godotenv.Load(".env", "golang_ai_agent/agent/native/.env")
	rb := &RequestBody{
		Model: "deepseek-flash",
		Messages: []*Message{
			&Message{Role: "system", Content: ""},
			&Message{Role: "user", Content: "给我讲讲谷歌的Agent Development Kit、LangChain、LangGraph、Agent Executor、Agent Substrate、Genkit、Gemini CLI、字节跳动的Eino、DeerFlow、HiAgent和AgentKit"},
		},
		Stream: true,
	}
	ChatWithDeepseek(rb)
}
