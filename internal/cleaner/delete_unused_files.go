package cleaner

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

const filesPerIteration int64 = 100

func (g *GarbageCollector) deleteUnusedFiles() {
	log.Info().Msg("Garbage collection session started.")

	var (
		ids            []string
		id             string
		err            error
		requeue        bool
		interrupt      bool
		filesProcessed int
	)

	ctx := context.Background()
	bound := time.Now().UTC().Unix() - int64(g.cfg.deleteAfter)

	for {
		ids, err = g.storage.UnusedFiles(ctx, bound, filesPerIteration)

		if err != nil {
			log.Err(err).Msg("Failed to receive unused files.")
			break
		}

		if len(ids) == 0 {
			break
		}

		for _, id = range ids {
			filesProcessed++

			err = g.storage.DeleteFile(ctx, id)
			requeue = err != nil

			if err != nil {
				log.Err(err).Msg(fmt.Sprintf("Failed to delete unused file '%s'.", id))
			}

			err = g.storage.RemoveFromUnused(ctx, id)

			if err != nil {
				log.Err(err).Msg(fmt.Sprintf("Failed to remove file '%s' from unused.", id))
				interrupt = true // possible dead loop - try later in new session
				continue
			}

			if requeue {
				err = g.storage.MarkAsUnused(ctx, id) // requeue with new timestamp (try to delete later)

				if err != nil {
					log.Err(err).Msg(fmt.Sprintf("Failed to requeue unused file '%s'.", id))
					continue
				}

				log.Info().Msg(fmt.Sprintf("Unused file '%s' requeued.", id))
			}
		}

		if interrupt {
			log.Info().Msg("Garbage collection session interrupted due to occurred errors.")
			break
		}
	}

	log.Info().Msg(fmt.Sprintf("Garbage collection session finished. Unused files processed: %d.", filesProcessed))
}
