// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package togo

import (
	"database/sql"
	"time"
)

type Purchase struct {
	ID          int64
	WagerID     int64
	BuyingPrice sql.NullString
	BoughtAt    time.Time
}

type Wager struct {
	ID                  int64
	TotalWagerValue     int32
	Odds                int32
	SellingPercentage   int32
	SellingPrice        string
	CurrentSellingPrice string
	PercentageSold      sql.NullInt32
	AmountSold          sql.NullString
	PlacedAt            time.Time
}
