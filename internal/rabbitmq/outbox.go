package rabbitmq

import (
	"database/sql"
	"log"
	"time"
)

const maxAttempts = 3
const interval = 5 * time.Second

func StartOutboxDispatcher(db *sql.DB, producer Producer) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			rows, err := db.Query("SELECT id, payload, dispatch_attempts FROM outbox WHERE dispatched_at IS NULL ORDER BY created_at LIMIT 50")
			if err != nil {
				log.Printf("Error querying outbox: %v", err)
				continue
			}

			var events []struct {
				ID               int
				Payload          []byte
				DispatchAttempts int
			}

			for rows.Next() {
				var e struct {
					ID               int
					Payload          []byte
					DispatchAttempts int
				}
				if err := rows.Scan(&e.ID, &e.Payload, &e.DispatchAttempts); err != nil {
					log.Printf("Error scanning outbox row: %v", err)
					continue
				}
				events = append(events, e)
			}
			rows.Close()

			for _, e := range events {
				err := producer.PublishRaw(e.Payload)
				if err != nil {
					log.Printf("Failed to publish event %d: %v", e.ID, err)

					_, updateErr := db.Exec("UPDATE outbox SET dispatch_attempts = dispatch_attempts + 1 WHERE id = $1", e.ID)
					if updateErr != nil {
						log.Printf("Failed to increment dispatch_attempts for event %d: %v", e.ID, updateErr)
						continue
					}

					newAttempts := e.DispatchAttempts + 1
					if newAttempts >= maxAttempts {
						_, deadErr := db.Exec(
							"UPDATE outbox SET event_type = 'dead_letter', dispatched_at = NOW() WHERE id = $1",
							e.ID,
						)
						if deadErr != nil {
							log.Printf("Failed to mark event %d as dead_letter: %v", e.ID, deadErr)
						} else {
							log.Printf("Event %d marked as dead_letter after %d attempts", e.ID, newAttempts)
						}
					}
					continue
				}

				_, err = db.Exec("UPDATE outbox SET dispatched_at = NOW() WHERE id = $1", e.ID)
				if err != nil {
					log.Printf("Failed to update outbox dispatched_at for event %d: %v", e.ID, err)
				}
			}
		}
	}()
}
