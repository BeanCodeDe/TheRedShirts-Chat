package core

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (core CoreFacade) startCleanUp() {
	log.Info("Start auto cleanup of messages")
	s := gocron.NewScheduler(time.UTC)

	s.Every(60).Seconds().Do(func() {
		correlationId := uuid.NewString()
		logger := log.WithFields(log.Fields{
			"Scavenger": correlationId,
		})

		tx, err := core.db.StartTransaction()
		if err != nil {
			logger.Warnf("something went wrong while creating transaction: %v", err)
			return
		}
		defer tx.Rollback()

		if err := tx.DeleteMessages(time.Now().Add(-30 * time.Second)); err != nil {
			log.Warn("Error while deleting old messages: %v", err)
			return
		}
		tx.Commit()
	})

	//s.StartAsync()
}
