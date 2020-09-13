package apprule

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"proteinreminder/internal/app/enterpriserule"
	"proteinreminder/internal/pkg/config"
	"time"
)

// Implements Repository interface with PostgreSQL.
type PostgresRepository struct {
	config config.Config
}

// Get sqlx.DB
func getPostgreDb(ctx context.Context, dataSourceName string) (db *sqlx.DB, err error) {
	db, err = sqlx.ConnectContext(ctx, "postgres", dataSourceName)
	if err != nil {
		return
	}
	return
}

func NewPostgresRepository(config config.Config) Repository {
	return &PostgresRepository{
		config,
	}
}

// Find protein event by user id.
func (r *PostgresRepository) FindProteinEvent(ctx context.Context, userId string) (event *enterpriserule.ProteinEvent, err error) {
	db, err := getPostgreDb(ctx, r.config.Get("POSTGRES_DATASOURCE"))
	if err != nil {
		return
	}

	event = &enterpriserule.ProteinEvent{}
	if err = db.GetContext(ctx, event, "SELECT * FROM protein_event WHERE user_id=$1", userId); err != nil {
		return nil, nil
	}
	return
}

// Find protein event from "from" to "to".
func (r *PostgresRepository) FindProteinEventByTime(ctx context.Context, from, to time.Time) (results []*enterpriserule.ProteinEvent, err error) {
	db, err := getPostgreDb(ctx, r.config.Get("POSTGRES_DATASOURCE"))
	if err != nil {
		return
	}

	values := []enterpriserule.ProteinEvent{}
	if err = db.GetContext(ctx, &values, "SELECT * FROM protein_event WHERE $1 <= utc_time_to_drink AND utc_time_to_drink < $2", from, to); err != nil {
		return nil, nil
	}
	for _, v := range values {
		results = append(results, &v)
	}
	return
}

// Save ProteinEvent to DB.
//
// Return error and the slice of ProteinEvent saved successfully.
func (r *PostgresRepository) SaveProteinEvent(ctx context.Context, events []*enterpriserule.ProteinEvent) (saved []*enterpriserule.ProteinEvent, err error) {
	db, err := getPostgreDb(ctx, r.config.Get("POSTGRES_DATASOURCE"))
	if err != nil {
		return
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, event := range events {
		_, err = tx.NamedExecContext(ctx, `
			INSERT INTO protein_event (user_id, utc_time_to_drink, drink_time_interval_sec)
			VALUES (:user_id, :utc_time_to_drink, :drink_time_interval_sec)
			ON CONFLICT (:user_id) DO UPDATE
			SET  utc_time_to_drink = :utc_time_to_drink, drink_time_interval_sec = :drink_time_interval_sec
		`, event)
		if err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()

	saved = events
	return
}
