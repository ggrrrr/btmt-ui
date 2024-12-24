package repo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	peoplepbv1 "github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb/v1"
)

type (
	testCase struct {
		test string
		run  func(t *testing.T)
	}
)

func TestSave(t *testing.T) {

	expectedDuration := time.Duration(400 * time.Millisecond)

	ctx := context.Background()
	cfg := mgo.MgoTestCfg("test_people")
	testDb, err := mgo.New(ctx, cfg)
	require.NoError(t, err)
	defer testDb.Close(ctx)

	testRepo := New(cfg.Collection, testDb)

	err = testDb.DB().Collection(cfg.Collection).Drop(ctx)
	require.NoError(t, err)

	tests := []testCase{
		{
			test: "happy get not found",
			run: (func(t *testing.T) {
				noRec, err := testRepo.GetById(ctx, primitive.NewObjectID().Hex())
				assert.NoError(t, err)
				assert.Nil(t, noRec)
			}),
		},
		{
			test: "happy save get",
			run: func(t *testing.T) {
				expected := &peoplepbv1.Person{
					IdNumbers: map[string]string{"pin": "pin1"},
					Name:      "ggrrrr",
					Emails:    map[string]string{"": "asdasd@asd"},
					FullName:  "varban krushev",
					Labels:    []string{"tours:bike", "tours:hike", "kids"},
					Phones:    map[string]string{"mobile": "123123123"},
					Attr:      map[string]string{"food": "veg"},
					Gender:    "male",
				}
				ts := time.Now()

				err = testRepo.Save(ctx, expected)
				require.NoError(t, err)
				assert.True(t, expected.Id != "")
				assert.True(t, !expected.CreatedAt.AsTime().IsZero(), "Created Time must be set")
				assert.WithinDuration(t, ts, expected.CreatedAt.AsTime(), expectedDuration)

				actual, err := testRepo.GetById(ctx, expected.Id)
				require.NoError(t, err)
				require.NotNil(t, actual)
				TestPerson(t, actual, expected, expectedDuration)
			},
		},
		{
			test: "update all",
			run: func(t *testing.T) {
				p1 := &peoplepbv1.Person{}
				err = testRepo.Save(ctx, p1)
				assert.NoError(t, err)
				assert.True(t, p1.Id != "")
				assert.True(t, !p1.CreatedAt.AsTime().IsZero())

				expected := &peoplepbv1.Person{
					Id:           p1.Id,
					PrimaryEmail: "login@Email",
					Name:         "ggrrrr",
					Emails:       map[string]string{"": "asdasd@asd"},
					FullName:     "not varban krushev",
					IdNumbers:    map[string]string{"pin": "pin1"},
					Labels:       []string{"tours:bike", "tours:hike", "kids"},
					Phones:       map[string]string{"mobile": "123123123"},
					Attr:         map[string]string{"food": "veg"},
					Gender:       "male",
					Dob:          &peoplepbv1.Dob{Year: 1978, Month: 2, Day: 2},
					CreatedAt:    p1.CreatedAt,
				}

				err = testRepo.Update(ctx, expected)
				require.NoError(t, err)
				actual, err := testRepo.GetById(ctx, p1.Id)
				require.NoError(t, err)
				require.NotNil(t, actual)

				expected.UpdatedAt = timestamppb.Now()

				TestPerson(t, expected, actual, expectedDuration)

				expected = &peoplepbv1.Person{
					Id:           p1.Id,
					PrimaryEmail: " ",
					Name:         " ",
					Emails:       map[string]string{"": "asdasd@asd"},
					FullName:     " ",
					IdNumbers:    map[string]string{"pin ": "pin1 "},
					Labels:       []string{"tours:bike", "tours:hike", "kids"},
					Phones:       map[string]string{"mobile": "123123123"},
					Attr:         map[string]string{"food": "veg"},
					Gender:       " ",
					Dob:          &peoplepbv1.Dob{Year: 1978, Month: 2, Day: 2},
					CreatedAt:    p1.CreatedAt,
				}

				err = testRepo.Update(ctx, expected)
				require.NoError(t, err)
				actual, err = testRepo.GetById(ctx, p1.Id)
				require.NoError(t, err)

				expected.PrimaryEmail = ""
				expected.Name = ""
				expected.FullName = ""
				expected.Gender = ""
				expected.IdNumbers = map[string]string{"pin": "pin1"}
				expected.UpdatedAt = timestamppb.Now()

				TestPerson(t, expected, actual, expectedDuration)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			tc.run(t)
		})
	}

}

// https://github.com/mongodb/mongo-go-driver/blob/v1.12.1/examples/documentation_examples/examples.go
func TestList(t *testing.T) {
	ctx := context.Background()
	cfg := mgo.MgoTestCfg("test_people")
	testDb, err := mgo.New(ctx, cfg)
	require.NoError(t, err)
	defer testDb.Close(ctx)

	testRepo := New(cfg.Collection, testDb)

	err = testDb.DB().Collection(cfg.Collection).Drop(ctx)
	require.NoError(t, err)

	testRepo.CreateIndex(ctx)

	newData := map[string]*peoplepbv1.Person{
		"ggrrrr": {
			IdNumbers: map[string]string{"pin": "ggrrrrpin"},
			Name:      "ggrrrr",
			Emails:    map[string]string{"default": "ggrrrr@gmail.com"},
			// Emails:    [str]"ggrrrr@gmail.com",
			FullName: "ggrrrr varban krushev",
			// DateOfBirth: time.Date(1978, 2, 13, 0, 0, 0, 0, time.Local),
			Labels: []string{"tours:snow", "instructor:kids"},
			Phones: map[string]string{"mobile": "99009900"},
			Attr:   map[string]string{"food": "veg"},
			Gender: "male",
		},
		"mandajiev": {
			IdNumbers: map[string]string{"pin": "mandajievpin"},
			Name:      "mandajiev",
			Emails:    map[string]string{"default": "mandajiev@yahoo.com"},
			FullName:  "mandajiev asdasd asdasd",
			// DateOfBirth: time.Date(1990, 4, 23, 0, 0, 0, 0, time.Local),
			Labels: []string{"tours:bike", "volunteer:mtb", "bike:mtb"},
			Phones: map[string]string{"mobile": "223123123"},
			Attr:   map[string]string{"sleep": "no-tent"},
			Gender: "male",
		},
		"uniq": {
			IdNumbers: map[string]string{"pin": "NOPIN"},
			Name:      "uniq",
			Emails:    map[string]string{"default": "pesho@yahoo.com"},
			FullName:  "NONONO DDDD",
			// DateOfBirth: time.Date(1990, 4, 23, 0, 0, 0, 0, time.Local),
			Labels: []string{"shit"},
			Phones: map[string]string{"mobile": "somephone"},
			Attr:   map[string]string{"other": "noother"},
			Gender: "male",
		},
	}
	list, err := testRepo.List(ctx, nil)
	require.NoError(t, err)
	if len(list) == 0 {
		for _, p := range newData {
			err = testRepo.Save(ctx, p)
			require.NoError(t, err)
		}
		printMap("NEW DATA", newData)
	}

	tests := []testCase{
		{
			test: "happy list emails 1rc",
			run: (func(tt *testing.T) {
				filter, err := NewFilter(AddTexts("NONONO"))
				require.NoError(tt, err)
				list, err = testRepo.List(ctx, filter)
				require.NoError(tt, err)
				require.Equal(tt, 1, len(list), "records")
				for _, p := range list {
					fmt.Printf("%+v \n", p)
					newData[p.Name].Id = p.Id
					TestPerson(tt, p, newData[p.Name], 0)
				}
				printList("LIST", list)
			}),
		}, {
			test: "happy list all 3 records",
			run: (func(tt *testing.T) {
				list, err = testRepo.List(ctx, nil)
				require.NoError(tt, err)
				assert.Equal(tt, 3, len(list), "records")
				for _, p := range list {
					newData[p.Name].Id = p.Id
					TestPerson(tt, p, newData[p.Name], 0)
				}
				printList("LIST", list)
			}),
		}, {
			test: "happy list all two empty filter",
			run: (func(tt *testing.T) {
				filter, err := NewFilter()
				require.NoError(tt, err)

				list, err = testRepo.List(ctx, filter)
				require.NoError(tt, err)
				assert.Equal(tt, 3, len(list), "records")
				for _, p := range list {
					newData[p.Name].Id = p.Id
					TestPerson(tt, p, newData[p.Name], 0)
				}
				printList("LIST", list)
			}),
		},
		{
			test: "test regexpr for label add or  labels filter 1 rec",
			run: (func(tt *testing.T) {
				filter, err := NewFilter(AddLabels("instructor"), AddLabels("mtb"))
				require.NoError(tt, err)

				list, err = testRepo.List(ctx, filter)
				require.NoError(tt, err)
				assert.Equal(tt, 1, len(list), "records")
				for _, p := range list {
					newData[p.Name].Id = p.Id
					TestPerson(tt, p, newData[p.Name], 0)
				}
				printList("LIST", list)
			}),
		},
		{
			test: "test regexpr for label add or  labels filter 2 rec",
			run: (func(tt *testing.T) {
				// t.Skip("1")
				filter, err := NewFilter(AddLabels("instructor"), AddLabels("bike:mtb"))
				require.NoError(tt, err)

				list, err = testRepo.List(ctx, filter)
				require.NoError(tt, err)
				assert.Equal(tt, len(list), 2)
				for _, p := range list {
					newData[p.Name].Id = p.Id
					TestPerson(tt, p, newData[p.Name], 0)
				}
				printList("LIST", list)
			}),
		},
		{
			test: "test text search 0 rec",
			run: (func(tt *testing.T) {
				// t.Skip("1")
				filter, err := NewFilter(AddTexts("ggrrrrpin"), AddTexts("pin"))
				require.NoError(tt, err)

				list, err = testRepo.List(ctx, filter)
				require.NoError(tt, err)
				assert.Equal(tt, len(list), 0)
				for _, p := range list {
					newData[p.Name].Id = p.Id
					TestPerson(tt, p, newData[p.Name], 0)
				}
				printList("LIST", list)
			}),
		},
		{
			test: "test text search 2 rec",
			run: (func(tt *testing.T) {
				// t.Skip("1")
				filter, err := NewFilter(AddTexts("ggrrrr"), AddTexts("asdasd"))
				require.NoError(tt, err)

				list, err = testRepo.List(ctx, filter)
				require.NoError(tt, err)
				assert.Equal(tt, len(list), 2)
				for _, p := range list {
					newData[p.Name].Id = p.Id
					TestPerson(tt, p, newData[p.Name], 0)
				}
				printList("LIST", list)
			}),
		},
		{
			test: "test phone and label 1 rec",
			run: (func(tt *testing.T) {
				// t.Skip("1")
				filter, err := NewFilter(AddLabels("volunteer"), AddPhones("223123123"))
				require.NoError(tt, err)

				list, err = testRepo.List(ctx, filter)
				require.NoError(tt, err)
				assert.Equal(tt, len(list), 1)
				for _, p := range list {
					newData[p.Name].Id = p.Id
					TestPerson(tt, p, newData[p.Name], 0)
				}
				printList("LIST", list)
			}),
		},
		{
			test: "test phones and labels 2 rec",
			run: (func(tt *testing.T) {
				// t.Skip("1")
				filter, err := NewFilter(AddLabels("volunteer", "tours"), AddPhones("223123123", "99009900"))
				require.NoError(tt, err)

				list, err = testRepo.List(ctx, filter)
				require.NoError(tt, err)
				assert.Equal(tt, len(list), 2)
				for _, p := range list {
					newData[p.Name].Id = p.Id
					TestPerson(tt, p, newData[p.Name], 0)
				}
				printList("LIST", list)
			}),
		},
		{
			test: "test text and phone 1 rec",
			run: (func(tt *testing.T) {
				// t.Skip("1")
				filter, err := NewFilter(AddTexts("asdasd"), AddPhones("223123123", "99009900"))
				require.NoError(tt, err)

				list, err = testRepo.List(ctx, filter)
				require.NoError(tt, err)
				assert.Equal(tt, len(list), 1)
				for _, p := range list {
					newData[p.Name].Id = p.Id
					TestPerson(tt, p, newData[p.Name], 0)
				}
				printList("LIST", list)
			}),
		},
		{
			test: "pin 1 rec",
			run: (func(tt *testing.T) {
				// t.Skip("1")
				filter, err := NewFilter(AddPINs("mandajievpin"))
				require.NoError(tt, err)

				list, err = testRepo.List(ctx, filter)
				require.NoError(tt, err)
				assert.Equal(tt, len(list), 1)
				for _, p := range list {
					newData[p.Name].Id = p.Id
					TestPerson(tt, p, newData[p.Name], 0)
				}
				printList("LIST", list)
			}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.test, func(tt *testing.T) {
			tc.run(tt)
		})
	}

}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	cfg := mgo.MgoTestCfg("test-people")
	testDb, err := mgo.New(ctx, cfg)
	require.NoError(t, err)
	defer testDb.Close(ctx)

	testRepo := New(cfg.Collection, testDb)

	err = testDb.DB().Collection(cfg.Collection).Drop(ctx)
	require.NoError(t, err)

	tests := []testCase{
		{
			test: "happy get nil",
			run: (func(t *testing.T) {
				noRec, err := testRepo.GetById(ctx, primitive.NewObjectID().Hex())
				assert.NoError(t, err)
				assert.Nil(t, noRec)
			}),
		}, {
			test: "happy save get",
			run: func(t *testing.T) {
				p1 := &peoplepbv1.Person{
					IdNumbers: map[string]string{"pin": "pin1"},
					Name:      "ggrrrr",
					Emails:    map[string]string{"": "asdasd@asd"},
					FullName:  "varban krushev",
					Labels:    []string{"tours:bike", "tours:hike", "kids"},
					Phones:    map[string]string{"mobile": "123123123"},
					Attr:      map[string]string{"food": "veg"},
					Gender:    "male",
				}
				err = testRepo.Save(ctx, p1)
				require.NoError(t, err)
				require.True(t, p1.Id != "")
				require.True(t, !p1.CreatedAt.AsTime().IsZero(), "Created Time must be set")

				p2, err := testRepo.GetById(ctx, p1.Id)
				require.NoError(t, err)
				require.True(t, !p2.CreatedAt.AsTime().IsZero())
				p1.CreatedAt = p2.CreatedAt
				TestPerson(t, p2, p1, 10)
				fmt.Printf("got %v \n", p2)
			},
		},
		{
			test: "update",
			run: func(t *testing.T) {
				p1 := &peoplepbv1.Person{
					Name:     "ggrrrr",
					Emails:   map[string]string{"": "asdasd@asd"},
					FullName: "not varban krushev",
				}
				err = testRepo.Save(ctx, p1)
				assert.NoError(t, err)
				assert.True(t, p1.Id != "")
				assert.True(t, !p1.CreatedAt.AsTime().IsZero())

				p2 := &peoplepbv1.Person{
					Id:        p1.Id,
					IdNumbers: map[string]string{"pin": "pin1"},
					Labels:    []string{"tours:bike", "tours:hike", "kids"},
					Phones:    map[string]string{"mobile": "123123123"},
					Attr:      map[string]string{"food": "veg"},
					Gender:    "male",
				}

				err = testRepo.Update(ctx, p2)
				p3, err := testRepo.GetById(ctx, p1.Id)
				require.NoError(t, err)
				p3.CreatedAt = p1.CreatedAt
				assert.Equal(t, p3.Name, p1.Name)
				assert.Equal(t, p3.Emails, p1.Emails)
				assert.Equal(t, p3.FullName, p1.FullName)
				assert.Equal(t, p3.IdNumbers, p2.IdNumbers)
				assert.Equal(t, p3.Labels, p2.Labels)
				assert.Equal(t, p3.Phones, p2.Phones)
				assert.Equal(t, p3.Attr, p2.Attr)
				assert.Equal(t, p3.Gender, p2.Gender)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			tc.run(t)
		})
	}

}

func printList(name string, list []*peoplepbv1.Person) {
	fmt.Printf("%s: START---------------\n", name)
	for _, v := range list {
		fmt.Printf("%s: %#v\n", name, v)
	}
	fmt.Printf("%s: END.\n\n", name)
}

func printMap(name string, list map[string]*peoplepbv1.Person) {
	fmt.Printf("%s: START---------------\n", name)
	for _, v := range list {
		fmt.Printf("%s: %#v\n", name, v)
	}
	fmt.Printf("%s: END.\n\n", name)
}
