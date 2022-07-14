package services

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/bxcodec/faker/v3"
	mockcommon "github.com/duongcongtoai/toytoytoy/mock/common"
	mockserv "github.com/duongcongtoai/toytoytoy/mock/services"
	"github.com/duongcongtoai/toytoytoy/sqlc/togo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func RandomWager(id int64) togo.Wager {
	var ret togo.Wager
	err := faker.FakeData(&ret)
	if err != nil {
		panic(err)
	}
	ret.ID = int64(id)
	return ret
}

func TestBuyWager(t *testing.T) {
	wagerRepo := mockserv.NewWagerRepo(t)
	purchaseRepo := mockserv.NewPurchaseRepo(t)
	db := mockcommon.NewDBX(t)
	mockTx := mockcommon.NewTx(t)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)
	db.On("BeginTx", mock.Anything, mock.Anything).Return(mockTx, nil)
	svc := PurchaseSvc{db: db, wagerRepo: wagerRepo, purchaseRepo: purchaseRepo}

	var (
		nonExistWagerID int64 = 0
		newWagerID      int64 = 1
		oldWagerID      int64 = 2
		newWagerEnt           = RandomWager(newWagerID)
		oldWagerEnt           = RandomWager(oldWagerID)

		validBuyingPrice     = 9.0
		exceedingBuyingPrice = 11.0
	)

	// never bought
	newWagerEnt.SellingPrice = "10.0"
	newWagerEnt.CurrentSellingPrice = "10.0"
	newWagerEnt.AmountSold = sql.NullString{}

	// bought
	oldWagerEnt.SellingPrice = "20.0"
	oldWagerEnt.CurrentSellingPrice = "15.0"
	oldWagerEnt.AmountSold = sql.NullString{Valid: true, String: "5.0"}
	oldWagerEnt.SellingPercentage = 25

	wagerRepo.On("GetWagerForUpdate", mock.Anything, mockTx, nonExistWagerID).
		Return(togo.Wager{}, fmt.Errorf("non exist wager"))
	wagerRepo.On("GetWagerForUpdate", mock.Anything, mockTx, newWagerID).
		Return(newWagerEnt, nil)
	wagerRepo.On("GetWagerForUpdate", mock.Anything, mockTx, oldWagerID).
		Return(oldWagerEnt, nil)

	t.Run("non existing wager", func(t *testing.T) {
		_, err := svc.BuyWager(context.Background(), nonExistWagerID, 1.0)
		assert.Equal(t, err, fmt.Errorf("non exist wager"))
	})
	t.Run("exceeding buying price", func(t *testing.T) {
		_, err := svc.BuyWager(context.Background(), newWagerID, exceedingBuyingPrice)
		assert.Equal(t, err, ErrBuyingPriceExceedSellingPrice)
	})
	t.Run("negative buying price", func(t *testing.T) {
		_, err := svc.BuyWager(context.Background(), newWagerID, -1.0)
		assert.Equal(t, err, ErrInvalidBuyingPrice)
	})
	// TODO: reorganize so it is less boilerplating
	t.Run("success first buy", func(t *testing.T) {
		createPurchaseParam := togo.CreatePurchaseParams{
			WagerID:     newWagerID,
			BuyingPrice: sql.NullString{String: fmt.Sprintf("%.2f", validBuyingPrice), Valid: true},
		}
		updateWagerParam := togo.UpdateWagerParams{
			ID:                  newWagerID,
			CurrentSellingPrice: "1.00",
			PercentageSold:      sql.NullInt32{Int32: 90, Valid: true},
			AmountSold:          sql.NullString{String: "9.00", Valid: true},
		}
		returnedPurchase := togo.Purchase{
			WagerID:     newWagerID,
			BuyingPrice: sql.NullString{String: "9.00", Valid: true},
		}
		wagerRepo.On("UpdateWager", mock.Anything, mockTx, updateWagerParam).
			Return(nil)
		purchaseRepo.On("CreatePurchase", mock.Anything, mockTx, mock.MatchedBy(func(param togo.CreatePurchaseParams) bool {
			// time is underterministic
			createPurchaseParam.BoughtAt = param.BoughtAt
			return reflect.DeepEqual(createPurchaseParam, param)
		})).Return(returnedPurchase, nil)

		purchase, err := svc.BuyWager(context.Background(), newWagerID, validBuyingPrice)
		assert.NoError(t, err)
		assert.Equal(t, returnedPurchase, purchase)
	})

	t.Run("success second buy", func(t *testing.T) {
		createPurchaseParam := togo.CreatePurchaseParams{
			WagerID:     oldWagerID,
			BuyingPrice: sql.NullString{String: fmt.Sprintf("%.2f", validBuyingPrice), Valid: true},
		}
		updateWagerParam := togo.UpdateWagerParams{
			ID:                  oldWagerID,
			CurrentSellingPrice: "6.00",
			PercentageSold:      sql.NullInt32{Int32: 70, Valid: true},
			AmountSold:          sql.NullString{String: "14.00", Valid: true},
		}
		returnedPurchase := togo.Purchase{
			WagerID:     oldWagerID,
			BuyingPrice: sql.NullString{String: "9.00", Valid: true},
		}
		wagerRepo.On("UpdateWager", mock.Anything, mockTx, updateWagerParam).
			Return(nil)
		purchaseRepo.On("CreatePurchase", mock.Anything, mockTx, mock.MatchedBy(func(param togo.CreatePurchaseParams) bool {
			// time is underterministic
			createPurchaseParam.BoughtAt = param.BoughtAt
			return reflect.DeepEqual(createPurchaseParam, param)
		})).Return(returnedPurchase, nil)

		purchase, err := svc.BuyWager(context.Background(), oldWagerID, validBuyingPrice)
		assert.NoError(t, err)
		assert.Equal(t, returnedPurchase, purchase)
	})

}
