package finam

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrNotFound = errors.New("404 Not Found")

// APIError структура ошибки в api
type APIError struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
	Status  int      `json:"-"` // статус HTTP ответа
}

func (e APIError) Error() string {
	return fmt.Sprintf("<APIError> httpStatus=%v, code=%d, msg=%s, details=%s", e.Status, e.Code, e.Message, e.Details)
	//return fmt.Sprintf("<APIError> httpStatus=%v, code=%d, msg=%s", e.Status, e.Code, e.Message)
}

func (e APIError) HTTPStatus() int {
	return e.Status
}

func makeAPIError(resp *Response) error {
	apiErr := new(APIError)
	apiErr.Status = resp.StatusCode
	// обработать ошибку StatusNotFound
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}

	err := resp.DecodeJSON(apiErr)
	if err != nil {
		log.Error("makeAPIError DecodeJSON", "err", err.Error())
		//apiErr.Code = resp.Status //strconv.Itoa(resp.StatusCode())
		apiErr.Message = http.StatusText(resp.StatusCode)
	}
	return apiErr

}

// func makeAPIError(resp *Response) error {
// 	apiErr := &APIError{}
// 	apiErr.Status = resp.StatusCode()
// 	// обработать ошибку StatusNotFound
// 	if resp.StatusCode() == http.StatusNotFound {
// 		return ErrNotFound
// 	}
// 	err := json.Unmarshal(resp.Body(), apiErr)
// 	if err != nil {
// 		log.Error("callAPI json.Unmarshal", "err", err.Error())
// 		apiErr.Code = resp.StatusCode()
// 		apiErr.Message = http.StatusText(resp.StatusCode())
// 	}
// 	return *apiErr
// }
