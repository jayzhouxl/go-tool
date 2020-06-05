package tool

import (
	"reflect"
	"strconv"
	"testing"
)

func TestSetStructField(t *testing.T) {
	el := struct {
		P0 int
		P1 uint32
		P2 string
		P3 string `abc:"P4"`
	}{}
	value := []interface{}{20, "20", int32(20), uint32(20)}
	for index, v := range value {
		for i := 0; i < 4; i++ {
			key := "P" + strconv.Itoa(i)
			err := SetStructField(&el, key, v)
			if err != nil {
				t.Fatalf("%dth set %s err:%s", index, key, err.Error())
			}
			if i == 0 && el.P0 != 20 {
				t.Errorf("%dth set %s err", index, key)
			}
			if i == 1 && el.P1 != uint32(20) {
				t.Errorf("%dth set %s err", index, key)
			}
			if i > 1 {
				if reflect.ValueOf(el).Field(i).String() != "20" {
					t.Errorf("%dth set %s err", index, key)
				}
			}
		}
	}

}
