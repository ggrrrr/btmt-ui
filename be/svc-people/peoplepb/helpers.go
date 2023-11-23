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
		Id:         p.Id,
		IdNumbers:  p.IdNumbers,
		LoginEmail: p.LoginEmail,
		Emails:     p.Emails,
		Name:       p.Name,
		FullName:   p.FullName,
		Gender:     p.Gender,
		Phones:     p.Phones,
		Labels:     p.Labels,
		Attr:       p.Attr,
		CreatedAt:  timestamppb.New(p.CreatedTime),
	}
	if p.Age != nil {
		out.Age = fmt.Sprintf("%d", *p.Age)
	}
	if p.DOB != nil {
		out.Dob = &Dob{}
		out.Dob.Year = uint32(p.DOB.Year)
		out.Dob.Month = uint32(p.DOB.Month)
		out.Dob.Day = uint32(p.DOB.Day)
	}
	return &out
}

func (p *Person) ToPerson() *ddd.Person {
	out := ddd.Person{
		Id:          p.Id,
		IdNumbers:   p.IdNumbers,
		LoginEmail:  p.LoginEmail,
		Emails:      p.Emails,
		Name:        p.Name,
		FullName:    p.FullName,
		Gender:      p.Gender,
		Phones:      p.Phones,
		Labels:      p.Labels,
		Attr:        p.Attr,
		CreatedTime: p.CreatedAt.AsTime(),
	}
	if p.Dob != nil {
		out.DOB = &ddd.Dob{}
		out.DOB.Year = int(p.Dob.Year)
		out.DOB.Month = int(p.Dob.Month)
		out.DOB.Day = int(p.Dob.Day)
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
