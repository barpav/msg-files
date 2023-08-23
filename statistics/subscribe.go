package statistics

import "fmt"

func (c *Client) Subscribe(handler func(fileId string, inUse bool)) error {
	stats, err := c.channel.Consume(c.queue.Name, "", true, false, false, false, nil)

	if err != nil {
		return fmt.Errorf("failed to subscribe on files usage statistics: %w", err)
	}

	go func() {
		info := FileUsage{}

		for s := range stats {
			if info.deserialize(s.Body) != nil {
				continue
			}

			handler(info.FileId, info.InUse)
		}
	}()

	return nil
}
