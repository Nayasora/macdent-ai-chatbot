package main

import "macdent-ai-chatbot/internal/api"
import "macdent-ai-chatbot/internal/configs"

func main() {
	config := configs.NewConfig(
		configs.NewEnv(".env"),
	)

	server := api.NewServer(config)
	server.Setup()
	server.Run()
}
