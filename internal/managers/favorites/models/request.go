package models

// AddFavoriteRequest тело запроса добавления объявления в избранное
type AddFavoriteRequest struct {
	AdvertisementId string `json:"advertisement_id"`
}
