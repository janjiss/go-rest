package paginator

import (
	"encoding/base64"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IDer interface {
	GetID() uuid.UUID
}

func EncodeCursor[T IDer](object T) string {
	id := object.GetID()
	return base64.StdEncoding.EncodeToString(id[:])
}

func decodeCursor(cursor string) (uuid.UUID, error) {
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.FromBytes(decoded)
}

func FetchNextPage[T IDer](db *gorm.DB, cursor string, limit int) ([]T, error) {
	var items []T

	if cursor != "" {
		lastID, err := decodeCursor(cursor)
		if err != nil {
			return nil, err
		}
		db = db.Where("id > ?", lastID)
	}

	err := db.Order("id ASC").Limit(limit).Find(&items).Error
	return items, err
}

func FetchPrevPage[T IDer](db *gorm.DB, cursor string, limit int) ([]T, error) {
	var items []T

	if cursor != "" {
		firstID, err := decodeCursor(cursor)
		if err != nil {
			return nil, err
		}
		db = db.Where("id < ?", firstID)
	}

	err := db.Order("id DESC").Limit(limit).Find(&items).Error
	return items, err
}
