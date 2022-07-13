package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/duongcongtoai/toytoytoy/internal/services"
	"github.com/duongcongtoai/toytoytoy/sqlc/togo"
	"github.com/labstack/echo/v4"
)

type APIConf struct {
	JWTSecret string
	TokenExp  time.Duration
}

type HttpAPI struct {
	conf        APIConf
	wagerSvc    services.WagerSvc
	purchaseSvc services.PurchaseSvc
}

func BindAPI(conf APIConf, e *echo.Echo, wagerSvc services.WagerSvc, purchaseSvc services.PurchaseSvc) *HttpAPI {
	result := &HttpAPI{
		conf:        conf,
		wagerSvc:    wagerSvc,
		purchaseSvc: purchaseSvc,
	}
	e.POST("/buy/:wager_id", result.BuyWager)
	e.POST("/wagers", result.PlaceWager)
	e.GET("/wagers", result.ListWagers)
	return result
}

type PlaceWagerReq struct {
	TotalWagerValue   int     `json:"total_wager_value"`
	Odds              int     `json:"odds"`
	SellingPercentage int     `json:"selling_percentage"`
	SellingPrice      float64 `json:"selling_price"`
}

func (h *HttpAPI) PlaceWager(ctx echo.Context) error {
	var req PlaceWagerReq
	err := ctx.Bind(&req)
	if err != nil {
		return err
	}
	wager, err := h.wagerSvc.PlaceWager(ctx.Request().Context(), DTOToWagerModel(req))
	if err != nil {
		return err
	}
	wagerDTO, err := ModelToWagerItem(wager)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, wagerDTO)
}

func (h *HttpAPI) ListWagers(ctx echo.Context) error {
	page := ctx.QueryParam("page")
	limit := ctx.QueryParam("limit")
	intpage, err := strconv.Atoi(page)
	if err != nil {
		return err
	}
	intlimit, err := strconv.Atoi(limit)
	if err != nil {
		return err
	}
	wagers, err := h.wagerSvc.GetWagers(ctx.Request().Context(), int32(intlimit), int32(intpage))
	if err != nil {
		return err
	}
	wagerDtos, err := ModelToWagerItems(wagers)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, wagerDtos)
}
func DTOToWagerModel(item PlaceWagerReq) togo.Wager {
	return togo.Wager{
		TotalWagerValue:   int32(item.TotalWagerValue),
		Odds:              int32(item.Odds),
		SellingPrice:      fmt.Sprintf("%.2f", item.SellingPrice),
		SellingPercentage: int32(item.SellingPercentage),
	}
}
func ModelToWagerItem(item togo.Wager) (WagerItem, error) {
	fcursellingPrice, err := strconv.ParseFloat(item.CurrentSellingPrice, 64)
	if err != nil {
		return WagerItem{}, err
	}
	dto := WagerItem{
		ID:                  int64(item.ID),
		TotalWagerValue:     int(item.TotalWagerValue),
		Odds:                int(item.Odds),
		SellingPercentage:   int(item.SellingPercentage),
		CurrentSellingPrice: fcursellingPrice,
		PlacedAt:            item.PlacedAt.Format(time.RFC3339),
	}

	if item.PercentageSold.Valid {
		intPercentageSold := int(item.PercentageSold.Int32)
		dto.PercentageSold = &intPercentageSold
	}
	if item.AmountSold.Valid {
		famountSold, err := strconv.ParseFloat(item.AmountSold.String, 64)
		if err != nil {
			return WagerItem{}, err
		}
		dto.AmountSold = &famountSold
	}
	return dto, nil
}
func ModelToWagerItems(wagers []togo.Wager) ([]WagerItem, error) {
	ret := make([]WagerItem, 0, len(wagers))
	for _, item := range wagers {
		fcursellingPrice, err := strconv.ParseFloat(item.CurrentSellingPrice, 64)
		if err != nil {
			return nil, err
		}
		dto := WagerItem{
			ID:                  int64(item.ID),
			TotalWagerValue:     int(item.TotalWagerValue),
			Odds:                int(item.Odds),
			SellingPercentage:   int(item.SellingPercentage),
			CurrentSellingPrice: fcursellingPrice,
			PlacedAt:            item.PlacedAt.Format(time.RFC3339),
		}

		if item.PercentageSold.Valid {
			intPercentageSold := int(item.PercentageSold.Int32)
			dto.PercentageSold = &intPercentageSold
		}
		if item.AmountSold.Valid {
			famountSold, err := strconv.ParseFloat(item.AmountSold.String, 64)
			if err != nil {
				return nil, err
			}
			dto.AmountSold = &famountSold
		}
		ret = append(ret, dto)
	}
	return ret, nil
}

type WagerItem struct {
	ID                  int64    `json:"id"`
	TotalWagerValue     int      `json:"total_wager_value"`
	Odds                int      `json:"odds"`
	SellingPercentage   int      `json:"selling_percentage"`
	SellingPrice        float64  `json:"selling_price"`
	CurrentSellingPrice float64  `json:"current_selling_price"`
	PercentageSold      *int     `json:"percentage_sold"`
	AmountSold          *float64 `json:"amount_sold"`
	PlacedAt            string   `json:"placed_at"`
}

type BuyWagerReq struct {
	BuyingPrice float64 `json:"buying_price"`
}
type BuyWagerRes struct {
	ID          int64   `json:"id"`
	WagerID     int64   `json:"wager_id"`
	BuyingPrice float64 `json:"buying_price"`
	BoughtAt    string  `json:"bought_at"`
}

func (h *HttpAPI) BuyWager(ctx echo.Context) error {
	wagerID := ctx.Param("wager_id")
	intWagerID, err := strconv.Atoi(wagerID)
	if err != nil {
		return err
	}
	var req BuyWagerReq
	err = ctx.Bind(&req)
	if err != nil {
		return err
	}
	purchase, err := h.purchaseSvc.BuyWager(ctx.Request().Context(), int64(intWagerID), req.BuyingPrice)
	if err != nil {
		return err
	}
	dto, err := ModelToPurchaseDTO(purchase)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, dto)
}
func ModelToPurchaseDTO(model togo.Purchase) (BuyWagerRes, error) {
	fbuyingPrice, err := strconv.ParseFloat(model.BuyingPrice.String, 64)
	if err != nil {
		return BuyWagerRes{}, err
	}
	return BuyWagerRes{
		ID:          model.ID,
		WagerID:     model.WagerID,
		BuyingPrice: fbuyingPrice,
		BoughtAt:    model.BoughtAt.Format(time.RFC3339),
	}, nil
}
