package model

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestNullJson(t *testing.T) {
	type MyCustom struct {
		F1 string `json:"f1"`
		F2 int    `json:"f2"`
	}
	type MyStruct struct {
		ID        int64           `json:"id"`
		Age       int             `json:"age"`
		Hobby     Null[string]    `json:"hobby"`
		Create    Null[time.Time] `json:"create"`
		Price     Null[float64]   `json:"price"`
		MyCustom  Null[MyCustom]  `json:"my_custom"`
		MyCustom2 MyCustom        `json:"myCustom2"`
	}
	myStruct := MyStruct{
		ID:     1,
		Age:    18,
		Hobby:  NewNull[string]("football", true),
		Create: NewNull[time.Time](time.Now(), true),
		Price:  NewNull[float64](9.99, false),
	}
	fmt.Printf("%v\n", myStruct)
	jsonData, err := json.Marshal(myStruct)
	if err != nil {
		t.Errorf("Error marshaling JSON: %v", err)
	}
	t.Logf("JSON Data: %s", jsonData)
}
