package services

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/duongcongtoai/toytoytoy/sqlc/togo"
)

type DBX interface {
	togo.DBTX
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Begin() (*sql.Tx, error)
}

func NewWagerSvc(db DBX, wagerRepo WagerRepo) *WagerSvc {
	return &WagerSvc{db: db, wagerRepo: wagerRepo}
}

type WagerSvc struct {
	// we accept leaking implementation details, because eventually, we require underneath storage implementation to support transaction
	db        DBX
	wagerRepo WagerRepo
}

type WagerRepo interface {
	GetWager(ctx context.Context, db togo.DBTX, id int64) (togo.Wager, error)
	GetWagerForUpdate(ctx context.Context, db togo.DBTX, id int64) (togo.Wager, error)
	GetWagers(ctx context.Context, db togo.DBTX, arg togo.GetWagersParams) ([]togo.Wager, error)
	CreateWager(ctx context.Context, db togo.DBTX, arg togo.CreateWagerParams) (togo.Wager, error)
	UpdateWager(ctx context.Context, db togo.DBTX, arg togo.UpdateWagerParams) error
}

func (s *WagerSvc) GetWager(ctx context.Context, id int64) (togo.Wager, error) {
	return s.wagerRepo.GetWager(ctx, s.db, id)
}
func (s *WagerSvc) GetWagers(ctx context.Context, limit int32, offset int32) ([]togo.Wager, error) {
	return s.wagerRepo.GetWagers(ctx, s.db, togo.GetWagersParams{Limit: limit, Offset: offset})
}
func (s *WagerSvc) PlaceWager(ctx context.Context, wager togo.Wager) (togo.Wager, error) {
	err := ValidateWager(wager)
	if err != nil {
		return togo.Wager{}, err
	}
	wager, err = s.wagerRepo.CreateWager(ctx, s.db, togo.CreateWagerParams{
		TotalWagerValue:     wager.TotalWagerValue,
		Odds:                wager.Odds,
		SellingPercentage:   wager.SellingPercentage,
		SellingPrice:        wager.SellingPrice,
		CurrentSellingPrice: wager.SellingPrice,
		PlacedAt:            time.Now(),
	})
	return wager, err
}

func ValidateWager(wager togo.Wager) error {
	fsellingPrice, err := strconv.ParseFloat(wager.SellingPrice, 64)
	if err != nil {
		return InvalidReqErr(err.Error())
	}
	priceFloor := float64(wager.TotalWagerValue) * float64(wager.SellingPercentage) / 100
	if fsellingPrice <= priceFloor {
		return ErrInvalidSellingPrice
	}
	if fsellingPrice < 0.0 {
		return ErrInvalidSellingPrice
	}
	return nil
}
