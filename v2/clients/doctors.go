package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"macdent-ai-chatbot/v2/utils"
	"net/http"
	"net/url"
)

type GetDoctorsRequest struct {
	AccessToken string `json:"access_token"`
	FullName    string `json:"name"`
}

type Specialty struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Doctor struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Specialties []Specialty `json:"specialnosti"`
}

type GetDoctorsResponse struct {
	Doctors  []Doctor `json:"doctors"`
	Count    string   `json:"count"`
	AtPage   int      `json:"atPage"`
	MaxPage  int      `json:"maxPage"`
	Response int      `json:"response"`
}

func GetDoctors(request *GetDoctorsRequest) (*GetDoctorsResponse, *utils.UserErrorResponse) {
	logger := utils.NewLogger("clients:getDoctors")

	baseURL := BaseURL + "doctor/find"

	params := url.Values{}
	params.Add("access_token", request.AccessToken)
	if request.FullName != "" {
		params.Add("name", request.FullName)
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	logger.Infof("получение списка врачей по URL: %s", fullURL)

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
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

	var response GetDoctorsResponse
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
