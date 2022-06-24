package jsonmask

import (
	"fmt"
	"reflect"
	"testing"
)

type address struct {
	Nation   string `json:"nation"`
	District string `json:"district"`
	Street   string `json:"street"`
}

type person struct {
	Name     string   `json:"name"`
	ID       string   `json:"id"`
	Age      int32    `json:"age"`
	Gender   string   `json:"gender"`
	Email    string   `json:"email"`
	PhoneNos []string `json:"phoneNos"`

	Address *address `json:"address"`
}

// 一般来说，这个方法不需要和 person 类型本身位于同一个源码文件，在同一个包下即可。
func (p *person) NeedMaskFields() map[string]MaskHandler {
	return map[string]MaskHandler{
		".name":             MaskHandlerFunc(Name),
		".id":               MaskHandlerFunc(ToAsterisks),
		".phoneNos":         MaskHandlerFunc(Telephone),
		".address.district": MaskHandlerFunc(ToAsterisks),
		".address.street":   MaskHandlerFunc(ToEmpty),
	}
}

func TestExample(t *testing.T) {
	p := &person{
		Name:     "Joey",
		ID:       "125125125",
		Age:      26,
		Gender:   "male",
		Email:    "joey@catchmeifyoucan.com",
		PhoneNos: []string{"12244445555", "12244446666"},
		Address:  &address{Nation: "China", District: "GD", Street: "NS"},
	}

	bytes, err := JsonMasked(p)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bytes))
}

func TestJsonMasked(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{{
		name: "1",
		args: args{&person{
			Name:     "Joey",
			ID:       "125125125",
			Age:      26,
			Gender:   "male",
			Email:    "joey@catchmeifyoucan.com",
			PhoneNos: []string{"12244445555", "12244446666"},
			Address:  &address{Nation: "China", District: "GD", Street: "NS"},
		}},
		want:    []byte(`{"address":{"district":"**","nation":"China","street":""},"age":26,"email":"joey@catchmeifyoucan.com","gender":"male","id":"*********","name":"J*y","phoneNos":["122****5555","122****6666"]}`),
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JsonMasked(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonMasked() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonMasked() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaskMapOrSlice(t *testing.T) {
	needMaskFields := map[string]MaskHandler{
		".name":             MaskHandlerFunc(Name),
		".address.district": MaskHandlerFunc(ToAsterisks),
		".address.street":   MaskHandlerFunc(ToEmpty),
	}
	rawData := map[string]interface{}{
		"name": "joey",
		"address": map[string]interface{}{
			"nation":   "China",
			"district": "GD",
			"street":   "NS",
		},
	}
	rawData2 := map[string]interface{}{
		"name": "joey",
		"address": map[string]interface{}{
			"nation":   "China",
			"district": "GD",
			"street":   "NS",
		},
	}
	wantMap := map[string]interface{}{
		"name": "j*y",
		"address": map[string]interface{}{
			"nation":   "China",
			"district": "**",
			"street":   "",
		},
	}

	type args struct {
		needMaskFields map[string]MaskHandler
		rawData        interface{}
	}
	tests := []struct {
		name           string
		args           args
		wantMaskedData interface{}
	}{{
		name:           "map",
		args:           args{needMaskFields, rawData},
		wantMaskedData: wantMap,
	}, {
		name:           "slice",
		args:           args{needMaskFields, []interface{}{rawData2}},
		wantMaskedData: []interface{}{wantMap},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMaskedData := MaskMapOrSlice(tt.args.needMaskFields, tt.args.rawData); !reflect.DeepEqual(gotMaskedData, tt.wantMaskedData) {
				t.Errorf("MaskMapOrSlice() = %v, want %v", gotMaskedData, tt.wantMaskedData)
			}
		})
	}
}
