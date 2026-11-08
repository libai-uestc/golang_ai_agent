package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/rs/xid"
)

var (
	session sync.Map
)

func ParseUrlParams(rawQuery string) map[string]string {
	params := make(map[string]string, 10)
	args := strings.Split(rawQuery, "&")
	for _, ele := range args {
		arr := strings.Split(ele, "=")
		if len(arr) == 2 {
			key, _ := url.QueryUnescape(arr[0]) // url参数反转义
			value, _ := url.QueryUnescape(arr[1])
			params[key] = value
		}
	}
	return params
}

func Chat(w http.ResponseWriter, r *http.Request) {
	log.Println("call Chat")
	ctx := context.Background()
	w.Header().Add("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	params := ParseUrlParams(r.URL.RawQuery)
	messages := make([]adk.Message, 0, 10)
	sid := params["session"]
	if v, exists := session.Load(sid); exists {
		history := v.([]adk.Message)
		messages = append(messages, history...)
	}
	messages = append(messages, &schema.Message{Role: schema.User, Content: params["msg"]})
	iter := GetRunner().Run(ctx, messages)
	answer := strings.Builder{}
	for {
		event, ok := iter.Next()
		if !ok {
			fmt.Fprint(w, "data: [DONE]\n\n")
			flusher.Flush()
			messages = append(messages, &schema.Message{Role: schema.Assistant, Content: answer.String()})
			session.Store(sid, messages)
			break
		}
		if event.Err != nil {
			log.Printf("read llm response failed: %s", event.Err)
			break
		}
		s := event.Output.MessageOutput.MessageStream
		for {
			msg, err := s.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("get output message failed: %s", event.Err)
				break
			}
			if msg != nil {
				fmt.Print(msg.Content)
				answer.WriteString(msg.Content)
				fmt.Fprintf(w, "data: %s\n\n", strings.ReplaceAll(msg.Content, "\n", "<br>"))
				flusher.Flush()
			}
		}
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(writer http.ResponseWriter, request *http.Request) {
		tmpl, err := template.ParseFiles("./web/chat.html")
		if err != nil {
			fmt.Println("create template failed:", err)
			return
		}
		sid := xid.New().String()
		err = tmpl.Execute(writer, map[string]string{"url": "http://127.0.0.1:5678/chat", "session": sid})
		if err != nil {
			log.Fatal(err)
		}
	})
	mux.HandleFunc("GET /favicon.ico", func(writer http.ResponseWriter, request *http.Request) {
		file, err := os.Open("./web/迈克尔乔丹.png")
		if err != nil {
			return
		}
		defer file.Close()
		io.Copy(writer, file)
	})
	mux.HandleFunc("GET /chat", Chat)
	log.Println("going to start web server on port 5678")
	if err := http.ListenAndServe("127.0.0.1:5678", mux); err != nil {
		panic(err)
	}
}
