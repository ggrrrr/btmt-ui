package repo

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	peoplepbv1 "github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb/v1"
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
	FieldCreatedAt  string = "created_at"
	FieldUpdatedAt  string = "updated_at"
)

type (
	dob struct {
		Year  uint32 `bson:"year"`
		Month uint32 `bson:"month"`
		Day   uint32 `bson:"day"`
	}
	person struct {
		Id         primitive.ObjectID `bson:"_id"`
		IdNumbers  []string           `bson:"id_numbers"`
		PIN        string             `bson:"pin"`
		LoginEmail string             `bson:"login_email"`
		Emails     []string           `bson:"emails"`
		Name       string             `bson:"name"`
		FullName   string             `bson:"full_name"`
		DOB        *dob               `bson:"dob"`
		Gender     string             `bson:"gender"`
		Phones     []string           `bson:"phones"`
		Labels     []string           `bson:"labels"`
		Attr       []string           `bson:"attr"`
		CreatedAt  primitive.DateTime `bson:"created_at"`
		UpdatedAt  primitive.DateTime `bson:"updated_at"`
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

func fromDob(fromDob *peoplepbv1.Dob) *dob {
	if fromDob == nil {
		return nil
	}
	return &dob{
		Year:  fromDob.Year,
		Month: fromDob.Month,
		Day:   fromDob.Day,
	}
}

func toDob(fromDob *dob) *peoplepbv1.Dob {
	if fromDob == nil {
		return nil
	}
	return &peoplepbv1.Dob{
		Year:  fromDob.Year,
		Month: fromDob.Month,
		Day:   fromDob.Day,
	}
}

func fromPerson(p *peoplepbv1.Person) (*person, error) {
	id, err := mgo.ConvertFromId(p.Id)
	if err != nil {
		return nil, app.BadRequestError("invalid person.id", err)
	}

	out := person{
		Id:         id,
		IdNumbers:  toSlice(p.IdNumbers),
		LoginEmail: p.LoginEmail,
		Emails:     toSlice(p.Emails),
		Name:       p.Name,
		FullName:   p.FullName,
		DOB:        fromDob(p.Dob),
		Gender:     p.Gender,
		Phones:     toSlice(p.Phones),
		Labels:     p.Labels,
		Attr:       toSlice(p.Attr),
	}

	if !p.CreatedAt.AsTime().IsZero() {
		out.CreatedAt = mgo.FromTimeOrNow(p.CreatedAt.AsTime())
	}

	if !p.UpdatedAt.AsTime().IsZero() {
		out.UpdatedAt = mgo.FromTimeOrNow(p.UpdatedAt.AsTime())
	}

	return &out, nil
}

func (p *person) toProto() *peoplepbv1.Person {
	out := &peoplepbv1.Person{
		Id:         p.Id.Hex(),
		IdNumbers:  toMap(p.IdNumbers),
		LoginEmail: p.LoginEmail,
		Emails:     toMap(p.Emails),
		Name:       p.Name,
		FullName:   p.FullName,
		Gender:     p.Gender,
		Dob:        toDob(p.DOB),
		Phones:     toMap(p.Phones),
		Labels:     p.Labels,
		Attr:       toMap(p.Attr),
	}

	if !p.CreatedAt.Time().IsZero() {
		out.CreatedAt = timestamppb.New(p.CreatedAt.Time())
	}
	if !p.UpdatedAt.Time().IsZero() {
		out.UpdatedAt = timestamppb.New(p.UpdatedAt.Time())
	}

	return out
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
