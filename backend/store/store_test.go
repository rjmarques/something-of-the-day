package store
import (
	"time"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rjmarques/something-of-the-day/model"
)

func TestAdd(t *testing.T) {
	s := NewStore()

	somethings := []*model.Something{
		&model.Something{
			Id: 123,
			CreatedAt: time.Now(),
			Text: "Lorem Ipsum",
		},
	}
	s.Add(somethings)

	saved := s.GetRand() // nothing else is in the store, so can only return the previous struct
	assert.NotNil(t, saved)
	assert.Equal(t, somethings[0].Id, saved.Id)
	assert.Equal(t, somethings[0].CreatedAt, saved.CreatedAt)
	assert.Equal(t, somethings[0].Text, saved.Text)
}


func TestGetLatest(t *testing.T) {
	s := NewStore()

	somethings := []*model.Something{
		&model.Something{
			Id: 123,
			CreatedAt: time.Now(),
			Text: "Lorem Ipsum",
		},
		&model.Something{
			Id: 456,
			CreatedAt: time.Now(),
			Text: "Carpe Diem",
		},
		&model.Something{
			Id: 789,
			CreatedAt: time.Now(),
			Text: "Per aspera ad astra",
		},
	}
	s.Add(somethings)

	latest := s.GetLatest() // nothing else is in the store, so can only return the previous struct
	assert.NotNil(t, latest)
	assert.Equal(t, somethings[2].Id, latest.Id)
	assert.Equal(t, somethings[2].CreatedAt, latest.CreatedAt)
	assert.Equal(t, somethings[2].Text, latest.Text)
}
