package handler

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Service interface {
	Run(ctx context.Context) error
}

type service struct {
	cg      sarama.ConsumerGroup
	topics  []string
	handler Handler
}

func New(cg sarama.ConsumerGroup, topics []string, handler Handler) Service {
	return &service{
		cg:      cg,
		topics:  topics,
		handler: handler,
	}
}

func (s *service) Run(ctx context.Context) error {
	log.WithFields(log.Fields{
		"component": "ConsumerGroupHandler",
	}).Trace("Run() called")

	cgh := &consumerGroupHandler{
		handler: s.handler,
	}

	for {
		log.WithFields(log.Fields{
			"component": "ConsumerGroupHandler",
		}).Trace("starting message consumption")

		if err := s.cg.Consume(ctx, s.topics, cgh); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return err
			}
		}

		if ctx.Err() != nil {
			if errors.Is(ctx.Err(), context.Canceled) {
				log.WithFields(log.Fields{
					"component": "ConsumerGroupHandler",
				}).Warn("context cancelled on message consumption")

				return nil
			}
			return errors.Wrap(ctx.Err(), "context error received")
		}
	}
}
