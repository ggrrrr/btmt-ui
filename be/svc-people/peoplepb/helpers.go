package peoplepb

import (
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/ddd"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (f *ListRequest) ToFilter() map[string][]string {
	out := map[string][]string{}
	if f.Filters == nil {
		return out
	}
	for k, v := range f.Filters {
		if v != nil {
			out[k] = v.GetList()
		}
	}
	return out
}

func FromPerson(p *ddd.Person) *Person {
	out := Person{
		Id:       p.Id,
		Pin:      p.PIN,
		Email:    p.Email,
		Name:     p.Name,
		FullName: p.FullName,
		// Dob:      p.DateOfBirth,
		Gender:    p.Gender,
		Phones:    p.Phones,
		Labels:    p.Labels,
		Attr:      p.Attr,
		CreatedAt: timestamppb.New(p.CreatedTime),
		Age:       fmt.Sprintf("%d", p.Age),
	}
	if p.DateOfBirth != nil {
		out.Dob.Year = uint32(p.DateOfBirth.Year)
		out.Dob.Month = uint32(p.DateOfBirth.Month)
		out.Dob.Day = uint32(p.DateOfBirth.Day)
	}
	return &out
}

func (p *Person) ToPerson() *ddd.Person {
	out := ddd.Person{
		Id:          p.Id,
		PIN:         p.Pin,
		Email:       p.Email,
		Name:        p.Name,
		FullName:    p.FullName,
		Gender:      p.Gender,
		Phones:      p.Phones,
		Labels:      p.Labels,
		Attr:        p.Attr,
		CreatedTime: p.CreatedAt.AsTime(),
	}
	if p.Dob != nil {
		out.DateOfBirth.Year = int(p.Dob.Year)
		out.DateOfBirth.Month = int(p.Dob.Month)
		out.DateOfBirth.Day = int(p.Dob.Day)
	}
	return &out
}

func (p *SaveRequest) ToPerson() *ddd.Person {
	if p.Data == nil {
		return &ddd.Person{}
	}
	return p.Data.ToPerson()
}

func (p *UpdateRequest) ToPerson() *ddd.Person {
	if p.Data == nil {
		return &ddd.Person{}
	}
	return p.Data.ToPerson()
}
