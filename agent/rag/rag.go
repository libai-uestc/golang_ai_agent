package rag

import (
	"context"
	"fmt"
	llm "libai/golang_ai_agent/agent"
	"log"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	qdrant_retriever "github.com/cloudwego/eino-ext/components/retriever/qdrant"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

type Doc struct {
	Id      string    `json:"id,omitempty"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Vector  []float64 `json:"vector,omitempty"`
}

const (
	QaCollection = "qa"
)

var (
	docMap = make(map[string]*Doc, 300)
)

func IndexDocument(markdownFile string) {
	ctx := context.Background()

	textParser := parser.TextParser{}
	loader, _ := file.NewFileLoader(ctx, &file.FileLoaderConfig{
		Parser: textParser,
	})
	docs, err := loader.Load(ctx, document.Source{URI: markdownFile})
	if err != nil {
		log.Fatal(err)
	}
	splitTransformer, _ := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"#":  "",
			"##": "",
		},
	})
	transformedDocs, _ := splitTransformer.Transform(ctx, docs)

	points := make([]*qdrant.PointStruct, 0, len(transformedDocs))

	for _, doc := range transformedDocs {
		DocId := uuid.NewString()
		var title, content string
		lines := strings.SplitN(doc.Content, "\n", 2)
		if strings.HasPrefix(lines[0], "## ") {
			title = lines[0][3:]
			content = lines[1]
		} else {
			continue
		}

		vectors := Embedding([]string{title, title, content})
		vector, _ := AvgOfVector(vectors)
		docMap[DocId] = &Doc{Id: DocId, Title: title, Content: content, Vector: vector}
		points = append(points, &qdrant.PointStruct{
			Id: &qdrant.PointId{PointIdOptions: &qdrant.PointId_Uuid{Uuid: DocId}},
			Vectors: &qdrant.Vectors{VectorsOptions: &qdrant.Vectors_Vector{
				Vector: &qdrant.Vector{Vector: &qdrant.Vector_Dense{Dense: &qdrant.DenseVector{Data: ToFloat32(vector)}}},
			}},
		})
	}

	GetQdrantClient().DeleteCollection(ctx, QaCollection) // 清空之前的数据
	if exists, err := GetQdrantClient().CollectionExists(ctx, QaCollection); err != nil {
		log.Fatal(err)
	} else {
		if !exists {
			err = GetQdrantClient().CreateCollection(ctx, &qdrant.CreateCollection{
				CollectionName: QaCollection,
				VectorsConfig: &qdrant.VectorsConfig{
					Config: &qdrant.VectorsConfig_Params{
						Params: &qdrant.VectorParams{
							Size:     VECTOR_DIM,
							Distance: qdrant.Distance_Cosine,
						},
					},
				},
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	result, err := GetQdrantClient().Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: QaCollection,
		Points:         points,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("upsert status %d", result.Status)
}

func RetrieveDocument(query string, limit int) []*Doc {
	ctx := context.Background()
	ScoreThreshold := 0.5
	retriever, _ := qdrant_retriever.NewRetriever(ctx, &qdrant_retriever.Config{
		Client:         GetQdrantClient(),
		Embedding:      GetEmbedder(),
		TopK:           limit,
		Collection:     QaCollection,
		ScoreThreshold: &ScoreThreshold,
	})
	neighbors, _ := retriever.Retrieve(ctx, query)
	log.Printf("retrieve %d docs", len(neighbors))

	result := make([]*Doc, 0, limit)
	for _, doc := range neighbors {
		id := doc.ID
		fmt.Println(id, doc.Score())
		if v, exists := docMap[id]; exists {
			result = append(result, &Doc{
				Title:   v.Title,
				Content: v.Content,
			})
		}
	}
	return result
}

func ChatBot(question string) string {
	ctx := context.Background()

	docs := RetrieveDocument(question, 4)
	knowledge := ""
	for _, doc := range docs {
		s, _ := sonic.MarshalString(doc)
		knowledge += s + "\n"
	}
	fmt.Println("【资料】", knowledge)

	model := llm.CreateDeepseekModel()
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Model:       model,
		Name:        "customer_service",
		Description: "客服机器人",
		Instruction: "请根据我提供的资料，回答用户的问题。",
	})
	if err != nil {
		log.Fatal(err)
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: agent,
	})
	iter := runner.Query(ctx, fmt.Sprintf("资料如下：%s\n用户问题：%s", knowledge, question))
	var lastMsg adk.Message
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatal(event.Err)
		}
		msg, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			log.Fatal(err)
		}
		lastMsg = msg
	}
	if lastMsg.Role == schema.Assistant && len(lastMsg.Content) > 0 {
		return lastMsg.Content
	} else {
		return "对不起，这个问题我还在学习中……"
	}
}
