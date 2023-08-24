package stores

import (
	"go.sia.tech/renterd/webhooks"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	dbWebhook struct {
		Model

		Module string `gorm:"uniqueIndex:idx_url_module_event;NOT NULL"`
		Event  string `gorm:"uniqueIndex:idx_url_module_event;NOT NULL"`
		URL    string `gorm:"uniqueIndex:idx_url_module_event;NOT NULL"`
	}
)

func (dbWebhook) TableName() string {
	return "webhooks"
}

func (s *SQLStore) DeleteHook(wb webhooks.Webhook) error {
	return s.retryTransaction(func(tx *gorm.DB) error {
		res := tx.Exec("DELETE FROM webhooks WHERE module = ? AND event = ? AND url = ?",
			wb.Module, wb.Event, wb.URL)
		if res.Error != nil {
			return res.Error
		} else if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (s *SQLStore) AddHook(wb webhooks.Webhook) error {
	return s.retryTransaction(func(tx *gorm.DB) error {
		res := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&dbWebhook{
			Module: wb.Module,
			Event:  wb.Event,
			URL:    wb.URL,
		})
		if res.Error != nil {
			return res.Error
		} else if res.RowsAffected == 0 {
			return ErrRecordExists
		}
		return nil
	})
}

func (s *SQLStore) Webhooks() ([]webhooks.Webhook, error) {
	var dbWebhooks []dbWebhook
	if err := s.db.Find(&dbWebhooks).Error; err != nil {
		return nil, err
	}
	var whs []webhooks.Webhook
	for _, wb := range dbWebhooks {
		whs = append(whs, webhooks.Webhook{
			Module: wb.Module,
			Event:  wb.Event,
			URL:    wb.URL,
		})
	}
	return whs, nil
}
