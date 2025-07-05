package clients

import (
	"encoding/json"
	"io"
	"macdent-ai-chatbot/internal/utils"
	"net/http"
)

type CreatePatientRequest struct {
	AccessToken string `json:"access_token"`
	Name        string `json:"name"`
}

type PatientInfo struct {
	Name    string `json:"name"`
	ID      int    `json:"id"`
	Comment string `json:"comment"`
}

type CreatePatientResponse struct {
	Patient  PatientInfo `json:"patient"`
	Response int         `json:"response"`
}

func CreatePatient(request CreatePatientRequest) (*CreatePatientResponse, *utils.UserErrorResponse) {
	logger := utils.NewLogger("clients:createPatient")

	baseURL := BaseURL + "patient/add"
	logger.Infof("создание пациента через URL: %s", baseURL)

	reqURL := baseURL + "?access_token=" + request.AccessToken + "&name=" + request.Name

	// Создаем запрос без тела
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		logger.Errorf("создание запроса: %v", err)
		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка при создании запроса к API",
			"Пожалуйста, попробуйте позже или обратитесь в службу поддержки.",
		)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("выполнение запроса: %v", err)
		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка при выполнении запроса к API",
			"Пожалуйста, попробуйте позже или обратитесь в службу поддержки.",
		)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Fatalf("закрытие тела ответа: %v", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("чтение ответа: %v", err)
		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка при чтении ответа от API",
			"Пожалуйста, попробуйте позже или обратитесь в службу поддержки.",
		)
	}

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("получен неверный статус запроса: %d", resp.StatusCode)
		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка при выполнении запроса к API",
			"Пожалуйста, попробуйте позже или обратитесь в службу поддержки.",
		)
	}

	var response CreatePatientResponse
	if err := json.Unmarshal(body, &response); err != nil {
		logger.Errorf("десериализация ответа: %v", err)
		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка при обработке ответа от API",
			"Пожалуйста, попробуйте позже или обратитесь в службу поддержки.",
		)
	}

	if response.Response != 1 {
		logger.Errorf("API вернул неуспешный ответ: %d", response.Response)
		return nil, utils.NewUserErrorResponse(
			500,
			"API вернул ошибку",
			"Пожалуйста, попробуйте позже или обратитесь в службу поддержки.",
		)
	}

	return &response, nil
}
