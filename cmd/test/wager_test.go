package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/duongcongtoai/toytoytoy/internal/services"
	transhttp "github.com/duongcongtoai/toytoytoy/internal/transport/http"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestBuyWagerConcurrently(t *testing.T) {
	wagerReq := `
{
    "total_wager_value": 20,
    "odds": 7,
    "selling_percentage": 50,
    "selling_price": 20.0
}`
	res, err := http.Post(fmt.Sprintf("http://%s/wagers", getServerAddr()), "application/json", strings.NewReader(wagerReq))
	assert.NoError(t, err)
	defer res.Body.Close()
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	var resItem = transhttp.WagerItem{}
	raw, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(raw, &resItem)
	assert.NoError(t, err)
	createdWager := resItem.ID
	// spawn 25 goroutines,each purchasese 1.0
	buyReq := `{"buying_price": 1.0}`
	errGr := &errgroup.Group{}
	var totalErr int32
	for i := 0; i < 25; i++ {
		errGr.Go(func() error {
			res, err := http.Post(fmt.Sprintf("http://%s/buy/%d", getServerAddr(), createdWager), "application/json", strings.NewReader(buyReq))
			assert.NoError(t, err)

			defer res.Body.Close()
			if res.StatusCode != http.StatusCreated {
				raw, _ := ioutil.ReadAll(res.Body)
				var cont = map[string]string{}
				json.Unmarshal(raw, &cont)

				// concurrent access to update the same wager should be locked, and there should always be 5 failed req
				if cont["error"] == services.ErrBuyingPriceExceedSellingPrice.Desc {
					atomic.AddInt32(&totalErr, 1)
				}
			}
			return nil
		})
	}
	errGr.Wait()
	assert.Equal(t, int32(5), totalErr)
}

func TestBuyWager(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		name         string
		buyingPrice  float64
		expectedCode int
		expectedErr  string
	}
	commonWagerReq := `
{
    "total_wager_value": 20,
    "odds": 7,
    "selling_percentage": 50,
    "selling_price": 20.0
}
	`
	tcases := []TestCase{
		{
			name:         "success",
			buyingPrice:  0.5,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid buying price",
			buyingPrice:  21.0,
			expectedCode: http.StatusBadRequest,
			expectedErr:  "INVALID BUYING PRICE",
		},
		/** TODO
		- invalid wager_id
		- negative buying price
		- ...
		**/
	}
	for _, item := range tcases {
		t.Run(item.name, func(t *testing.T) {
			// create wager
			res, err := http.Post(fmt.Sprintf("http://%s/wagers", getServerAddr()), "application/json", strings.NewReader(commonWagerReq))
			assert.NoError(t, err)
			res.Body.Close()
			assert.Equal(t, item.expectedCode, res.StatusCode)

			// purchase wager

		})

	}
}

func TestCreateWager(t *testing.T) {
	type TestCase struct {
		name               string
		createReq          string
		expectedCreateRes  string
		expectedCreateCode int
	}
	t.Parallel()
	tcases := []TestCase{
		{
			name: "success",
			createReq: `
{
    "total_wager_value": 10,
    "odds": 7,
    "selling_percentage": 5,
    "selling_price": 1.0
}`,
			expectedCreateRes: `
{
    "id": "",
    "total_wager_value": 10,
    "odds": 7,
    "selling_percentage": 5,
    "selling_price": 1.0,
    "current_selling_price": 1.0,
    "percentage_sold": null,
    "amount_sold": null,
    "placed_at": "" 
}`,
			expectedCreateCode: http.StatusCreated,
		},
		{
			name: "invalid selling price",
			createReq: `
{
    "total_wager_value": 10,
    "odds": 7,
    "selling_percentage": 50,
    "selling_price": 4.0
}`,
			expectedCreateRes:  `{"error" : "INVALID SELLING PRICE"}`,
			expectedCreateCode: http.StatusBadRequest,
		},
		{
			name: "invalid selling price",
			createReq: `
{
    "total_wager_value": 10,
    "odds": 7,
    "selling_percentage": 50,
    "selling_price": -1.0
}`,
			expectedCreateRes:  `{"error" : "INVALID SELLING PRICE"}`,
			expectedCreateCode: http.StatusBadRequest,
		},
	}
	for _, item := range tcases {
		t.Run(item.name, func(t *testing.T) {
			res, err := http.Post(fmt.Sprintf("http://%s/wagers", getServerAddr()), "application/json", strings.NewReader(item.createReq))
			assert.NoError(t, err)
			defer res.Body.Close()
			rawRes, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			// id and placed_at is undeterministic
			equal, diff, err := compJsonStr(string(rawRes), item.expectedCreateRes, "id", "placed_at")
			assert.NoError(t, err)
			assert.True(t, equal, diff)
			assert.Equal(t, item.expectedCreateCode, res.StatusCode)

		})

	}
}

func compJsonStr(a, b string, ignoredFields ...string) (bool, string, error) {
	mapA := map[string]interface{}{}
	mapB := map[string]interface{}{}
	err := json.Unmarshal([]byte(a), &mapA)
	if err != nil {
		return false, "", err
	}
	err = json.Unmarshal([]byte(b), &mapB)
	if err != nil {
		return false, "", err
	}
	for _, ignored := range ignoredFields {
		delete(mapA, ignored)
		delete(mapB, ignored)
	}
	return cmp.Equal(mapA, mapB), cmp.Diff(mapA, mapB), nil
}
