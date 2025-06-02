package main

import (
	"macdent-ai-chatbot/api"
	"macdent-ai-chatbot/config"
	"macdent-ai-chatbot/store"
)

func main() {
	cfg := config.NewConfig()

	_ = store.NewQdrantClient(cfg.Qdrant)

	api.NewServer(cfg.ApiServer).
		Run()
}
