package main

import "macdent-ai-chatbot/store"

func main() {
	_ = store.New(&store.QdrantConfig{
		Host: "localhost",
		Port: 6333,
	})
}
