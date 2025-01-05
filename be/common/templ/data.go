package templ

import templv1 "github.com/ggrrrr/btmt-ui/be/common/templ/v1"

type (
	templPerson struct {
		Id       string
		Name     string
		Email    string
		FullName string
	}

	templData struct {
		Person templPerson
		Items  map[string]any
	}
)

func fromV1(from *templv1.Data) templData {
	out := templData{
		Items: map[string]any{},
	}
	if from == nil {
		return out
	}
	if from.Person != nil {
		out.Person.Id = from.Person.Id
		out.Person.Name = from.Person.Name
		out.Person.Email = from.Person.Email
		out.Person.FullName = from.Person.FullName
	}
	if from.Items != nil {
		for k := range from.Items {
			d := from.Items[k].AsMap()
			out.Items[k] = d
		}
	}
	return out

}
