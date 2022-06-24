package personpb

import (
	"fmt"
	jsonmask "github.com/JoeyLearnsToCode/json-mask"
	"testing"
)

func TestPersonJsonMask(t *testing.T) {
	p := &Person{
		Name:      "Joey",
		Age:       28,
		Telephone: "15599998888",
		Id:        "U151987981749871",
		Address: &Address{
			Nation:   "China",
			District: "GD",
			Street:   "NS",
		},
	}

	fmt.Println(string(jsonmask.JsonMaskedQuietly(p)))
}
