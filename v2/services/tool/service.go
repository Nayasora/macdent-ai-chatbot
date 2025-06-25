package tool

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/openai/openai-go"
	"macdent-ai-chatbot/v2/clients"
	"macdent-ai-chatbot/v2/models"
	"macdent-ai-chatbot/v2/utils"
)

type Service struct {
	Agent  *models.Agent
	logger *log.Logger
}

type AgentArgumentError struct {
	Message string `json:"message"`
}

func NewService(agent *models.Agent) *Service {
	logger := utils.NewLogger("tool")

	return &Service{
		Agent:  agent,
		logger: logger,
	}
}

func (s *Service) HasToolCalls(toolCalls []openai.ChatCompletionMessageToolCall) bool {

	if len(toolCalls) == 0 {
		s.logger.Infof("вызов инструментов: не обнаружено")
		return false
	}

	s.logger.Infof("вызов инструментов: обнаружено %d вызовов", len(toolCalls))
	return true
}

func (s *Service) ExecuteToolCalls(messages []openai.ChatCompletionMessageParamUnion, toolMessage openai.ChatCompletionMessage) []openai.ChatCompletionMessageParamUnion {
	var toolResults []openai.ChatCompletionMessageParamUnion

	toolResults = append(toolResults, messages...)
	toolResults = append(toolResults, toolMessage)

	for _, toolCall := range toolMessage.ToolCalls {
		switch toolCall.Function.Name {
		case "get_doctors":
			s.logger.Info("вызов инструмента get_doctors")
			doctorsResponse, errorResponse := clients.GetDoctors(&clients.GetDoctorsRequest{
				AccessToken: s.Agent.Metadata.AccessToken,
			})

			if errorResponse != nil {
				s.logger.Errorf("обработка инструмента get_doctors: %v", errorResponse)
				errorJSON, err := json.Marshal(errorResponse)
				if err != nil {
					s.logger.Errorf("создание json: %v", err)
					continue
				}

				toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(errorJSON)))
				continue
			}

			responseJSON, err := json.Marshal(doctorsResponse)
			if err != nil {
				s.logger.Errorf("создание json: %v", err)
				continue
			}
			toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(responseJSON)))
		case "create_patient":
			s.logger.Info("вызов инструмента create_patient")

			var args map[string]interface{}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
				s.logger.Errorf("разбор аргументов инструмента create_patient: %v", err)

				errorObj := AgentArgumentError{Message: "Не удалось разобрать аргументы"}
				errorJSON, err := json.Marshal(errorObj)
				if err != nil {
					s.logger.Errorf("создание json: %v", err)
					continue
				}

				toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(errorJSON)))
				continue
			}

			name, ok := args["name"].(string)
			if !ok || name == "" {
				s.logger.Errorf("отсутствует имя пациента")

				errorObj := AgentArgumentError{Message: "Необходимо указать имя пациента"}
				errorJSON, err := json.Marshal(errorObj)
				if err != nil {
					s.logger.Errorf("создание json: %v", err)
					continue
				}

				toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(errorJSON)))
				continue
			}

			patientRequest := clients.CreatePatientRequest{
				AccessToken: s.Agent.Metadata.AccessToken,
				Name:        name,
			}

			patientResponse, errorResponse := clients.CreatePatient(patientRequest)

			if errorResponse != nil {
				s.logger.Errorf("обработка инструмента create_patient: %v", errorResponse)
				errorJSON, err := json.Marshal(errorResponse)
				if err != nil {
					s.logger.Errorf("создание json: %v", err)
					continue
				}

				toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(errorJSON)))
				continue
			}

			responseJSON, err := json.Marshal(patientResponse)
			if err != nil {
				s.logger.Errorf("создание json: %v", err)
				continue
			}

			toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(responseJSON)))
		case "get_schedule":
			s.logger.Info("вызов инструмента get_schedule")
			scheduleResponse, errorResponse := clients.GetSchedule(&clients.GetScheduleRequest{
				AccessToken: s.Agent.Metadata.AccessToken,
			})

			if errorResponse != nil {
				s.logger.Errorf("обработка инструмента get_schedule: %v", errorResponse)
				errorJSON, err := json.Marshal(errorResponse)
				if err != nil {
					s.logger.Errorf("создание json: %v", err)
					continue
				}

				toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(errorJSON)))
				continue
			}

			responseJSON, err := json.Marshal(scheduleResponse)
			if err != nil {
				s.logger.Errorf("создание json: %v", err)
				continue
			}
			toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(responseJSON)))
		case "create_appointment":
			s.logger.Info("вызов инструмента create_appointment")

			var args map[string]interface{}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
				s.logger.Errorf("разбор аргументов инструмента create_appointment: %v", err)

				errorObj := AgentArgumentError{Message: "Не удалось разобрать аргументы"}
				errorJSON, err := json.Marshal(errorObj)
				if err != nil {
					s.logger.Errorf("создание json: %v", err)
					continue
				}

				toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(errorJSON)))
				continue
			}

			patientID, ok := args["patient"].(float64)
			if !ok {
				s.logger.Errorf("неверный формат ID пациента: %v", args["patient"])

				errorObj := AgentArgumentError{Message: "Не указан ID пациента. Сначала создайте пациента с помощью инструмента create_patient и используйте полученный ID"}
				errorJSON, err := json.Marshal(errorObj)
				if err != nil {
					s.logger.Errorf("создание json: %v", err)
					continue
				}

				toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(errorJSON)))
				continue
			}

			doctorID, ok := args["doctor"].(float64)
			if !ok {
				s.logger.Errorf("неверный формат ID врача: %v", args["doctor"])

				errorObj := AgentArgumentError{Message: "Не указан ID врача. Сначала получите список врачей с помощью инструмента get_doctors и выберите нужный ID"}
				errorJSON, err := json.Marshal(errorObj)
				if err != nil {
					s.logger.Errorf("создание json: %v", err)
					continue
				}

				toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(errorJSON)))
				continue
			}

			date, ok := args["date"].(string)
			if !ok {
				s.logger.Errorf("неверный формат даты: %v", args["date"])

				errorObj := AgentArgumentError{Message: "Неверный формат даты"}
				errorJSON, err := json.Marshal(errorObj)
				if err != nil {
					s.logger.Errorf("создание json: %v", err)
					continue
				}

				toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(errorJSON)))
				continue
			}

			start, ok := args["start"].(string)
			if !ok {
				s.logger.Errorf("неверный формат времени начала: %v", args["start"])

				errorObj := AgentArgumentError{Message: "Неверный формат времени начала"}
				errorJSON, err := json.Marshal(errorObj)
				if err != nil {
					s.logger.Errorf("создание json: %v", err)
					continue
				}

				toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(errorJSON)))
				continue
			}

			end, ok := args["end"].(string)
			if !ok {
				s.logger.Errorf("неверный формат времени окончания: %v", args["end"])

				errorObj := AgentArgumentError{Message: "Неверный формат времени окончания"}
				errorJSON, err := json.Marshal(errorObj)
				if err != nil {
					s.logger.Errorf("создание json: %v", err)
					continue
				}

				toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(errorJSON)))
				continue
			}

			appointmentRequest := clients.CreateAppointmentRequest{
				AccessToken:          s.Agent.Metadata.AccessToken,
				DoctorID:             int(doctorID),
				PatientID:            int(patientID),
				AppointmentDate:      date,
				AppointmentStartTime: start,
				AppointmentEndTime:   end,
			}

			appointmentResponse, errorResponse := clients.CreateAppointment(appointmentRequest)

			if errorResponse != nil {
				s.logger.Errorf("обработка инструмента create_appointment: %v", errorResponse)
				errorJSON, err := json.Marshal(errorResponse)
				if err != nil {
					s.logger.Errorf("создание json: %v", err)
					continue
				}

				toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(errorJSON)))
				continue
			}

			responseJSON, err := json.Marshal(appointmentResponse)
			if err != nil {
				s.logger.Errorf("создание json: %v", err)
				continue
			}

			toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(responseJSON)))
		}
	}

	s.logger.Infof("Выполнение инструментов завершено, добавлено %d сообщений", len(toolResults)-len(messages)-1)
	return toolResults
}
