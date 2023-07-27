package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"
)

type Token struct {
    Hash []byte `json:"-"`
    Plaintext string `json:"token"`
    UserID int64 `json:"-"`
    Expiry time.Time `json:"expiry"`
}

type TokenModel struct {
    DB *sql.DB
}

// generateToken generates a token for purposes of authentication for a user.
func generateToken(userID int64, lifetime time.Duration) (*Token, error) {
    token := &Token{
        UserID: userID,
        Expiry: time.Now().Add(lifetime),
    }

    randomBytes := make([]byte, 16)

    _, err := rand.Read(randomBytes)
    if err != nil {
        return nil, err
    }

    token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
    hash := sha256.Sum256([]byte(token.Plaintext))
    token.Hash = hash[:]

    return token, nil
}

// AddToken adds an auth token for the user to the database.
func (m TokenModel) AddToken(userID int64, lifetime time.Duration) error {
    token, err := generateToken(userID, lifetime)
    if err != nil {
        return err
    }
    
    stmt := `INSERT INTO token (hash, userID, expiry)
    VALUES ($1, $2, $3)`

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    _, err = m.DB.ExecContext(ctx, stmt, token.Hash, token.UserID, token.Expiry)
    return err
}
