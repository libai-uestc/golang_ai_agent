package rag

// import (
// 	"context"
// 	// "fmt"
// 	"log"
// 	"sync"

// 	"github.com/cloudwego/eino-ext/components/retriever/qdrant"
// 	// qdrant_indexer "github.com/cloudwego/eino-ext/components/indexer/qdrant"
// 	// qdrant_retriever "github.com/cloudwego/eino-ext/components/retriever/qdrant"
// 	// "github.com/qdrant/go-client/qdrant"
// 	// "github.com/cloudwego/eino/schema"
// 	// "github.com/google/uuid"
// )

// var (
// 	qdCli     *qdrant.Client
// 	qdCliOnce sync.Once
// )

// func GetQdrantClient() *qdrant.Client {
// 	qdCliOnce.Do(func() {
// 		var err error
// 		qdCli, err = qdrant.NewClient(&qdrant.Config{
// 			Host: "localhost",
// 			Port: 6334,
// 		})
// 		if err != nil {
// 			log.Fatalf("failed to create qdrant client: %s", err)
// 		}
// 	})
// 	return qdCli
// }

// // 为大量文档创建索引，用于后续的检索。即先把文档转为向量，再存入向量数据库，入库的时候自然会创建索引
// func SearchByVector() {
// 	ctx := context.Background()

// 	CollectionName := "test_collection"
// 	if exists, _ := GetQdrantClient().CollectionExists(ctx, CollectionName); exists {
// 		if err := GetQdrantClient().DeleteCollection(ctx, CollectionName); err != nil {
// 			log.Fatal(err)
// 		} else {
// 			log.Println("delete collection", CollectionName)
// 		}
// 	}
// }
