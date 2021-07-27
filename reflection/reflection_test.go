package reflection

import (
	"reflect"
	"testing"
)

type Person struct {
	Name    string
	Profile Profile
}

type Profile struct {
	Age  int
	City string
}

func TestWalk(t *testing.T) {
	cases := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{"Struct with two string field",
			struct {
				Name string
				City string
			}{"Manu", "London"},
			[]string{"Manu", "London"},
		},
		{
			"struct with non string field",
			struct {
				Name string
				Age  int
			}{"Manu", 00},
			[]string{"Manu"},
		},
		{
			"nested fields",
			Person{
				"Manu",
				Profile{00, "London"},
			},
			[]string{"Manu", "London"},
		},
		{
			"Pointers to things",
			&Person{
				"Manu",
				Profile{00, "London"},
			},
			[]string{"Manu", "London"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			var got []string
			walk(test.Input, func(input string) {
				got = append(got, input)
			})
			if !reflect.DeepEqual(got, test.ExpectedCalls) {
				t.Errorf("got %q want %q", got, test.ExpectedCalls)
			}
		})
	}
}