package favorites

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/helpers/retry"
	"api-tests-template/internal/managers/favorites/models"
)

// AddFavorite добавляет объявление в избранное и проверяет ожидаемый статус-код
func AddFavorite(t *testing.T, token, advertisementId string, expectedStatusCode int) string {
	rawResponse := retry.On503(func() *http.Response {
		return favorites.HttpPostFavorite(t, token, models.AddFavoriteRequest{AdvertisementId: advertisementId})
	})
	defer rawResponse.Body.Close()

	require.Equalf(t, expectedStatusCode, rawResponse.StatusCode,
		"HTTP status code должен быть %d", expectedStatusCode)

	responseBytes, err := io.ReadAll(rawResponse.Body)
	require.NoError(t, err)

	return string(responseBytes)
}

// GetFavorites возвращает список избранного и проверяет ожидаемый статус-код
func GetFavorites(t *testing.T, token string, expectedStatusCode int) string {
	rawResponse := retry.On503(func() *http.Response {
		return favorites.HttpGetFavorites(t, token)
	})
	defer rawResponse.Body.Close()

	require.Equalf(t, expectedStatusCode, rawResponse.StatusCode,
		"HTTP status code должен быть %d", expectedStatusCode)

	responseBytes, err := io.ReadAll(rawResponse.Body)
	require.NoError(t, err)

	return string(responseBytes)
}

// DeleteFavorite удаляет объявление из избранного и проверяет ожидаемый статус-код
func DeleteFavorite(t *testing.T, token, advertisementId string, expectedStatusCode int) string {
	rawResponse := retry.On503(func() *http.Response {
		return favorites.HttpDeleteFavorite(t, token, advertisementId)
	})
	defer rawResponse.Body.Close()

	require.Equalf(t, expectedStatusCode, rawResponse.StatusCode,
		"HTTP status code должен быть %d", expectedStatusCode)

	responseBytes, err := io.ReadAll(rawResponse.Body)
	require.NoError(t, err)

	return string(responseBytes)
}
