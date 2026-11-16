package favorites

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/advertisement/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

// В сценариях избранного участвуют два пользователя:
// пользователь A (TEST_LOGIN/TEST_PASSWORD) создаёт объявления,
// пользователь B (TEST_LOGIN_SECOND/TEST_PASSWORD_SECOND) добавляет их в избранное.
// Это нужно, потому что добавить в избранное можно только чужое объявление (своё — 403),
// а удалить объявление может только его владелец.
type TestSuite struct {
	suite.Suite
	userAData authModels.LoginOkResponse
	userBData authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя A (владелец объявлений) с кредами из переменных окружения")
	s.userAData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))

	base.Precondition("Авторизация пользователя B (добавляет объявления в избранное) с кредами из переменных окружения")
	s.userBData = auth.Login(s.T(), os.Getenv("TEST_LOGIN_SECOND"), os.Getenv("TEST_PASSWORD_SECOND"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// createAdvertisementAsUserA создаёт объявление от имени пользователя A и возвращает его id
func (s *TestSuite) createAdvertisementAsUserA(title string) string {
	photo := models.Photo{FileName: "photo.png", Content: utils.GenerateTestImage()}

	body := advertisement.CreateAdvertisement(s.T(), s.userAData.Token, title, "Описание объявления", 500, 1,
		[]models.Photo{photo}, http.StatusCreated)

	id := gjson.Get(body, "id").String()
	require.NotEmpty(s.T(), id, "id созданного объявления не должен быть пустым")

	return id
}

// favoritesContainAdvertisement проверяет, есть ли объявление с данным id в списке избранного
func favoritesContainAdvertisement(favoritesBody, advertisementId string) bool {
	for _, item := range gjson.Get(favoritesBody, "items").Array() {
		if gjson.Get(item.String(), "advertisement.id").String() == advertisementId {
			return true
		}
	}
	return false
}

// Positive: полный цикл работы с избранным — добавление, удаление из избранного,
// а также каскадное исчезновение из избранного при удалении самого объявления
func (s *TestSuite) TestFavoritesFullCyclePositive() {
	title := "Объявление для избранного " + utils.RandomString(8)

	var adId string
	s.Run("Пользователь A создаёт объявление", func() {
		adId = s.createAdvertisementAsUserA(title)
	})

	adDeleted := false
	defer func() {
		if !adDeleted {
			base.Precondition(fmt.Sprintf("Удаляем созданное объявление %s (очистка после теста)", adId))
			advertisement.DeleteAdvertisementById(s.T(), s.userAData.Token, adId, http.StatusNoContent)
		}
	}()

	s.Run("Пользователь B добавляет чужое объявление в избранное", func() {
		body := favorites.AddFavorite(s.T(), s.userBData.Token, adId, http.StatusCreated)
		require.Equal(s.T(), adId, gjson.Get(body, "advertisement.id").String(),
			"id добавленного в избранное объявления не совпадает")
	})

	s.Run("Проверяем, что объявление появилось в списке избранного пользователя B", func() {
		body := favorites.GetFavorites(s.T(), s.userBData.Token, http.StatusOK)
		require.True(s.T(), favoritesContainAdvertisement(body, adId), "объявление не найдено в избранном")
	})

	s.Run("Пользователь B удаляет объявление из избранного", func() {
		favorites.DeleteFavorite(s.T(), s.userBData.Token, adId, http.StatusNoContent)
	})

	s.Run("Проверяем, что объявление пропало из избранного после удаления из избранного", func() {
		body := favorites.GetFavorites(s.T(), s.userBData.Token, http.StatusOK)
		require.False(s.T(), favoritesContainAdvertisement(body, adId), "объявление всё ещё в избранном после удаления")
	})

	s.Run("Пользователь B снова добавляет объявление в избранное", func() {
		favorites.AddFavorite(s.T(), s.userBData.Token, adId, http.StatusCreated)
	})

	s.Run("Пользователь A удаляет само объявление", func() {
		advertisement.DeleteAdvertisementById(s.T(), s.userAData.Token, adId, http.StatusNoContent)
		adDeleted = true
	})

	s.Run("Проверяем, что удалённое объявление пропало из избранного пользователя B", func() {
		body := favorites.GetFavorites(s.T(), s.userBData.Token, http.StatusOK)
		require.False(s.T(), favoritesContainAdvertisement(body, adId), "удалённое объявление всё ещё в избранном")
	})
}

// Negative: попытка добавить и удалить избранное по несуществующему ad_id.
// Реальный сервер возвращает 404 Not Found для обеих операций (хотя swagger его не документирует).
func (s *TestSuite) TestFavoriteNonexistentAdvertisementNegative() {
	nonexistentId := utils.RandomUUID()

	s.Run("Пытаемся добавить в избранное несуществующее объявление", func() {
		body := favorites.AddFavorite(s.T(), s.userBData.Token, nonexistentId, http.StatusNotFound)
		require.NotEmpty(s.T(), gjson.Get(body, "error").String(), "в ответе должно быть указано поле error")
	})

	s.Run("Пытаемся удалить из избранного несуществующее объявление", func() {
		body := favorites.DeleteFavorite(s.T(), s.userBData.Token, nonexistentId, http.StatusNotFound)
		require.NotEmpty(s.T(), gjson.Get(body, "error").String(), "в ответе должно быть указано поле error")
	})
}

// Negative: запросы к ручкам избранного без авторизации
func (s *TestSuite) TestFavoritesWithoutAuthNegative() {
	anyId := utils.RandomUUID()

	s.Run("Получение избранного без токена авторизации", func() {
		favorites.GetFavorites(s.T(), "", http.StatusUnauthorized)
	})

	s.Run("Добавление в избранное без токена авторизации", func() {
		favorites.AddFavorite(s.T(), "", anyId, http.StatusUnauthorized)
	})

	s.Run("Удаление из избранного без токена авторизации", func() {
		favorites.DeleteFavorite(s.T(), "", anyId, http.StatusUnauthorized)
	})
}

// Negative: нельзя добавить в избранное собственное объявление
func (s *TestSuite) TestFavoriteOwnAdvertisementNegative() {
	title := "Собственное объявление " + utils.RandomString(8)

	var adId string
	s.Run("Пользователь A создаёт объявление", func() {
		adId = s.createAdvertisementAsUserA(title)
	})

	defer func() {
		base.Precondition(fmt.Sprintf("Удаляем созданное объявление %s", adId))
		advertisement.DeleteAdvertisementById(s.T(), s.userAData.Token, adId, http.StatusNoContent)
	}()

	s.Run("Пользователь A пытается добавить своё же объявление в избранное", func() {
		body := favorites.AddFavorite(s.T(), s.userAData.Token, adId, http.StatusForbidden)
		require.NotEmpty(s.T(), gjson.Get(body, "error").String(), "в ответе должно быть указано поле error")
	})
}

// Negative: нельзя повторно добавить одно и то же объявление в избранное
func (s *TestSuite) TestFavoriteDuplicateNegative() {
	title := "Объявление для дубля избранного " + utils.RandomString(8)

	var adId string
	s.Run("Пользователь A создаёт объявление", func() {
		adId = s.createAdvertisementAsUserA(title)
	})

	defer func() {
		base.Precondition(fmt.Sprintf("Убираем объявление %s из избранного и удаляем его", adId))
		favorites.DeleteFavorite(s.T(), s.userBData.Token, adId, http.StatusNoContent)
		advertisement.DeleteAdvertisementById(s.T(), s.userAData.Token, adId, http.StatusNoContent)
	}()

	s.Run("Пользователь B добавляет объявление в избранное первый раз", func() {
		favorites.AddFavorite(s.T(), s.userBData.Token, adId, http.StatusCreated)
	})

	s.Run("Пользователь B повторно добавляет то же объявление в избранное", func() {
		body := favorites.AddFavorite(s.T(), s.userBData.Token, adId, http.StatusConflict)
		require.NotEmpty(s.T(), gjson.Get(body, "error").String(), "в ответе должно быть указано поле error")
	})
}
