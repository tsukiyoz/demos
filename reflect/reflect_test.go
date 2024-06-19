package reflect

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReflectBasicUsage(t *testing.T) {
	var x float64 = 3.4

	val, typ := reflect.ValueOf(x), reflect.TypeOf(x)

	t.Logf("type: %v\n", typ)
	t.Logf("value: %v\n", val)
	t.Logf("Kind is float64: %v\n", typ.Kind() == reflect.Float64)
	t.Logf("Type is float64: %v\n", typ == reflect.TypeOf(float64(0)))
}

func TestReflectSetValue(t *testing.T) {
	var x float64 = 3.4

	t.Logf("Original value is %v\n", x)

	reflect.ValueOf(&x).Elem().SetFloat(7.1) // x必须可寻址的

	t.Logf("Modified value is %v\n", x)
}

func TestReflectInTagAnalyse(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	p := &Person{
		Name: "tsukiyo",
		Age:  30,
	}

	typ := reflect.TypeOf(*p)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		t.Logf("Field: %s, Tag: %s\n", field.Name, field.Tag.Get("json"))
	}
}

func TestBuildGeneralFunction(t *testing.T) {
	debug := func(input interface{}) {
		var iprint func(input interface{}, path string)
		iprint = func(input interface{}, path string) {
			val := reflect.ValueOf(input)
			path += val.Type().Name() + "=>"

			if val.Kind() == reflect.Func {
				fmt.Printf("this is a func")
				return
			}

			// 如果输入是指针，则获取指针指向的元素
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}

			// 处理非结构体
			if val.Kind() != reflect.Struct {
				fmt.Printf("%s: %v\n", path, val)
			}

			// 遍历结构体的每个字段
			for i := 0; i < val.NumField(); i++ {
				// 获取字段的值
				valueField := val.Field(i)

				// 获取字段的类型
				typeField := val.Type().Field(i)

				// 处理结构体、指针、其他类型
				if valueField.Kind() == reflect.Struct {
					iprint(valueField.Interface(), path+typeField.Name)
				} else if valueField.Kind() == reflect.Pointer {
					iprint(valueField.Interface(), path+"*"+typeField.Name)
				} else if valueField.Kind() == reflect.Func {
					fmt.Printf("%s: %v\n", path+typeField.Name, valueField.Type())
				} else {
					fmt.Printf("%s: %v\n", path+typeField.Name, valueField.Interface())
				}
			}
		}
		path := ""
		iprint(input, path)
	}

	type Person struct {
		Name string
		Age  int
		Data *struct {
			X, Y  int
			IData struct {
				ID string
			}
		}
		f func(x, y int)
	}

	debug(Person{
		Name: "tsukiyo",
		Age:  18,
		Data: &struct {
			X, Y  int
			IData struct{ ID string }
		}{X: 1, Y: 2, IData: struct{ ID string }{ID: "123"}},
		f: func(x, y int) {
			fmt.Printf("hello, world")
		},
	})
}
