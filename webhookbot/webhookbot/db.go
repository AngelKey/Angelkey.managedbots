package webhookbot

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"

	"github.com/keybase/managed-bots/base"
)

type DB struct {
	*base.DB
}

func NewDB(db *sql.DB) *DB {
	return &DB{
		DB: base.NewDB(db),
	}
}

func (d *DB) makeID(name, convID string) (string, error) {
	secret, err := base.RandBytes(16)
	if err != nil {
		return "", err
	}
	cdat, err := hex.DecodeString(convID)
	if err != nil {
		return "", err
	}
	h := hmac.New(sha256.New, secret)
	_, _ = h.Write(cdat)
	_, _ = h.Write([]byte(name))
	return base.URLEncoder().EncodeToString(h.Sum(nil)[:20]), nil
}

func (d *DB) Create(name, convID string) (string, error) {
	id, err := d.makeID(name, convID)
	if err != nil {
		return "", err
	}
	err = d.RunTxn(func(tx *sql.Tx) error {
		if _, err := tx.Exec(`
			INSERT INTO hooks
			(id, name, conv_id)
			VALUES
			(?, ?, ?)
		`, id, name, convID); err != nil {
			return err
		}
		return nil
	})
	return id, err
}

func (d *DB) GetHook(id string) (res webhook, err error) {
	row := d.DB.QueryRow(`
		SELECT conv_id, name FROM hooks WHERE id = ?
	`, id)
	if err := row.Scan(&res.convID, &res.name); err != nil {
		return res, err
	}
	return res, nil
}

type webhook struct {
	id     string
	convID string
	name   string
}

func (d *DB) List(convID string) (res []webhook, err error) {
	rows, err := d.DB.Query(`
		SELECT id, name FROM hooks WHERE conv_id = ?
	`, convID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var hook webhook
		hook.convID = convID
		if err := rows.Scan(&hook.id, &hook.name); err != nil {
			return res, err
		}
		res = append(res, hook)
	}
	return res, nil
}

func (d *DB) Remove(name, convID string) error {
	return d.RunTxn(func(tx *sql.Tx) error {
		_, err := tx.Exec(`
			DELETE FROM hooks WHERE conv_id = ? AND name = ?
		`, convID, name)
		return err
	})
}
