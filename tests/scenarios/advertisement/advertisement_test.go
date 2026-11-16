package advertisement

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/advertisement/models"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

// tooLongTitleLength длина title, заведомо превышающая допустимый на сайте лимит
const tooLongTitleLength = 300

type TestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// newValidPhoto возвращает валидную тестовую фотографию для прикрепления к объявлению
func newValidPhoto() models.Photo {
	return models.Photo{FileName: "photo.png", Content: utils.GenerateTestImage()}
}

// Positive: создание объявления со всеми возможными параметрами.
// Проверяем доступность объявления по id, в поиске и фотографии на сервере
func (s *TestSuite) TestCreateAdvertisementPositive() {
	title := "Тестовое объявление " + utils.RandomString(8)
	description := "Описание тестового объявления " + utils.RandomString(16)
	price := 1500
	quantity := 2

	var createBody string
	s.Run("Создаём объявление со всеми возможными параметрами", func() {
		createBody = advertisement.CreateAdvertisement(s.T(), s.loginData.Token, title, description, price, quantity,
			[]models.Photo{newValidPhoto()}, http.StatusCreated)
	})

	var adId string
	s.Run("Проверяем, что в ответе вернулись переданные параметры", func() {
		adId = gjson.Get(createBody, "id").String()
		require.NotEmpty(s.T(), adId, "id созданного объявления не должен быть пустым")
		require.Equal(s.T(), title, gjson.Get(createBody, "title").String(), "title не совпадает")
		require.Equal(s.T(), description, gjson.Get(createBody, "description").String(), "description не совпадает")
		require.EqualValues(s.T(), price, gjson.Get(createBody, "price").Int(), "price не совпадает")
		require.Len(s.T(), gjson.Get(createBody, "photos").Array(), 1, "ожидается одна загруженная фотография")
	})

	defer func() {
		base.Precondition(fmt.Sprintf("Удаляем созданное объявление %s", adId))
		advertisement.DeleteAdvertisementById(s.T(), s.loginData.Token, adId, http.StatusNoContent)
	}()

	s.Run("Проверяем, что объявление доступно через GET /advertisement по id и возвращает переданные параметры", func() {
		getBody := advertisement.GetAdvertisementById(s.T(), adId, http.StatusOK)
		require.Equal(s.T(), title, gjson.Get(getBody, "title").String(), "title не совпадает")
		require.Equal(s.T(), description, gjson.Get(getBody, "description").String(), "description не совпадает")
		require.EqualValues(s.T(), price, gjson.Get(getBody, "price").Int(), "price не совпадает")
		require.EqualValues(s.T(), quantity, gjson.Get(getBody, "quantity").Int(), "quantity не совпадает")
	})

	s.Run("Проверяем, что объявление находится в поиске по полному вхождению названия", func() {
		searchBody := advertisements.SearchAdvertisements(s.T(), title, 5, http.StatusOK)
		items := gjson.Get(searchBody, "items").Array()

		var found bool
		for _, item := range items {
			if gjson.Get(item.String(), "id").String() == adId {
				found = true
				require.Equal(s.T(), title, gjson.Get(item.String(), "title").String(), "title найденного объявления не совпадает")
			}
		}
		require.True(s.T(), found, "созданное объявление не найдено в выдаче поиска по названию")
	})

	s.Run("Проверяем, что фотография загружена на сервер и доступна по ручке фотографий объявления", func() {
		photosBody := advertisements.GetAdvertisementPhotos(s.T(), s.loginData.Token, adId, http.StatusOK)
		photos := gjson.Parse(photosBody).Array()

		require.Len(s.T(), photos, 1, "ожидается ровно одна фотография")
		require.NotEmpty(s.T(), gjson.Get(photos[0].String(), "url").String(), "url фотографии не должен быть пустым")
	})
}

// Negative: попытка создать объявление с пустым title
func (s *TestSuite) TestCreateAdvertisementEmptyTitleNegative() {
	base.Precondition("Пытаемся создать объявление с пустым title")

	body := advertisement.CreateAdvertisement(s.T(), s.loginData.Token, "", "Описание объявления", 100, 1,
		[]models.Photo{newValidPhoto()}, http.StatusBadRequest)

	require.NotEmpty(s.T(), gjson.Get(body, "error").String(), "в ответе должно быть указано поле error")
}

// Negative: попытка создать объявление с отрицательной ценой
func (s *TestSuite) TestCreateAdvertisementNegativePriceNegative() {
	base.Precondition("Пытаемся создать объявление с отрицательной ценой")

	body := advertisement.CreateAdvertisement(s.T(), s.loginData.Token, "Объявление "+utils.RandomString(6),
		"Описание объявления", -100, 1, []models.Photo{newValidPhoto()}, http.StatusBadRequest)

	require.NotEmpty(s.T(), gjson.Get(body, "error").String(), "в ответе должно быть указано поле error")
}

// Negative: попытка создать объявление со слишком длинным title
func (s *TestSuite) TestCreateAdvertisementTooLongTitleNegative() {
	base.Precondition("Пытаемся создать объявление со слишком длинным title")
	longTitle := strings.Repeat("а", tooLongTitleLength)

	body := advertisement.CreateAdvertisement(s.T(), s.loginData.Token, longTitle, "Описание объявления", 100, 1,
		[]models.Photo{newValidPhoto()}, http.StatusBadRequest)

	require.NotEmpty(s.T(), gjson.Get(body, "error").String(), "в ответе должно быть указано поле error")
}

// Negative: запрос на создание объявления без токена авторизации
func (s *TestSuite) TestCreateAdvertisementWithoutTokenNegative() {
	base.Precondition("Пытаемся создать объявление без токена авторизации")

	body := advertisement.CreateAdvertisement(s.T(), "", "Объявление "+utils.RandomString(6), "Описание объявления", 100, 1,
		[]models.Photo{newValidPhoto()}, http.StatusUnauthorized)

	require.NotEmpty(s.T(), gjson.Get(body, "error").String(), "в ответе должно быть указано поле error")
}

// Negative: запрос на создание объявления с некорректным токеном авторизации
func (s *TestSuite) TestCreateAdvertisementIncorrectTokenNegative() {
	base.Precondition("Пытаемся создать объявление с некорректным токеном авторизации")

	body := advertisement.CreateAdvertisement(s.T(), "incorrect_token", "Объявление "+utils.RandomString(6),
		"Описание объявления", 100, 1, []models.Photo{newValidPhoto()}, http.StatusUnauthorized)

	require.NotEmpty(s.T(), gjson.Get(body, "error").String(), "в ответе должно быть указано поле error")
}
