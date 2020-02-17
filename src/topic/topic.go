package topic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type ListenerConfig struct {
	Topic string
	BrokerHost string
}

func Listen(ctx context.Context, config ListenerConfig) error {
	qos := 0
	opts := MQTT.NewClientOptions()
	opts.AddBroker(config.BrokerHost)

	receiveCount := 0
	choke := make(chan [2]string)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		choke <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return errors.Wrap(token.Error(), "unable to connect to MQTT")
	}

	if token := client.Subscribe(config.Topic, byte(qos), nil); token.Wait() && token.Error() != nil {
		return errors.Wrap(token.Error(), "unable to subscribe to MQTT")
	}

	go func() {
		for {
			select {
			case incoming := <-choke:
				fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
				receiveCount++
				fmt.Printf("%v\n", receiveCount)
/*			case _ <- ctx.Done():
				close(choke)
				return
*/			}
		}
	}()

	return nil
}