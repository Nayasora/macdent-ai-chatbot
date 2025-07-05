package agent

import (
	"macdent-ai-chatbot/internal/databases"
)

type DeleteAgentRequest struct {
	AgentID string `json:"agent_id" validate:"required"`
}

func (s *Service) DeleteAgent(request *DeleteAgentRequest, postgres *databases.PostgresDatabase) {

}
