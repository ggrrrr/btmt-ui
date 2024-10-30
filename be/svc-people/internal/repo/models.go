package repo

import (
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/ddd"
)

const (
	FieldIdNumbers  string = "id_numbers"
	FieldLoginEmail string = "login_email"
	FieldEmails     string = "emails"
	FieldName       string = "name"
	FieldFullName   string = "full_name"
	FieldDoB        string = "dob"
	FieldGender     string = "gender"
	FieldPhones     string = "phones"
	FieldLabels     string = "labels"
	FieldAttr       string = "attr"
)

type (
	dob struct {
		Year  int `bson:"year"`
		Month int `bson:"month"`
		Day   int `bson:"day"`
	}
	person struct {
		Id          primitive.ObjectID `bson:"_id"`
		IdNumbers   []string           `bson:"id_numbers"`
		PIN         string             `bson:"pin"`
		LoginEmail  string             `bson:"login_email"`
		Emails      []string           `bson:"emails"`
		Name        string             `bson:"name"`
		FullName    string             `bson:"full_name"`
		DOB         *dob               `bson:"dob"`
		Gender      string             `bson:"gender"`
		Phones      []string           `bson:"phones"`
		Labels      []string           `bson:"labels"`
		Attr        []string           `bson:"attr"`
		CreatedTime primitive.DateTime `bson:"created_ts"`
	}
)

func toSlice(in map[string]string) []string {
	out := []string{}
	for kk, vv := range in {
		k := strings.Trim(kk, " ")
		v := strings.Trim(vv, " ")

		out = append(out, fmt.Sprintf("%s:%s", k, v))
	}
	return out
}

func toMap(in []string) map[string]string {
	out := map[string]string{}
	for _, vv := range in {
		kv := strings.Split(vv, ":")
		k := strings.Trim(kv[0], " ")
		v := strings.Trim(kv[0], " ")
		if len(kv) > 1 {
			k = strings.Trim(kv[0], " ")
			v = strings.Trim(kv[1], " ")
		}
		out[k] = v
	}
	return out
}

func fromDob(fromDob *ddd.Dob) *dob {
	if fromDob == nil {
		return nil
	}
	return &dob{
		Year:  fromDob.Year,
		Month: fromDob.Month,
		Day:   fromDob.Day,
	}
}

func toDob(fromDob *dob) *ddd.Dob {
	if fromDob == nil {
		return nil
	}
	return &ddd.Dob{
		Year:  fromDob.Year,
		Month: fromDob.Month,
		Day:   fromDob.Day,
	}
}

func fromPerson(p *ddd.Person) (*person, error) {
	id, err := mgo.ConvertFromId(p.Id)
	if err != nil {
		return nil, app.BadRequestError("invalid person.id", err)
	}

	out := person{
		Id:          id,
		IdNumbers:   toSlice(p.IdNumbers),
		LoginEmail:  p.LoginEmail,
		Emails:      toSlice(p.Emails),
		Name:        p.Name,
		FullName:    p.FullName,
		DOB:         fromDob(p.DOB),
		Gender:      p.Gender,
		Phones:      toSlice(p.Phones),
		Labels:      p.Labels,
		Attr:        toSlice(p.Attr),
		CreatedTime: mgo.FromTimeOrNow(p.CreatedTime),
	}

	return &out, nil
}

func (p *person) toPerson() *ddd.Person {
	var ts time.Time
	if p.CreatedTime > 0 {
		ts = p.CreatedTime.Time()
	}

	out := ddd.Person{
		Id:          p.Id.Hex(),
		IdNumbers:   toMap(p.IdNumbers),
		LoginEmail:  p.LoginEmail,
		Emails:      toMap(p.Emails),
		Name:        p.Name,
		FullName:    p.FullName,
		Gender:      p.Gender,
		DOB:         toDob(p.DOB),
		Phones:      toMap(p.Phones),
		Labels:      p.Labels,
		Attr:        toMap(p.Attr),
		CreatedTime: ts,
	}
	if p.DOB != nil {
		if p.DOB.Year > 0 {
			age := time.Now().Year() - p.DOB.Year
			out.Age = &age
		}
	}
	return &out
}

// if at least one of the DOB fields set then return true
func (d dob) isZero() bool {
	if d.Year > 0 {
		return false
	}
	if d.Month > 0 {
		return false
	}
	if d.Day > 0 {
		return false
	}
	return true
}
