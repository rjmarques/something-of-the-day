package datastore
import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/rjmarques/something-of-the-day/model"
)

const dbQueryTimeout = 30 * time.Second

type DAO struct {
	db *sql.DB
}

func NewDAO(db *sql.DB) *DAO {
	return &DAO{
		db: db,
	}
}

func (d *DAO) StoreSomethings(sms []*model.Something) error {
	log.Printf("saving %d something(s) to storage", len(sms))

	ctx, cancel := context.WithTimeout(context.Background(), dbQueryTimeout)
	defer cancel()

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error creating transaction: %v", err)
	}

	const sql = `INSERT INTO something.somethings(external_id, created_at, text) VALUES ($1, $2, $3) RETURNING id;`
	somethingStmt, err := tx.Prepare(sql)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer somethingStmt.Close()

	for _, s := range sms {
		var insertedId int64
		err := somethingStmt.QueryRow(s.Id, s.CreatedAt, s.Text).Scan(&insertedId)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert something %d: %v", s.Id, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error commiting transaction: %v", err)
	}

	return nil
}

func (d *DAO) GetSomethings() ([]*model.Something, error) {
	log.Println("getting somethings from storage")

	ctx, cancel := context.WithTimeout(context.Background(), dbQueryTimeout)
	defer cancel()

	const sql = `SELECT * FROM something.somethings ORDER BY created_at;`
	rows, err := d.db.QueryContext(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("error getting somethings: %v", err)
	}
	defer rows.Close()

	somethings := []*model.Something{}
	for rows.Next() {
		s := &model.Something{}

		var rowId int64
		err := rows.Scan(
			&rowId,
			&s.Id,
			&s.CreatedAt,
			&s.Text,
		)
		if err != nil {
			return nil, err
		}

		somethings = append(somethings, s)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error getting somethings: %v", err)
	}

	return somethings, nil
}

