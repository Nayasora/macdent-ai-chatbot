package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"macdent-ai-chatbot/v2/utils"
	"net/http"
	"net/url"
)

type GetScheduleRequest struct {
	AccessToken string `json:"access_token"`
}

type TimeInterval struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type DaySchedule map[string][]TimeInterval

type Schedule struct {
	ID         int         `json:"id"`
	Year       int         `json:"year"`
	Month      int         `json:"month"`
	Cabinet    string      `json:"cabinet"`
	DoctorID   int         `json:"doctor"`
	PerDayData DaySchedule `json:"perDayData"`
}

type GetScheduleResponse struct {
	Schedules []Schedule `json:"rasps"`
	Count     string     `json:"count"`
	AtPage    int        `json:"atPage"`
	MaxPage   int        `json:"maxPage"`
	Response  int        `json:"response"`
}

func GetSchedule(request *GetScheduleRequest) (*GetScheduleResponse, *utils.UserErrorResponse) {
	logger := utils.NewLogger("clients:GetSchedule")

	baseURL := BaseURL + "rasp/find"

	params := url.Values{}
	params.Add("access_token", request.AccessToken)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	logger.Infof("получение списка расписания по URL: %s", fullURL)

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

	var response GetScheduleResponse
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
