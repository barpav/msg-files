package statistics

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *Client) SendUsage(ctx context.Context, fileId string, inUse bool) error {
	usageStats := FileUsage{FileId: fileId, InUse: inUse}
	data, err := usageStats.serialize()

	if err == nil {
		err = c.channel.PublishWithContext(ctx, "", c.queue.Name, false, false, amqp.Publishing{Body: data})
	}

	if err != nil {
		return fmt.Errorf("failed to send file usage info: %w", err)
	}

	return nil
}
