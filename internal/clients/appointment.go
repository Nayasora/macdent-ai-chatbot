package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"macdent-ai-chatbot/internal/utils"
	"net/http"
	"net/url"
	"strconv"
)

type CreateAppointmentRequest struct {
	AccessToken          string `json:"access_token"`
	DoctorID             int    `json:"doctor"`
	PatientID            int    `json:"patient"`
	AppointmentDate      string `json:"date"`
	AppointmentStartTime string `json:"start"`
	AppointmentEndTime   string `json:"end"`
}

type AppointmentInfo struct {
	ID         int    `json:"id"`
	DoctorID   int    `json:"doctor"`
	PatientID  int    `json:"patient"`
	Date       string `json:"date"`
	StartTime  string `json:"start"`
	EndTime    string `json:"end"`
	Status     int    `json:"status"`
	Complaint  string `json:"zhaloba"`
	Comment    string `json:"comment"`
	IsFirst    bool   `json:"isFirst"`
	Cabinet    string `json:"cabinet"`
	ScheduleID string `json:"rasp"`
}

type CreateAppointmentResponse struct {
	Appointment AppointmentInfo `json:"zapis"`
	Response    int             `json:"response"`
}

func CreateAppointment(request CreateAppointmentRequest) (*CreateAppointmentResponse, *utils.UserErrorResponse) {
	logger := utils.NewLogger("clients:createAppointment")

	baseURL := BaseURL + "zapis/add"

	params := url.Values{}
	params.Add("access_token", request.AccessToken)
	params.Add("doctor", strconv.Itoa(request.DoctorID))
	params.Add("patient", strconv.Itoa(request.PatientID))
	params.Add("date", request.AppointmentDate)
	params.Add("start", request.AppointmentStartTime)
	params.Add("end", request.AppointmentEndTime)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	logger.Infof("создание записи на приём по URL: %s", fullURL)

	req, err := http.NewRequest(http.MethodPost, fullURL, nil)
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

	var response CreateAppointmentResponse
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
