package data

import (
    "time"
    "database/sql"
)

type Token struct {
    Hash []byte `json:"token"`
    UserID int64 `json:"-"`
    Expiry time.Time `json:"expiry"`
}

type TokenModel struct {
    DB *sql.DB
}
