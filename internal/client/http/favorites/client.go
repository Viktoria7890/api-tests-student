package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	favoritesModels "api-tests-template/internal/managers/favorites/models"

	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpPostFavorite добавляет объявление в избранное авторизованного пользователя
func HttpPostFavorite(t *testing.T, token string, request favoritesModels.AddFavoriteRequest) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).
		End().Response
}

// HttpGetFavorites возвращает список избранного авторизованного пользователя
func HttpGetFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Expect(t).
		End().Response
}

// HttpDeleteFavorite удаляет объявление из избранного авторизованного пользователя
func HttpDeleteFavorite(t *testing.T, token, advertisementId string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.FavoriteByIdPath(advertisementId)).
		Expect(t).
		End().Response
}
