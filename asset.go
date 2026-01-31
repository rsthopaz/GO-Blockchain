package main

import (
	"time"

	"github.com/google/uuid"
)

type Asset struct {
  ID        uuid.UUID `db:"id" json:"id"`
  Code      string    `db:"asset_code" json:"asset_code"`
  Name      string    `db:"asset_name" json:"asset_name"`
  Category  string    `db:"category" json:"category"`
  Value     float64   `db:"value" json:"value"`
  Status    string    `db:"status" json:"status"`
  CreatedAt time.Time `db:"created_at" json:"created_at"`
  UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
