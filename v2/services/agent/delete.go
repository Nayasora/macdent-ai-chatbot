package agent

import (
	"macdent-ai-chatbot/v2/databases"
)

type DeleteAgentRequest struct {
	AgentID string `json:"agent_id" validate:"required"`
}

func (s *Service) DeleteAgent(request *DeleteAgentRequest, postgres *databases.PostgresDatabase) {

}
