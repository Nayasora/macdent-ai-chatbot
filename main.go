package main

import "macdent-ai-chatbot/v2/api"
import "macdent-ai-chatbot/v2/configs"

func main() {
	config := configs.NewConfig(
		configs.NewEnv(".env"),
	)

	server := api.NewServer(config)
	server.Setup()
	server.Run()
}
