package main

import (
	"macdent-ai-chatbot/api"
	"macdent-ai-chatbot/config"
)

func main() {
	cfg := config.NewConfig()

	api.NewServer(cfg.ApiServer).Run()
}
