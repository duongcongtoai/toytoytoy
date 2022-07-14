package storage

import (
	"context"

	"github.com/duongcongtoai/toytoytoy/sqlc/togo"
	_ "github.com/go-sql-driver/mysql"
)

type WagerRepo struct {
}
type PurchaseRepo struct {
}

func (s *WagerRepo) GetWagerForUpdate(ctx context.Context, db togo.DBTX, id int64) (togo.Wager, error) {
	return togo.New(db).GetWagerForUpdate(ctx, id)
}

func (s *WagerRepo) GetWager(ctx context.Context, db togo.DBTX, id int64) (togo.Wager, error) {
	return togo.New(db).GetWager(ctx, id)
}
func (s *WagerRepo) GetWagers(ctx context.Context, db togo.DBTX, arg togo.GetWagersParams) ([]togo.Wager, error) {
	return togo.New(db).GetWagers(ctx, arg)
}

func (s *WagerRepo) CreateWager(ctx context.Context, db togo.DBTX, arg togo.CreateWagerParams) (togo.Wager, error) {
	ret, err := togo.New(db).CreateWager(ctx, arg)
	if err != nil {
		return togo.Wager{}, err
	}
	id, err := ret.LastInsertId()
	if err != nil {
		return togo.Wager{}, err
	}
	return togo.Wager{
		ID:                  id,
		TotalWagerValue:     arg.TotalWagerValue,
		Odds:                arg.Odds,
		SellingPercentage:   arg.SellingPercentage,
		SellingPrice:        arg.SellingPrice,
		CurrentSellingPrice: arg.CurrentSellingPrice,
		PercentageSold:      arg.PercentageSold,
		AmountSold:          arg.AmountSold,
		PlacedAt:            arg.PlacedAt,
	}, nil
}

func (s *WagerRepo) UpdateWager(ctx context.Context, db togo.DBTX, arg togo.UpdateWagerParams) error {
	_, err := togo.New(db).UpdateWager(ctx, arg)
	return err
}

func (s *PurchaseRepo) CreatePurchase(ctx context.Context, db togo.DBTX, arg togo.CreatePurchaseParams) (togo.Purchase, error) {
	ret, err := togo.New(db).CreatePurchase(ctx, arg)
	if err != nil {
		return togo.Purchase{}, nil
	}
	id, err := ret.LastInsertId()
	if err != nil {
		return togo.Purchase{}, err
	}
	return togo.Purchase{
		ID:          id,
		WagerID:     arg.WagerID,
		BuyingPrice: arg.BuyingPrice,
		BoughtAt:    arg.BoughtAt,
	}, nil
}
