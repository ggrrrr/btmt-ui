package ddd

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	peoplepbv1 "github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb/v1"
	// peoplepb "github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb/v1"
)

type (
	DOB struct {
		Year  int
		Month int
		Day   int
	}
	PinValidation struct {
		DOB    DOB
		Gender string
	}

	PeopleRepo interface {
		Save(ctx context.Context, p *peoplepbv1.Person) error
		Update(ctx context.Context, p *peoplepbv1.Person) error
		List(ctx context.Context, filter app.FilterFactory) ([]*peoplepbv1.Person, error)
		GetById(ctx context.Context, id string) (*peoplepbv1.Person, error)
	}
)
