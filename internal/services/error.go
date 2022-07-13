package services

import "net/http"

type AppError struct {
	Code int
	Desc string
}

func (a AppError) Error() string {
	return a.Desc
}

var (
	ErrInvalidSellingPrice = AppError{
		Code: http.StatusBadRequest,
		Desc: "INVALID SELLING PRICE",
	}
	ErrBuyingPriceExceedSellingPrice = AppError{
		Code: http.StatusBadRequest,
		Desc: "BUYING PRICE EXCEEDS SELLING PRICE",
	}
)

func InvalidReqErr(detail string) error {
	return AppError{
		Code: http.StatusBadRequest,
		Desc: detail,
	}

}
