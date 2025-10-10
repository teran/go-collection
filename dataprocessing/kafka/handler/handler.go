package handler

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	_ sarama.ConsumerGroupHandler = (*consumerGroupHandler)(nil)

	ErrMarkAcked = errors.New("skip message")
)

type Handler interface {
	Handle(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type consumerGroupHandler struct {
	handler Handler
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	log.WithFields(log.Fields{
		"component": "ConsumerGroupHandler",
	}).Trace("Setup() called")

	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.WithFields(log.Fields{
		"component": "ConsumerGroupHandler",
	}).Trace("Cleanup() called")

	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log.WithFields(log.Fields{
		"component": "ConsumerGroupHandler",
	}).Trace("ConsumeClaim() called")

	for {
		if err := func() error {
			defer session.Commit()

			select {
			case <-session.Context().Done():
				return errors.Wrap(session.Context().Err(), "error received from context")
			case message, ok := <-claim.Messages():
				if !ok {
					log.Warn("message channel was closed")
					return nil
				}

				log.WithFields(log.Fields{
					"topic":     message.Topic,
					"timestamp": message.Timestamp.Format(time.RFC3339),
					"offset":    message.Offset,
					"partition": message.Partition,
					"length":    len(message.Value),
				}).Debug("message consumed. Running handler ...")

				if err := h.handler.Handle(session.Context(), message); err != nil {
					if errors.Is(errors.Cause(err), ErrMarkAcked) {
						log.WithError(err).WithFields(log.Fields{
							"component": "ConsumerGroupHandler",
						}).Debug("handler returned ErrMarkAcked. Marking message ...")
						session.MarkMessage(message, "")
					}

					log.WithError(err).Error("error running handler. Not marking message")
					return errors.Wrap(err, "error running handler")
				}
				log.WithFields(log.Fields{
					"component": "ConsumerGroupHandler",
				}).Debug("handler completed without an error. Marking message ...")

				session.MarkMessage(message, "")
				return nil
			}
		}(); err != nil {
			log.WithError(err).Error("error running consumer group handler")
			return errors.Wrap(err, "error running consumer group handler")
		}
	}
}

func (h *consumerGroupHandler) Close() error {
	log.WithFields(log.Fields{
		"component": "ConsumerGroupHandler",
	}).Trace("Close() called")

	return nil
}
