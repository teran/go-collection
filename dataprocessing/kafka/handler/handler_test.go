package handler

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/sync/errgroup"

	"github.com/teran/go-docker-testsuite/applications/kafka"
)

func init() {
	log.SetLevel(log.TraceLevel)
	sarama.Logger = log.StandardLogger()
}

const testTopicName = "test-topic"

func (s *handlerTestSuite) TestRoundtrip() {
	url, err := s.kafka.GetBrokerURL(s.ctx)
	s.Require().NoError(err)

	g, ctx := errgroup.WithContext(s.ctx)
	g.SetLimit(10)

	handlerMock := &testHandler{
		cancelFn: s.cancelFn,
	}
	handlerMock.On("Handle", testTopicName, []byte("test")).Return(nil).Once()
	defer handlerMock.AssertExpectations(s.T())

	g.Go(func() error {
		producer, err := sarama.NewSyncProducer([]string{url}, newKafkaConfig())
		if err != nil {
			return errors.Wrap(err, "error creating new producer")
		}

		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic:     testTopicName,
			Partition: 0,
			Value:     sarama.StringEncoder("test"),
		})
		if err != nil {
			return errors.Wrap(err, "error producing message")
		}

		log.WithFields(log.Fields{
			"partition": partition,
			"offset":    offset,
			"topic":     testTopicName,
		}).Warnf("message sent")

		return nil
	})

	g.Go(func() error {
		cg, err := sarama.NewConsumerGroup([]string{url}, "test-group", newKafkaConfig())
		if err != nil {
			return errors.Wrap(err, "error creating consumer group")
		}

		cgh := New(cg, []string{testTopicName}, handlerMock)
		if err = cgh.Run(ctx); err != nil {
			return errors.Wrap(err, "error running consumer group handler")
		}

		return nil
	})

	err = g.Wait()
	s.Require().NoError(err)
}

func (s *handlerTestSuite) TestRoundtrip_ServiceError() {
	url, err := s.kafka.GetBrokerURL(s.ctx)
	s.Require().NoError(err)

	g, ctx := errgroup.WithContext(s.ctx)
	g.SetLimit(10)

	handlerMock := &testHandler{
		cancelFn: s.cancelFn,
	}
	handlerMock.On("Handle", testTopicName, []byte("test #1")).Return(errors.New("blah")).Once()
	handlerMock.On("Handle", testTopicName, []byte("test #1")).Return(nil).Once()
	handlerMock.On("Handle", testTopicName, []byte("test #2")).Return(nil).Once()
	defer handlerMock.AssertExpectations(s.T())

	g.Go(func() error {
		producer, err := sarama.NewSyncProducer([]string{url}, newKafkaConfig())
		if err != nil {
			return errors.Wrap(err, "error creating new producer")
		}

		for i := 1; i < 3; i++ {
			partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
				Topic:     testTopicName,
				Partition: 0,
				Value:     sarama.StringEncoder("test #" + strconv.Itoa(i)),
			})
			if err != nil {
				return errors.Wrap(err, "error producing message")
			}

			log.WithFields(log.Fields{
				"partition": partition,
				"offset":    offset,
				"topic":     testTopicName,
			}).Warnf("message sent")
		}

		return nil
	})

	g.Go(func() error {
		cg, err := sarama.NewConsumerGroup([]string{url}, "test-group", newKafkaConfig())
		if err != nil {
			return errors.Wrap(err, "error creating consumer group")
		}

		cgh := New(cg, []string{testTopicName}, handlerMock)
		if err = cgh.Run(ctx); err != nil {
			return errors.Wrap(err, "error running consumer group handler")
		}

		return nil
	})

	err = g.Wait()
	s.Require().NoError(err)
}

func (s *handlerTestSuite) TestRoundtrip_ServiceErrorMarkAcked() {
	url, err := s.kafka.GetBrokerURL(s.ctx)
	s.Require().NoError(err)

	g, ctx := errgroup.WithContext(s.ctx)
	g.SetLimit(10)

	handlerMock := &testHandler{
		cancelFn: s.cancelFn,
	}
	handlerMock.On("Handle", testTopicName, []byte("test #1")).Return(ErrMarkAcked).Once()
	handlerMock.On("Handle", testTopicName, []byte("test #2")).Return(nil).Once()
	defer handlerMock.AssertExpectations(s.T())

	g.Go(func() error {
		producer, err := sarama.NewSyncProducer([]string{url}, newKafkaConfig())
		if err != nil {
			return errors.Wrap(err, "error creating new producer")
		}

		for i := 1; i < 3; i++ {
			partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
				Topic:     testTopicName,
				Partition: 0,
				Value:     sarama.StringEncoder("test #" + strconv.Itoa(i)),
			})
			if err != nil {
				return errors.Wrap(err, "error producing message")
			}

			log.WithFields(log.Fields{
				"partition": partition,
				"offset":    offset,
				"topic":     testTopicName,
			}).Warnf("message sent")
		}

		return nil
	})

	g.Go(func() error {
		cg, err := sarama.NewConsumerGroup([]string{url}, "test-group", newKafkaConfig())
		if err != nil {
			return errors.Wrap(err, "error creating consumer group")
		}

		cgh := New(cg, []string{testTopicName}, handlerMock)
		if err = cgh.Run(ctx); err != nil {
			return errors.Wrap(err, "error running consumer group handler")
		}

		return nil
	})

	err = g.Wait()
	s.Require().NoError(err)
}

// Definitions ...
type handlerTestSuite struct {
	suite.Suite

	ctx      context.Context
	cancelFn context.CancelFunc
	kafka    kafka.Kafka
}

func (s *handlerTestSuite) SetupTest() {
	var err error

	s.ctx, s.cancelFn = context.WithTimeout(s.T().Context(), 30*time.Second)

	s.kafka, err = kafka.New(s.ctx)
	s.Require().NoError(err)
}

func (s *handlerTestSuite) TearDownTest() {
	s.cancelFn()
	_ = s.kafka.Close(s.T().Context())
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, &handlerTestSuite{})
}

func newKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.ClientID = "go-collection-test-suite"

	return config
}

type testHandler struct {
	mock.Mock

	cancelFn context.CancelFunc
}

func (m *testHandler) Handle(_ context.Context, msg *sarama.ConsumerMessage) error {
	args := m.Called(msg.Topic, msg.Value)
	err := args.Error(0)
	if err == nil {
		m.cancelFn()
	}

	return err
}
