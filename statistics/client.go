package statistics

import (
	"context"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	cfg     *config
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

const queueName = "msg.files.statistics"

func (c *Client) Connect() (err error) {
	c.cfg = &config{}
	c.cfg.read()

	c.conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", c.cfg.user, c.cfg.password, c.cfg.host, c.cfg.port))

	if err == nil {
		c.channel, err = c.conn.Channel()
	}

	if err == nil {
		c.queue, err = c.channel.QueueDeclare(queueName, true, false, false, false, nil)
	}

	if err != nil {
		return fmt.Errorf("can't connect to files usage statistics: %w", err)
	}

	return nil
}

func (c *Client) Disconnect(ctx context.Context) (err error) {
	var closeErr error
	closed := make(chan struct{}, 1)

	go func() {
		if c.channel != nil {
			closeErr = errors.Join(closeErr, c.channel.Close())
		}

		if c.conn != nil {
			closeErr = errors.Join(closeErr, c.conn.Close())
		}

		closed <- struct{}{}
	}()

	select {
	case <-closed:
		err = closeErr
	case <-ctx.Done():
		err = ctx.Err()
	}

	if err != nil {
		err = fmt.Errorf("failed to disconnect from file usage statistics: %w", err)
	}

	return err
}
