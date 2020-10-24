package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/pkg/config"
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

func NewPostgresRepository() updatetimerevent.Repository {
	return &PostgresRepository{}
}

// Find timer event by user id.
func (r *PostgresRepository) FindTimerEvent(ctx context.Context, userId string) (event *enterpriserule.TimerEvent, err error) {
	db, err := getPostgresDb(ctx, config.Get("DATABASE_URL", ""))
	if err != nil {
		return
	}

	defer db.Close()

	event = &enterpriserule.TimerEvent{}
	if err = db.GetContext(ctx, event, fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", config.Get("POSTGRES_TBL_TIMEREVENT", "")), userId); err != nil {
		return nil, nil
	}

	event.UtcTimeToDrink = event.UtcTimeToDrink.UTC()

	return
}

// Find timer event from "from" to "to".
func (r *PostgresRepository) FindTimerEventByTime(ctx context.Context, from, to time.Time) (results []*enterpriserule.TimerEvent, err error) {
	db, err := getPostgresDb(ctx, config.Get("DATABASE_URL", ""))
	if err != nil {
		return
	}

	defer db.Close()

	values := []enterpriserule.TimerEvent{}
	if err = db.SelectContext(ctx, &values, fmt.Sprintf("SELECT * FROM %s WHERE $1 <= utc_time_to_drink AND utc_time_to_drink <= $2", config.Get("POSTGRES_TBL_TIMEREVENT", "")), from, to); err != nil {
		return nil, nil
	}

	results = make([]*enterpriserule.TimerEvent, len(values))
	for i := range values {
		values[i].UtcTimeToDrink = values[i].UtcTimeToDrink.UTC()
		results[i] = &values[i]
	}
	return
}

// Save TimerEvent to DB.
//
// Return error and the slice of TimerEvent saved successfully.
func (r *PostgresRepository) SaveTimerEvent(ctx context.Context, events []*enterpriserule.TimerEvent) (saved []*enterpriserule.TimerEvent, err error) {
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

	table := config.Get("POSTGRES_TBL_TIMEREVENT", "")
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
