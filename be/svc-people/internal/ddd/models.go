package ddd

import (
	"context"
	"time"
)

type (
	PinValidation struct {
		Dob    Dob
		Gender string
	}

	Dob struct {
		Year  int
		Month int
		Day   int
	}
	Person struct {
		Id          string
		IdNumbers   map[string]string
		LoginEmail  string
		Name        string
		FullName    string
		DOB         *Dob
		Gender      string
		Emails      map[string]string
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
