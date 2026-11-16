package path

import "fmt"

const (
	AuthLoginPath        = "/auth/login"
	MyAdvertisementsPath = "/my/advertisements"
	AdvertisementPath    = "/advertisement"
	AdvertisementsPath   = "/advertisements"
	FavoritesPath        = "/favorites"
)

// AdvertisementPhotosPath возвращает путь ручки получения фотографий объявления по его id
func AdvertisementPhotosPath(advertisementId string) string {
	return fmt.Sprintf("/advertisements/%s/photos", advertisementId)
}

// FavoriteByIdPath возвращает путь ручки удаления объявления из избранного по его id
func FavoriteByIdPath(advertisementId string) string {
	return fmt.Sprintf("/favorites/%s", advertisementId)
}
