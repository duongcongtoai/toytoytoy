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
		validWagerID    int64 = 1
		wagerEnt              = RandomWager(validWagerID)

		validBuyingPrice     = 9.0
		exceedingBuyingPrice = 11.0
	)

	wagerEnt.SellingPrice = "10.0"
	wagerEnt.CurrentSellingPrice = "10.0"

	wagerRepo.On("GetWagerForUpdate", mock.Anything, mockTx, nonExistWagerID).
		Return(togo.Wager{}, fmt.Errorf("non exist wager"))
	wagerRepo.On("GetWagerForUpdate", mock.Anything, mockTx, validWagerID).
		Return(wagerEnt, nil)

	t.Run("non existing wager", func(t *testing.T) {
		_, err := svc.BuyWager(context.Background(), nonExistWagerID, 1.0)
		assert.Equal(t, err, fmt.Errorf("non exist wager"))
	})
	t.Run("exceeding buying price", func(t *testing.T) {
		_, err := svc.BuyWager(context.Background(), validWagerID, exceedingBuyingPrice)
		assert.Equal(t, err, ErrBuyingPriceExceedSellingPrice)
	})
	t.Run("negative buying price", func(t *testing.T) {
		_, err := svc.BuyWager(context.Background(), validWagerID, -1.0)
		assert.Equal(t, err, ErrInvalidBuyingPrice)
	})

	t.Run("success", func(t *testing.T) {
		createPurchaseParam := togo.CreatePurchaseParams{
			WagerID:     validWagerID,
			BuyingPrice: sql.NullString{String: fmt.Sprintf("%.2f", validBuyingPrice), Valid: true},
		}
		updateWagerParam := togo.UpdateWagerParams{
			ID:                  validWagerID,
			CurrentSellingPrice: "1.00",
			PercentageSold:      sql.NullInt32{Int32: 90, Valid: true},
			AmountSold:          sql.NullString{String: "9.00", Valid: true},
		}
		returnedPurchase := togo.Purchase{
			WagerID:     validWagerID,
			BuyingPrice: sql.NullString{String: "9.00", Valid: true},
		}
		wagerRepo.On("UpdateWager", mock.Anything, mockTx, updateWagerParam).
			Return(nil)
		purchaseRepo.On("CreatePurchase", mock.Anything, mockTx, mock.MatchedBy(func(param togo.CreatePurchaseParams) bool {
			// time is underterministic
			createPurchaseParam.BoughtAt = param.BoughtAt
			return reflect.DeepEqual(createPurchaseParam, param)
		})).Return(returnedPurchase, nil)

		purchase, err := svc.BuyWager(context.Background(), validWagerID, validBuyingPrice)
		assert.NoError(t, err)
		assert.Equal(t, returnedPurchase, purchase)
	})

}
