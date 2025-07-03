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
		select {
		case <-session.Context().Done():
			return session.Context().Err()
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
			}).Trace("message consumed. Running handler ...")

			if err := h.handler.Handle(session.Context(), message); err != nil {
				if errors.Is(errors.Cause(err), ErrMarkAcked) {
					session.MarkMessage(message, "")
				}

				return errors.Wrap(err, "error running handler")
			}
			log.Trace("handler completed without an error. Marking message ...")

			session.MarkMessage(message, "")
		}
	}
}

func (h *consumerGroupHandler) Close() error {
	log.WithFields(log.Fields{
		"component": "ConsumerGroupHandler",
	}).Trace("Close() called")

	return nil
}
