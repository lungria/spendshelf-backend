package mqtt

import (
	"context"
	"encoding/json"

	"github.com/lungria/spendshelf-backend/src/transactions"

	"go.uber.org/zap"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type ListenerConfig interface {
	GetTopic() string
	GetBrokerHost() string
}

const qos = 1

type Listener struct {
	topic  string
	client MQTT.Client
	logger *zap.SugaredLogger
	// message contains topic at [0] and message at [1]
	message chan [2]string
	repo    *transactions.Repository
}

func NewListener(config ListenerConfig, logger *zap.SugaredLogger, repo *transactions.Repository) *Listener {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(config.GetBrokerHost())

	message := make(chan [2]string)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		message <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := MQTT.NewClient(opts)
	return &Listener{
		client:  client,
		topic:   config.GetTopic(),
		logger:  logger,
		message: message,
		repo:    repo,
	}
}

// Listen MQTT messages until ctx cancelled. This method is blocking.
func (l *Listener) Listen(ctx context.Context) error {

	if token := l.client.Connect(); token.Wait() && token.Error() != nil {
		l.logger.Fatal("mqtt connection failed", zap.Error(token.Error()))
	}

	if token := l.client.Subscribe(l.topic, qos, nil); token.Wait() && token.Error() != nil {
		l.logger.Fatal("mqtt subscription failed", zap.Error(token.Error()))
	}

	for {
		select {
		case incoming := <-l.message:
			l.logger.Info("received message", zap.String("topic", incoming[0]))
			var t []transactions.Transaction
			err := json.Unmarshal([]byte(incoming[1]), &t)
			if err != nil {
				l.logger.Error("unable to unmarshal json: ", zap.Error(err))
				continue
			}
			err = l.repo.InsertMany(ctx, t)
			if err != nil {
				l.logger.Error("unable to save transaction: ", zap.Error(err))
			}
		case <-ctx.Done():
			close(l.message)
			return nil
		}
	}
}
