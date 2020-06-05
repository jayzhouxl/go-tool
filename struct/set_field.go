package tool

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

//根据字段名往结构体中填充值
//字段名可以是结构体字段名也可以是标签值
//无须传入标签字段名
//同时工具还会自动转换数据类型，使用者无须关心字段类型
//目前仅支持int相关(uint,int32...)和string
func SetStructField(obj interface{}, key string, value interface{}) error {
	optr := reflect.ValueOf(obj)
	if optr.Kind() != reflect.Ptr || optr.IsNil() {
		return fmt.Errorf("struct object must ptr")
	}
	o := optr.Elem()
	if o.Kind() != reflect.Struct {
		return fmt.Errorf("object is not vaild struct")
	}

	oField := o.FieldByName(key)
	if !oField.IsValid() { //如果根据field名字没找到
		kName := getFieldNameByTag(o.Type(), key) //再根据tag找
		if len(kName) > 0 {
			oField = o.FieldByName(kName)
		} else {
			return fmt.Errorf("%s is not in struct field", key)
		}
	}
	if oField.CanSet() {
		val := reflect.ValueOf(value)
		if !isSupportKind(val.Kind()) { //目前仅支持简单类型
			return fmt.Errorf("now just support simple kind field")
		}
		if val.Kind() == oField.Kind() {
			oField.Set(val)
		} else {
			val = getVaildValue(oField.Kind(), val)
			if val.IsValid() {
				oField.Set(val.Convert(oField.Type()))
			}
		}
	} else {
		return fmt.Errorf("%s can not set,please ensure it is addressable and was not obtained by the use of unexported struct fields", key)
	}
	return nil
}

func getFieldNameByTag(ot reflect.Type, key string) string {
	var name string
	for i := 0; i < ot.NumField(); i++ {
		if ot.Field(i).Anonymous {
			name = getFieldNameByTag(ot.Field(i).Type, key)
		} else if ot.Field(i).Tag != "" {
			if strings.Trim(strings.Split(strings.TrimSpace(string(ot.Field(i).Tag)), ":")[1], `"`) == key {
				name = ot.Field(i).Name
				break
			}
		}
	}
	return name
}

func isSupportKind(k reflect.Kind) bool {
	supportKind := map[reflect.Kind]bool{
		reflect.String: true,
		reflect.Int:    true,
		reflect.Int8:   true,
		reflect.Int16:  true,
		reflect.Int32:  true,
		reflect.Int64:  true,
		reflect.Uint:   true,
		reflect.Uint8:  true,
		reflect.Uint16: true,
		reflect.Uint32: true,
		reflect.Uint64: true,
	}
	return supportKind[k]
}

func getVaildValue(dstKind reflect.Kind, val reflect.Value) reflect.Value {
	switch dstKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val.Kind() != reflect.String { //value是int相关类型
			return val
		} else { //value是string
			v, err := strconv.ParseFloat(val.String(), 32)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s convert err:%s\n", val.String(), err.Error())
				return reflect.ValueOf(nil)
			}
			return reflect.ValueOf(v)
		}
	case reflect.String:
		//value 只能是int那些类型
		return reflect.ValueOf(fmt.Sprintf("%d", val.Interface()))
	}
	return reflect.ValueOf(nil)
}
