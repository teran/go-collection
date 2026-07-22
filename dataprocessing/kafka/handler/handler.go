package handler

import (
	"context"
	"math/rand"
	"time"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	maxRetries     = 3
	baseRetryDelay = 1 * time.Second
	maxRetryDelay  = 30 * time.Second
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

				if err := h.handleWithRetry(session, message); err != nil {
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

func (h *consumerGroupHandler) handleWithRetry(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) error {
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := h.handler.Handle(session.Context(), message)
		if err == nil {
			return nil
		}

		if errors.Is(errors.Cause(err), ErrMarkAcked) {
			log.WithError(err).WithFields(log.Fields{
				"component": "ConsumerGroupHandler",
				"attempt":   attempt + 1,
			}).Debug("handler returned ErrMarkAcked. Marking message ...")
			session.MarkMessage(message, "")
			return nil
		}

		if attempt < maxRetries-1 {
			delay := backoffDelay(attempt)
			log.WithError(err).WithFields(log.Fields{
				"component":   "ConsumerGroupHandler",
				"attempt":     attempt + 1,
				"max_retries": maxRetries,
				"retry_in":    delay.String(),
			}).Warning("handler returned error. Retrying ...")

			select {
			case <-session.Context().Done():
				return errors.Wrap(session.Context().Err(), "context cancelled during retry")
			case <-time.After(delay):
			}
		}
	}

	log.WithFields(log.Fields{
		"component":   "ConsumerGroupHandler",
		"max_retries": maxRetries,
	}).Error("handler exhausted retries. Not marking message")
	return errors.Errorf("handler failed after %d attempts", maxRetries)
}

func backoffDelay(attempt int) time.Duration {
	delay := baseRetryDelay * (1 << attempt) // 1s, 2s, 4s
	if delay > maxRetryDelay {
		delay = maxRetryDelay
	}
	// Add jitter: ±25%
	jitter := time.Duration(float64(delay) * (0.75 + 0.5*rand.Float64())) //nolint:gosec // jitter doesn't need crypto/rand
	return jitter
}

func (h *consumerGroupHandler) Close() error {
	log.WithFields(log.Fields{
		"component": "ConsumerGroupHandler",
	}).Trace("Close() called")

	return nil
}
