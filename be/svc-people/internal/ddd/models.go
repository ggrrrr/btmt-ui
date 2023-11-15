package ddd

import (
	"context"
	"time"
)

type (
	Dob struct {
		Year  int
		Month int
		Day   int
	}
	Person struct {
		Id          string
		PIN         string
		Email       string
		Name        string
		FullName    string
		DateOfBirth *Dob
		Gender      string
		Phones      map[string]string
		Labels      []string
		Attr        map[string]string
		Age         *int
		CreatedTime time.Time
	}

	PeopleRepo interface {
		Save(ctx context.Context, p *Person) error
		Update(ctx context.Context, p *Person) error
		List(ctx context.Context, filter FilterFactory) ([]Person, error)
		GetById(ctx context.Context, id string) (*Person, error)
	}

	FilterFactory interface {
		Create() any
	}
)
