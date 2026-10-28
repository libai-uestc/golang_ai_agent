package rag

import (
	"context"
	"fmt"
	"log"
	"sync"

	qdrant_indexer "github.com/cloudwego/eino-ext/components/indexer/qdrant"
	qdrant_retriever "github.com/cloudwego/eino-ext/components/retriever/qdrant"
	"github.com/qdrant/go-client/qdrant"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
)

var (
	qdCli     *qdrant.Client
	qdCliOnce sync.Once
)

func GetQdrantClient() *qdrant.Client {
	qdCliOnce.Do(func() {
		var err error
		qdCli, err = qdrant.NewClient(&qdrant.Config{
			Host: "localhost",
			Port: 6334,
		})
		if err != nil {
			log.Fatalf("failed to create qdrant client: %s", err)
		}
	})
	return qdCli
}

// 为大量文档创建索引，用于后续的检索。即先把文档转为向量，再存入向量数据库，入库的时候自然会创建索引
func SearchByVector() {
	ctx := context.Background()

	CollectionName := "test_collection"
	if exists, _ := GetQdrantClient().CollectionExists(ctx, CollectionName); exists {
		if err := GetQdrantClient().DeleteCollection(ctx, CollectionName); err != nil {
			log.Fatal(err)
		} else {
			log.Println("delete collection", CollectionName)
		}
	}

	indexer, err := qdrant_indexer.NewIndexer(ctx, &qdrant_indexer.Config{
		Client:     GetQdrantClient(),
		Embedding:  GetEmbedder(),
		VectorDim:  VECTOR_DIM,
		Distance:   qdrant.Distance_Cosine,
		Collection: CollectionName,
	})
	if err != nil {
		log.Fatalf("failed to create qdrant indexer: %s", err)
	}
	log.Println("create indexer")

	docs := []*schema.Document{
		{
			ID:      uuid.NewString(), // 对于qdrant来说，必须使用uuid
			Content: "美国9个国家公共假日，元旦：1月1日；马丁路德金日：一月的第三个星期一；总统日：二月的第三个星期一；阵亡将士纪念日：五月的最后一个星期一；独立日：7月4日；劳动节：九月的第一个星期一；退伍军人节：11月11日；感恩节：11月的第四个星期四；圣诞节：12月25日。\n",
			MetaData: map[string]any{ // 根据业务需求，可以往MetaData里添加任意数据，比如：文档的来源信息、文档的分数（用于排序）、文档的子索引（用于分层检索）等等
				"source": "sohu",
			},
		},
		{
			ID:      uuid.NewString(),
			Content: "美国平均退休年龄：66.33岁",
			MetaData: map[string]any{
				"source": "sina",
			},
		},
		{
			ID:      uuid.NewString(),
			Content: "美国平均退休年龄：66.33岁",
			MetaData: map[string]any{
				"source": "sohu",
			},
		},
	}
	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		log.Fatalf("failed to store documents to indexer: %s", err)
	}
	log.Printf("Store success, ids: %v", ids)

	scoreThresh := 0.3
	retriever, _ := qdrant_retriever.NewRetriever(ctx, &qdrant_retriever.Config{
		Client:         GetQdrantClient(),
		Embedding:      GetEmbedder(),
		TopK:           20,
		ScoreThreshold: &scoreThresh,
		Collection:     CollectionName,
	})
	query := "美国多少岁可以退休？"
	neighbors, _ := retriever.Retrieve(ctx, query,
		qdrant_retriever.WithFilter(&qdrant.Filter{
			Must: []*qdrant.Condition{
				qdrant.NewMatch("metadata.source", "sohu"),
			},
		}),
	)

	for i, doc := range neighbors {
		md := doc.MetaData["metadata"]
		mp := md.(map[string]*qdrant.Value)
		source := mp["source"].GetStringValue()
		fmt.Printf("%d  %.4f  %s  %s\n", i, doc.Score(), doc.ID, source)
	}
}
