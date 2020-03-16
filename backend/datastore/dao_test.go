package datastore
import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"

	"github.com/rjmarques/something-of-the-day/model"
)

func TestStoreSomethings(t *testing.T) {
	db := createDBConn(t)
	dao := NewDAO(db)

	somethings := []*model.Something{
		&model.Something{
			Id:        123,
			CreatedAt: time.Now().Add(-24 * time.Hour),
			Text:      "What's the best thing about Switzerland? I don't know, but the flag is a big plus.",
		},
		&model.Something{
			Id:        456,
			CreatedAt: time.Now(),
			Text:      "Most people think T-Rexes can't clap their hands because they have short arms whereas actually it's because they're dead.",
		},
	}

	err := dao.StoreSomethings(somethings)
	assert.NoError(t, err)

	for _, s := range somethings {
		saved := &model.Something{}
		err := db.QueryRow("SELECT external_id, created_at, text FROM something.somethings WHERE external_id = $1", s.Id).Scan(
			&saved.Id,
			&saved.CreatedAt,
			&saved.Text,
		)
		assert.NoError(t, err)
		assert.Equal(t, s.Id, saved.Id)
		assert.Equal(t, s.CreatedAt.Truncate(time.Second).UTC(), saved.CreatedAt.Truncate(time.Second).UTC())
		assert.Equal(t, s.Text, saved.Text)

		// cleanup
		_, err = db.Exec("DELETE FROM something.somethings WHERE external_id = $1;", saved.Id)
		assert.NoError(t, err)
	}
}

func TestGetSomethings(t *testing.T) {
	db := createDBConn(t)
	dao := NewDAO(db)

	// check if db is empty
	somethings, err := dao.GetSomethings()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(somethings))

	// add a couple somethings
	somethings = []*model.Something{
		&model.Something{
			Id:        123,
			CreatedAt: time.Now().Add(-24 * time.Hour),
			Text:      "What's the best thing about Switzerland? I don't know, but the flag is a big plus.",
		},
		&model.Something{
			Id:        456,
			CreatedAt: time.Now(),
			Text:      "Most people think T-Rexes can't clap their hands because they have short arms whereas actually it's because they're dead.",
		},
	}
	err = dao.StoreSomethings(somethings)
	assert.NoError(t, err)

	// check to see if they're read correctly
	read, err := dao.GetSomethings()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(read))
	for i, s := range somethings {
		assert.Equal(t, s.Id, read[i].Id)
		assert.Equal(t, s.CreatedAt.Truncate(time.Second).UTC(), read[i].CreatedAt.Truncate(time.Second).UTC())
		assert.Equal(t, s.Text, read[i].Text)
	}

	// cleanup
	_, err = db.Exec("TRUNCATE TABLE something.somethings;")
	assert.NoError(t, err)
}

// TestDBConn is a helper for integration tests
// build a connection from the POSTGRES_URL environment variable
// or skip the test if not available
func createDBConn(t *testing.T) *sql.DB {
	t.Helper()

	url := os.Getenv("POSTGRES_URL")
	if url == "" {
		t.Skip("Skipping integration test")
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		t.Fatalf("failed to open a connection to the database %v", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping database %v", err)
	}

	return db
}

