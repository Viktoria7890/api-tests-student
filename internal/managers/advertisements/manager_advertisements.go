package advertisements

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/helpers/retry"
)

// SearchAdvertisements выполняет поиск объявлений по названию и проверяет ожидаемый статус-код
func SearchAdvertisements(t *testing.T, search string, limit int, expectedStatusCode int) string {
	rawResponse := retry.On503(func() *http.Response {
		return advertisements.HttpGetAdvertisementsSearch(t, search, limit)
	})
	defer rawResponse.Body.Close()

	require.Equalf(t, expectedStatusCode, rawResponse.StatusCode,
		"HTTP status code должен быть %d", expectedStatusCode)

	responseBytes, err := io.ReadAll(rawResponse.Body)
	require.NoError(t, err)

	return string(responseBytes)
}

// GetAdvertisementPhotos возвращает список фотографий объявления и проверяет ожидаемый статус-код
func GetAdvertisementPhotos(t *testing.T, token, id string, expectedStatusCode int) string {
	rawResponse := retry.On503(func() *http.Response {
		return advertisements.HttpGetAdvertisementPhotos(t, token, id)
	})
	defer rawResponse.Body.Close()

	require.Equalf(t, expectedStatusCode, rawResponse.StatusCode,
		"HTTP status code должен быть %d", expectedStatusCode)

	responseBytes, err := io.ReadAll(rawResponse.Body)
	require.NoError(t, err)

	return string(responseBytes)
}
