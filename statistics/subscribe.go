package statistics

import "fmt"

func (c *Client) Subscribe(handler func(fileId string, inUse bool) (ok bool)) error {
	stats, err := c.channel.Consume(c.queue.Name, "", false, false, false, false, nil)

	if err != nil {
		return fmt.Errorf("failed to subscribe on files usage statistics: %w", err)
	}

	go func() {
		var sErr error
		info := FileUsage{}

		for s := range stats {
			sErr = info.deserialize(s.Body)

			if sErr != nil {
				s.Reject(false)
				continue
			}

			if handler(info.FileId, info.InUse) {
				s.Ack(false)
				continue
			}

			s.Reject(true)
		}
	}()

	return nil
}
