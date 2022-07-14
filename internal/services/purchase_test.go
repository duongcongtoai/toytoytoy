package services

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/bxcodec/faker/v3"
	mockserv "github.com/duongcongtoai/toytoytoy/mock/services"
	"github.com/duongcongtoai/toytoytoy/sqlc/togo"
	"github.com/stretchr/testify/mock"
)

func RandomWager(id int) togo.Wager {
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
	db := mockserv.NewDBX(t)
	db.On("BeginTx", mock.Anything, mock.Anything).Return(&sql.Tx{}, nil)
	mockTx := mock.AnythingOfType("sql.Tx")
	svc := PurchaseSvc{db: db, wagerRepo: wagerRepo, purchaseRepo: purchaseRepo}

	nonExistWagerID := 0
	validWagerID := 1
	wagerEnt := RandomWager(validWagerID)
	wagerEnt.SellingPrice = "10.0"
	validBuyingPrice := "9.0"
	exceedingBuyingPrice := "11.0"
	wagerRepo.On("GetWagerForUpdate", mock.Anything, mockTx, nonExistWagerID).
		Return(togo.Wager{}, fmt.Errorf("non exist wager"))
	wagerRepo.On("GetWagerForUpdate", mock.Anything, mockTx, validWagerID).
		Return(wagerEnt, nil)
	wagerRepo.On("UpdateWager", mock.Anything, mockTx, mock.Anything).
		Return(nil)
	purchaseRepo.On("CreatePurchase", mock.Anything, mockTx, mock.Anything)

	type TestCase struct {
		buyingPrice     float64
		sellingPriceGap float64
		amountSold      float64
		percentageSold  int
	}

}
