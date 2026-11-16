package advertisements

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"

	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpGetAdvertisementsSearch выполняет GET-запрос поиска объявлений по вхождению в название
func HttpGetAdvertisementsSearch(t *testing.T, search string, limit int) *http.Response {
	return apiRunner.GetRunner().Create().Get(path.AdvertisementsPath).
		Query("search", search).
		Query("limit", strconv.Itoa(limit)).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos выполняет GET-запрос получения фотографий объявления
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementPhotosPath(id)).
		Expect(t).
		End().Response
}
