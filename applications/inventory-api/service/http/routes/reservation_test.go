package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const (
	payload = `{
		"lines": [
			{
				"product": "PIPR-JACKET-SIZM",
				"quantity": 1
			},
			{
				"product": "PIPR-JACKET-SIZL",
				"quantity": 2
			}
		]
	}`
)

type createReservationUseCase struct{}

func (usecase *createReservationUseCase) CreateReservation(request ReservationRequest) (ReservationResponse, error) {
	reservation := ReservationResponse{
		ID:        "test",
		CreatedAt: "date",
		Lines:     request.Lines,
		Status:    ReservationReserved,
	}
	return reservation, nil
}

func TestCreateReservation(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/reservation", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	useCase := &createReservationUseCase{}
	handleCreateReservationHTTPRequest := CreateReservation(useCase)

	// Assertions
	if assert.NoError(t, handleCreateReservationHTTPRequest(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code, "status code should be 201 Created")

		// Assert that response is built on ApplyUseCase
		responseBodyString := rec.Body.String()
		var responseReservation ReservationResponse
		json.Unmarshal([]byte(responseBodyString), &responseReservation)

		applyCallResult, _ := useCase.CreateReservation(ReservationRequest{
			Lines: []OrderLine{{
				ProductID:       "PIPR-JACKET-SIZM",
				ProductQuantity: 1,
			}, {
				ProductID:       "PIPR-JACKET-SIZL",
				ProductQuantity: 2,
			}},
		})
		assert.Equal(t, applyCallResult, responseReservation, "response should equal to apply callback result")
	}
}
