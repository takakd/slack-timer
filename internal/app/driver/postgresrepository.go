package driver

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"proteinreminder/internal/app/enterpriserule"
	"proteinreminder/internal/app/usecase/updateproteinevent"
	"proteinreminder/internal/pkg/config"
	"time"
)

// Implements Repository interface with PostgreSQL.
type PostgresRepository struct {
}

// Get sqlx.DB
func getPostgresDb(ctx context.Context, dsn string) (db *sqlx.DB, err error) {
	db, err = sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		db = nil
		return
	}
	return
}

func NewPostgresRepository() updateproteinevent.Repository {
	return &PostgresRepository{}
}

// Find protein event by user id.
func (r *PostgresRepository) FindProteinEvent(ctx context.Context, userId string) (event *enterpriserule.ProteinEvent, err error) {
	db, err := getPostgresDb(ctx, config.Get("DATABASE_URL", ""))
	if err != nil {
		return
	}

	defer db.Close()

	event = &enterpriserule.ProteinEvent{}
	if err = db.GetContext(ctx, event, fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", config.Get("POSTGRES_TBL_PROTEINEVENT", "")), userId); err != nil {
		return nil, nil
	}

	event.UtcTimeToDrink = event.UtcTimeToDrink.UTC()

	return
}

// Find protein event from "from" to "to".
func (r *PostgresRepository) FindProteinEventByTime(ctx context.Context, from, to time.Time) (results []*enterpriserule.ProteinEvent, err error) {
	db, err := getPostgresDb(ctx, config.Get("DATABASE_URL", ""))
	if err != nil {
		return
	}

	defer db.Close()

	values := []enterpriserule.ProteinEvent{}
	if err = db.SelectContext(ctx, &values, fmt.Sprintf("SELECT * FROM %s WHERE $1 <= utc_time_to_drink AND utc_time_to_drink <= $2", config.Get("POSTGRES_TBL_PROTEINEVENT", "")), from, to); err != nil {
		return nil, nil
	}

	results = make([]*enterpriserule.ProteinEvent, len(values))
	for i := range values {
		values[i].UtcTimeToDrink = values[i].UtcTimeToDrink.UTC()
		results[i] = &values[i]
	}
	return
}

// Save ProteinEvent to DB.
//
// Return error and the slice of ProteinEvent saved successfully.
func (r *PostgresRepository) SaveProteinEvent(ctx context.Context, events []*enterpriserule.ProteinEvent) (saved []*enterpriserule.ProteinEvent, err error) {
	db, err := getPostgresDb(ctx, config.Get("DATABASE_URL", ""))
	if err != nil {
		return
	}

	defer db.Close()

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	table := config.Get("POSTGRES_TBL_PROTEINEVENT", "")
	for _, event := range events {
		_, err = tx.NamedExecContext(ctx, fmt.Sprintf(`
			INSERT INTO %s (user_id, utc_time_to_drink, drink_time_interval_min)
			VALUES (:user_id, :utc_time_to_drink, :drink_time_interval_min)
			ON CONFLICT (user_id) DO UPDATE
			SET utc_time_to_drink = :utc_time_to_drink,
				drink_time_interval_min = :drink_time_interval_min
		`, table), event)
		if err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()

	saved = events
	return
}
