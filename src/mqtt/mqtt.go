package mqtt

import (
	"context"

	"go.uber.org/zap"

	"github.com/pkg/errors"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type ListenerConfig interface {
	GetTopic() string
	GetBrokerHost() string
}

type Listener struct {
	topic  string
	client MQTT.Client
	qos    byte
	logger *zap.SugaredLogger
	// message contains topic at [0] and message at [1]
	message chan [2]string
}

func NewListener(config ListenerConfig, logger *zap.SugaredLogger) *Listener {
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
		qos:     0,
		logger:  logger,
		message: message}
}

// Listen MQTT messages until ctx cancelled. This method is blocking.
func (l *Listener) Listen(ctx context.Context) error {

	if token := l.client.Connect(); token.Wait() && token.Error() != nil {
		return errors.Wrap(token.Error(), "unable to connect to MQTT")
	}

	if token := l.client.Subscribe(l.topic, l.qos, nil); token.Wait() && token.Error() != nil {
		return errors.Wrap(token.Error(), "unable to subscribe to MQTT")
	}

	for {
		select {
		case incoming := <-l.message:
			l.logger.Info("received message", zap.String("topic", incoming[0]), zap.String("message", incoming[1]))
		case <-ctx.Done():
			close(l.message)
			return nil
		}
	}
}
