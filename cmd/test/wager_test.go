package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	name               string
	createReq          string
	expectedCreateRes  string
	expectedCreateCode int
	buyReq             string
	expectedBuyRes     string
	expectedBuyCode    int
}

func TestCreateWager(t *testing.T) {
	t.Parallel()
	tcases := []TestCase{
		{
			name: "success",
			createReq: `
{
    "total_wager_value": 10,
    "odds": 7,
    "selling_percentage": 5,
    "selling_price": 1.0,
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
    "selling_price": 4.0,
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
    "selling_price": -1.0,
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
			equal, err := compJson(rawRes, item.expectedCreateRes, "id", "placed_at")
			assert.NoError(t, err)
			assert.True(t, equal)
		})

	}
}

func compJson(a, b interface{}, ignoredFields ...string) (bool, error) {
	rawa, err := json.Marshal(a)
	if err != nil {
		return false, err
	}
	rawb, err := json.Marshal(b)
	if err != nil {
		return false, err
	}
	mapA := map[string]interface{}{}
	mapB := map[string]interface{}{}
	err = json.Unmarshal(rawa, &mapA)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(rawb, &mapB)
	if err != nil {
		return false, err
	}
	for _, ignored := range ignoredFields {
		delete(mapA, ignored)
		delete(mapB, ignored)
	}
	return reflect.DeepEqual(mapA, mapB), nil
}
