package services

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/duongcongtoai/toytoytoy/sqlc/togo"
)

type DBX interface {
	togo.DBTX
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Begin() (*sql.Tx, error)
}

type WagerSvc struct {
	// we accept leaking implementation details, because eventually, we require underneath storage implementation to support transaction
	db        DBX
	wagerRepo WagerRepo
}
type PurchaseSvc struct {
	db           DBX
	wagerRepo    WagerRepo
	purchaseRepo PurchaseRepo
}

type WagerRepo interface {
	GetWager(ctx context.Context, db togo.DBTX, id int64) (togo.Wager, error)
	GetWagers(ctx context.Context, db togo.DBTX, arg togo.GetWagersParams) ([]togo.Wager, error)
	CreateWager(ctx context.Context, db togo.DBTX, arg togo.CreateWagerParams) (togo.Wager, error)
	UpdateWager(ctx context.Context, db togo.DBTX, arg togo.UpdateWagerParams) error
}

type PurchaseRepo interface {
	CreatePurchase(ctx context.Context, arg togo.CreatePurchaseParams) (togo.Purchase, error)
}

func (s *WagerSvc) GetWager(ctx context.Context, id int64) (togo.Wager, error) {
	return s.wagerRepo.GetWager(ctx, s.db, id)
}
func (s *WagerSvc) GetWagers(ctx context.Context, limit int32, offset int32) ([]togo.Wager, error) {
	return s.wagerRepo.GetWagers(ctx, s.db, togo.GetWagersParams{Limit: limit, Offset: offset})
}
func (s *WagerSvc) PlaceWager(ctx context.Context, wager togo.Wager) (togo.Wager, error) {
	wager, err := s.wagerRepo.CreateWager(ctx, s.db, togo.CreateWagerParams{
		TotalWagerValue:     wager.TotalWagerValue,
		Odds:                wager.Odds,
		SellingPercentage:   wager.SellingPercentage,
		SellingPrice:        wager.SellingPrice,
		CurrentSellingPrice: wager.SellingPrice,
		PlacedAt:            time.Now(),
	})
	return wager, err
}

func (s *PurchaseSvc) BuyWager(ctx context.Context, wagerID int64, buyingPrice float64) (togo.Purchase, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return togo.Purchase{}, err
	}
	var purchase togo.Purchase
	err = ExecWithTx(ctx, tx, func(ctx context.Context, tx *sql.Tx) error {
		wager, err := s.wagerRepo.GetWager(ctx, tx, wagerID)
		if err != nil {
			return err
		}
		curSellingPrice, err := strconv.ParseFloat(wager.CurrentSellingPrice, 64)
		if err != nil {
			return err
		}

		curSellingPrice = curSellingPrice - buyingPrice
		if curSellingPrice < 0.0 {
			return fmt.Errorf("buying price exceed selling price")
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
		purchase, err = s.purchaseRepo.CreatePurchase(ctx, togo.CreatePurchaseParams{WagerID: wagerID,
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

func ExecWithTx(ctx context.Context, tx *sql.Tx, f func(context.Context, *sql.Tx) error) error {
	var err error
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	err = f(ctx, tx)
	return err
}
