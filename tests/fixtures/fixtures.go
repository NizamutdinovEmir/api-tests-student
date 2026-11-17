package fixtures

import (
	"bytes"
	"image"
	"image/png"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

const (
	MissingAdvertisementID = "00000000-0000-4000-8000-000000000000"
	InvalidToken           = "incorrect_token"
	MissingAuthMessage     = "Authorization header required"
	InvalidAuthMessage     = "Invalid or expired token"
	UnauthorizedCode       = "unauthorized"
	InvalidRequestCode     = "invalid_request"

	TitleField       = "title"
	DescriptionField = "description"
	PriceField       = "price"
	QuantityField    = "quantity"
)

// LoginSeller авторизует владельца тестовых объявлений.
func LoginSeller(t *testing.T) authModels.LoginOkResponse {
	t.Helper()
	t.Log("Подготовка: авторизуем продавца")
	return auth.Login(t, os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

// LoginBuyer авторизует отдельного покупателя для работы с избранным.
func LoginBuyer(t *testing.T) authModels.LoginOkResponse {
	t.Helper()
	t.Log("Подготовка: авторизуем покупателя")
	login, password := os.Getenv("TEST_BUYER_LOGIN"), os.Getenv("TEST_BUYER_PASSWORD")
	require.NotEmpty(t, login, "задайте TEST_BUYER_LOGIN в .env.override по инструкции README")
	require.NotEmpty(t, password, "задайте TEST_BUYER_PASSWORD в .env.override по инструкции README")
	return auth.Login(t, login, password)
}

// NewAdvertisementRequest готовит уникальное объявление со всеми полями и тестовым PNG.
func NewAdvertisementRequest(t *testing.T) models.CreateRequest {
	t.Helper()
	const photoSize = 32

	var photo bytes.Buffer
	require.NoError(t, png.Encode(&photo, image.NewRGBA(image.Rect(0, 0, photoSize, photoSize))), "создание тестового PNG")
	return models.CreateRequest{
		Fields: map[string]string{
			TitleField:       "API Test " + utils.RandomString(20),
			DescriptionField: "Объявление для автоматизированной проверки TestBoard",
			PriceField:       "199.5",
			QuantityField:    "3",
		},
		Photo: photo.Bytes(),
	}
}

// CreateAdvertisement создаёт объявление и возвращает функцию для его досрочного удаления.
func CreateAdvertisement(t *testing.T, sellerToken string, request models.CreateRequest) (models.Advertisement, func()) {
	t.Helper()
	t.Logf("Создаём объявление %q", request.Fields[TitleField])
	ad := advertisements.CreateAdvertisement(t, sellerToken, request)
	require.NotEmpty(t, ad.ID, "созданное объявление вернулось без ID")
	deleted := false
	deleteAd := func() {
		t.Helper()
		if deleted {
			return
		}
		t.Logf("Удаляем тестовое объявление %s", ad.ID)
		advertisements.DeleteAdvertisement(t, sellerToken, ad.ID, http.StatusNoContent)
		deleted = true
	}
	t.Cleanup(deleteAd)
	return ad, deleteAd
}
