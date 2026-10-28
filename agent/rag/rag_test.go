package rag_test

import (
	"fmt"
	"libai/golang_ai_agent/agent/rag"
	"strings"
	"testing"
)

func TestParseDocument(t *testing.T) {
	rag.ParseDocument("../../data/qa.md") // 相对路径是相对于本_test.go文件的路径
	fmt.Println(strings.Repeat("~", 100))
	rag.ParseDocument("../../data/qa.html")
	fmt.Println(strings.Repeat("~", 100))
	rag.ParseDocument("../../data/qa.pdf")
	fmt.Println(strings.Repeat("~", 100))
	rag.ParseDocument("../../data/df.xlsx")
	fmt.Println(strings.Repeat("~", 100))
	rag.ParseDocument("../../data/观沧海.docx")
	fmt.Println(strings.Repeat("~", 100))
}

// go test -v ./agent/rag -run=^TestParseDocument$ -count=1

func TestLoadDocument(t *testing.T) {
	rag.LoadDocument("../../data/qa.md") // 相对路径是相对于本_test.go文件的路径
	fmt.Println(strings.Repeat("~", 100))
	rag.LoadDocument("../../data/qa.html")
	fmt.Println(strings.Repeat("~", 100))
	rag.LoadDocument("../../data/qa.pdf")
	fmt.Println(strings.Repeat("~", 100))
	rag.LoadDocument("../../data/df.xlsx")
	fmt.Println(strings.Repeat("~", 100))
	rag.LoadDocument("../../data/观沧海.docx")
	fmt.Println(strings.Repeat("~", 100))
}

// go test -v ./agent/rag -run=^TestLoadDocument$ -count=1

func TestTransformDocument(t *testing.T) {
	rag.TransformDocument("../../data/qa.md")
}

// go test -v ./agent/rag -run=^TestTransformDocument$ -count=1

func TestWordSim(t *testing.T) {
	rag.WordSim()
}

// go test -v ./agent/rag -run=^TestWordSim$ -count=1
