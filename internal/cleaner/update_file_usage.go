package cleaner

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

func (g *GarbageCollector) UpdateFileUsage(fileId string, inUse bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	uses, err := g.storage.UpdateStats(ctx, fileId, inUse)

	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("Failed to update file '%s' usage statistics.", fileId))
		return
	}

	if uses == 0 {
		err = g.storage.MarkAsUnused(ctx, fileId)

		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("Failed to mark file '%s' as unused.", fileId))
			return
		}
	}

	if uses == 1 {
		err = g.storage.RemoveFromUnused(ctx, fileId)

		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("Failed to remove file '%s' from unused.", fileId))
			return
		}
	}

	log.Info().Msg(fmt.Sprintf("File '%s' usage statistics updated (uses: %d).", fileId, uses))
}
