package advertisementCreation

import (
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"
)

func TestCreateAdvertisementWithAllFields(t *testing.T) {
	base.SetupSuite()
	seller := fixtures.LoginSeller(t)
	request := fixtures.NewAdvertisementRequest(t)
	created, _ := fixtures.CreateAdvertisement(t, seller.Token, request)

	t.Run("данные объявления совпадают с отправленной формой", func(t *testing.T) {
		ad := advertisements.GetAdvertisement(t, seller.Token, created.ID)
		require.Equal(t, created.ID, ad.ID, "ID объявления")
		require.Equal(t, seller.User.Id, ad.UserID, "ID владельца объявления")
		require.Equal(t, request.Fields[fixtures.TitleField], ad.Title, "название объявления")
		require.Equal(t, request.Fields[fixtures.DescriptionField], ad.Description, "описание объявления")
		require.Equal(t, request.Fields[fixtures.PriceField], strconv.FormatFloat(ad.Price, 'f', -1, 64), "цена объявления")
		require.Equal(t, request.Fields[fixtures.QuantityField], strconv.Itoa(ad.Quantity), "количество товара")
	})

	t.Run("объявление находится по полному названию", func(t *testing.T) {
		title := request.Fields[fixtures.TitleField]
		list := advertisements.GetAdvertisementsList(t, seller.Token, title)
		found := slices.ContainsFunc(list.Items, func(ad models.Advertisement) bool {
			return ad.ID == created.ID && ad.Title == title
		})
		require.True(t, found, "поиск по %q не вернул объявление %s", title, created.ID)
	})

	t.Run("фотография доступна через ручку фотографий", func(t *testing.T) {
		photos := advertisements.GetAdvertisementPhotos(t, seller.Token, created.ID)
		require.Len(t, created.Photos, 1, "ответ создания должен содержать загруженную фотографию")
		require.Len(t, photos, 1, "ручка фотографий должна вернуть одну фотографию")
		require.NotEmpty(t, photos[0].ID, "ID загруженной фотографии")
		require.NotEmpty(t, photos[0].URL, "URL загруженной фотографии")
		require.Equal(t, models.Photo{ID: created.Photos[0].ID, URL: created.Photos[0].URL, SortOrder: 0},
			photos[0], "ручка фотографий должна вернуть фото из ответа создания")
		require.Equal(t, request.Photo, advertisements.DownloadAdvertisementPhoto(t, photos[0].URL),
			"по URL из ручки фотографий должен скачиваться отправленный PNG")
	})
}
