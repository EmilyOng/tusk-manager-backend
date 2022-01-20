package utils

import (
	"net/http"

	"github.com/EmilyOng/cvwo/backend/models"
)

func MakeResponseCode(response models.Response) int {
	if len(response.Error) == 0 {
		return http.StatusOK
	}
	return http.StatusInternalServerError
}

func MakeResponseErr(err models.ErrorCode) models.Response {
	return models.Response{Error: string(err)}
}
