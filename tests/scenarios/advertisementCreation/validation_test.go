package advertisementCreation

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	apiResponse "api-tests-template/internal/helpers/api-response"
	"api-tests-template/internal/managers/advertisements"
	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"
)

func TestRejectInvalidAdvertisement(t *testing.T) {
	base.SetupSuite()
	seller := fixtures.LoginSeller(t)
	const titleLimit = 50
	const descriptionLimit = 500

	cases := []struct {
		name    string
		field   string
		value   string
		omit    bool
		message string
	}{
		{"без названия", fixtures.TitleField, "", true, "title is required"},
		{"без описания", fixtures.DescriptionField, "", true, "description is required"},
		{"без фотографии", "", "", true, "at least 1 photo is required"},
		{"пустое название", fixtures.TitleField, "", false, "title is required"},
		{"пустое описание", fixtures.DescriptionField, "", false, "description is required"},
		{"слишком длинное название", fixtures.TitleField, strings.Repeat("x", titleLimit+1), false,
			fmt.Sprintf("title must not exceed %d characters provided %d", titleLimit, titleLimit+1)},
		{"слишком длинное описание", fixtures.DescriptionField, strings.Repeat("x", descriptionLimit+1), false, "description must be between 1 and 500 characters"},
		{"отрицательная цена", fixtures.PriceField, "-1", false, "price must be between 0 and 1000000 inclusive"},
		{"цена не число", fixtures.PriceField, "bad", false, "cannot convert string value bad to float"},
		{"нулевое количество", fixtures.QuantityField, "0", false, "quantity must be between 1 and 100 inclusive"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := fixtures.NewAdvertisementRequest(t)
			switch {
			case tc.field == "":
				request.Photo = nil
			case tc.omit:
				delete(request.Fields, tc.field)
			default:
				request.Fields[tc.field] = tc.value
			}
			t.Log("Отправляем форму с одним нарушением")
			got := advertisements.RejectCreateAdvertisement(t, seller.Token, request, http.StatusBadRequest)
			require.Equal(t, apiResponse.Error{Code: fixtures.InvalidRequestCode, Message: tc.message}, got)
		})
	}
}
