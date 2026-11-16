package advertisement

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisement"
	"api-tests-template/internal/helpers/retry"
	"api-tests-template/internal/managers/advertisement/models"
)

// CreateAdvertisement создаёт объявление со всеми переданными параметрами и проверяет ожидаемый статус-код
func CreateAdvertisement(t *testing.T, token, title, description string, price, quantity int,
	photos []models.Photo, expectedStatusCode int) string {
	rawResponse := retry.On503(func() *http.Response {
		return advertisement.HttpPostAdvertisement(t, token, title, description, price, quantity, photos)
	})
	defer rawResponse.Body.Close()

	require.Equalf(t, expectedStatusCode, rawResponse.StatusCode,
		"HTTP status code должен быть %d", expectedStatusCode)

	responseBytes, err := io.ReadAll(rawResponse.Body)
	require.NoError(t, err)

	return string(responseBytes)
}

// GetAdvertisementById возвращает объявление по id и проверяет ожидаемый статус-код
func GetAdvertisementById(t *testing.T, id string, expectedStatusCode int) string {
	rawResponse := retry.On503(func() *http.Response {
		return advertisement.HttpGetAdvertisementById(t, id)
	})
	defer rawResponse.Body.Close()

	require.Equalf(t, expectedStatusCode, rawResponse.StatusCode,
		"HTTP status code должен быть %d", expectedStatusCode)

	responseBytes, err := io.ReadAll(rawResponse.Body)
	require.NoError(t, err)

	return string(responseBytes)
}

// DeleteAdvertisementById удаляет объявление по id и проверяет ожидаемый статус-код
func DeleteAdvertisementById(t *testing.T, token, id string, expectedStatusCode int) string {
	rawResponse := retry.On503(func() *http.Response {
		return advertisement.HttpDeleteAdvertisementById(t, token, id)
	})
	defer rawResponse.Body.Close()

	require.Equalf(t, expectedStatusCode, rawResponse.StatusCode,
		"HTTP status code должен быть %d", expectedStatusCode)

	responseBytes, err := io.ReadAll(rawResponse.Body)
	require.NoError(t, err)

	return string(responseBytes)
}
