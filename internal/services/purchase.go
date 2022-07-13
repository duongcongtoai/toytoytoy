package services

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/duongcongtoai/toytoytoy/internal/util"
	"github.com/duongcongtoai/toytoytoy/sqlc/togo"
)

func NewPurchaseSvc(db DBX, wagerRepo WagerRepo, purchaseRepo PurchaseRepo) *PurchaseSvc {
	return &PurchaseSvc{db: db, wagerRepo: wagerRepo, purchaseRepo: purchaseRepo}
}

type PurchaseSvc struct {
	db           DBX
	wagerRepo    WagerRepo
	purchaseRepo PurchaseRepo
}

type PurchaseRepo interface {
	CreatePurchase(ctx context.Context, db togo.DBTX, arg togo.CreatePurchaseParams) (togo.Purchase, error)
}

func (s *PurchaseSvc) BuyWager(ctx context.Context, wagerID int64, buyingPrice float64) (togo.Purchase, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return togo.Purchase{}, err
	}
	var purchase togo.Purchase
	err = util.ExecWithTx(ctx, tx, func(ctx context.Context, tx *sql.Tx) error {
		wager, err := s.wagerRepo.GetWagerForUpdate(ctx, tx, wagerID)
		if err != nil {
			return err
		}
		curSellingPrice, err := strconv.ParseFloat(wager.CurrentSellingPrice, 64)
		if err != nil {
			return err
		}

		curSellingPrice = curSellingPrice - buyingPrice
		if curSellingPrice < 0.0 {
			return ErrBuyingPriceExceedSellingPrice
		}
		amountSold := 0.0
		if wager.AmountSold.Valid {
			amountSold, err = strconv.ParseFloat(wager.AmountSold.String, 64)
			if err != nil {
				return err
			}
		}

		amountSold += buyingPrice
		originalSellingPrice, err := strconv.ParseFloat(wager.SellingPrice, 64)
		if err != nil {
			return err
		}
		percentageSold := math.RoundToEven(amountSold / originalSellingPrice * 100)
		err = s.wagerRepo.UpdateWager(ctx, tx, togo.UpdateWagerParams{
			CurrentSellingPrice: fmt.Sprintf("%.2f", curSellingPrice),
			AmountSold:          sql.NullString{String: fmt.Sprintf("%.2f", amountSold), Valid: true},
			PercentageSold:      sql.NullInt32{Valid: true, Int32: int32(percentageSold)},
			ID:                  wagerID,
		})
		if err != nil {
			return err
		}
		purchase, err = s.purchaseRepo.CreatePurchase(ctx, tx, togo.CreatePurchaseParams{WagerID: wagerID,
			BuyingPrice: sql.NullString{Valid: true, String: fmt.Sprintf("%.2f", buyingPrice)},
			BoughtAt:    time.Now(),
		})
		if err != nil {
			return err
		}
		return err
	})
	return purchase, err
}
