package reflect

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm"
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

			if !val.IsValid() || val.IsZero() {
				return
			}

			if val.Kind() == reflect.Func {
				path += reflect.TypeOf(input).String()
				fmt.Printf("%s: %v\n", path, val.Type())
				return
			}

			// 如果输入是指针，则获取指针指向的元素
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
				return
			}

			// 处理非结构体
			if val.Kind() != reflect.Struct {
				if path != "" {
					fmt.Printf("%s: %v\n", path, val)
				} else {
					fmt.Printf("%v\n", val)
				}
				return
			}

			path += val.Type().Name() + "=>"

			// 遍历结构体的每个字段
			for i := 0; i < val.NumField(); i++ {
				// 获取字段的值
				valueField := val.Field(i)

				// 获取字段的类型
				typeField := val.Type().Field(i)

				// 处理结构体、指针、其他类型
				if !valueField.IsValid() || valueField.IsZero() {
					continue
				} else if valueField.Kind() == reflect.Struct {
					iprint(valueField.Interface(), path+typeField.Name)
				} else if valueField.Kind() == reflect.Pointer {
					iprint(valueField.Interface(), path+"*"+typeField.Name)
				} else if valueField.Kind() == reflect.Func {
					fmt.Printf("%s: %v\n", path+typeField.Name, valueField.Type().String())
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
		Zero int
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
		Zero: 0,
		Data: &struct {
			X, Y  int
			IData struct{ ID string }
		}{X: 1, Y: 2, IData: struct{ ID string }{ID: "123"}},
		f: func(x, y int) {
			fmt.Printf("hello, world")
		},
	})

	var x int = 10
	debug(x)

	var a Person
	debug(a)

	var f func() = func() {
		fmt.Printf("this is a func")
	}
	debug(f)

	var pf *func()
	debug(pf)
}

func FuncInCompany(qs *gorm.DB, query interface{}, results interface{}) (count *int64, err error) {
	var selectFields []string
	var expandFields []string
	var orderFields []string
	filterFields := make(map[string]interface{})
	limit := 20
	skip := 0
	countEnable := false
	sqlKey := ""
	sql := ""
	var sqlValue []interface{}

	values := reflect.ValueOf(query)
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		v := values.Field(i)
		if !v.IsValid() || (v.IsZero()) {
			continue
		}

		k := types.Field(i).Name
		tag := types.Field(i).Tag
		if k == "Count" {
			if v.Type().String() != "bool" {
				continue
			}
			countEnable = v.Bool()
		} else if k == "Select" {
			if v.Type().String() != "string" {
				continue
			}
			selectFields = strings.Split(v.String(), ",")
		} else if k == "Expand" {
			if v.Type().String() != "string" {
				continue
			}
			expandFields = strings.Split(v.String(), ",")
		} else if k == "Order" {
			if v.Type().String() != "string" {
				continue
			}
			orderFields = strings.Split(v.String(), ",")
		} else if k == "Limit" {
			if v.Type().String() != "int64" {
				continue
			}
			limit = int(v.Int())
		} else if k == "Skip" {
			if v.Type().String() != "int64" {
				continue
			}
			skip = int(v.Int())
		} else {
			if tag.Get("column") == "" {
				continue
			}
			specialFields := strings.Split(k, "__")
			if len(specialFields) > 1 {
				k = specialFields[0]
				field := specialFields[1]

				if tag.Get("column") != "" {
					k = tag.Get("column")
				}

				if len(sqlKey) > 0 {
					sqlKey += " AND "
				}
				if field == "Contains" {
					if v.Type().String() != "string" {
						continue
					}
					sqlKey += k + " LIKE ?"
					sqlValue = append(sqlValue, "%"+v.String()+"%")
				} else if field == "In" {
					sqlKey += k + " IN ?"
					values := strings.Split(v.String(), ",")
					sqlValue = append(sqlValue, values)
				} else if field == "Gte" {
					sqlKey += k + " >= ?"
					sqlValue = append(sqlValue, v.Interface())
				} else if field == "Lte" {
					sqlKey += k + " <= ?"
					sqlValue = append(sqlValue, v.Interface())
				} else if field == "Gt" {
					sqlKey += k + " > ?"
					sqlValue = append(sqlValue, v.Interface())
				} else if field == "Lt" {
					sqlKey += k + " < ?"
					sqlValue = append(sqlValue, v.Interface())
				} else if field == "Not" {
					sqlKey += k + " != ?"
					sqlValue = append(sqlValue, v.Interface())
				} else if field == "IsNull" {
					if v.Type().String() != "string" {
						continue
					}
					if v.String() == "true" {
						sql += k + " IS NULL"
					} else {
						sql += k + " IS NOT NULL"
					}
				}
			} else {
				if tag.Get("column") != "" {
					k = tag.Get("column")
				}
				filterFields[k] = v.Interface()
			}
		}
	}
	if len(orderFields) < 1 {
		orderFields = append(orderFields, "created_at desc")
	}
	if limit > 500 {
		limit = 500
	}
	if len(selectFields) > 0 {
		qs = qs.Select(selectFields)
	}
	for _, v := range orderFields {
		qs = qs.Order(v)
	}
	for _, v := range expandFields {
		qs = qs.Preload(v)
	}
	qs = qs.Where(filterFields)
	if len(sqlKey) > 0 {
		qs = qs.Where(sqlKey, sqlValue...)
	}
	if len(sql) > 0 {
		qs = qs.Where(sql)
	}
	if countEnable {
		count = new(int64)
		err = qs.Count(count).Error
	}

	err = qs.Limit(limit).Offset(skip).Find(results).Error

	return
}

func TestFuncInCompany(t *testing.T) {
	var in *func()
	var result []string
	FuncInCompany(nil, in, &result)
}

func TestPtrFunc(t *testing.T) {
	var f *func()
	v := reflect.ValueOf(f)
	if !v.IsNil() {
		e := v.Elem()
		t.Logf("%v\n", e)
	} else {
		t.Logf("unknown\n")
	}
}

// IsAssigned 检查对象是否有被赋值，不检查非导出字段
func IsAssigned(s interface{}) bool {
	v := reflect.ValueOf(s)
	kind := v.Kind()

	switch kind {
	case reflect.Func:
		return v.IsNil()
	case reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return IsAssigned(v.Elem().Interface())
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if !field.CanInterface() {
				continue
			}
			if !IsAssigned(field.Interface()) {
				return false
			}
		}
		return true
	default:
		return v.IsZero()
	}
}

func TestIsStructEmpty(t *testing.T) {
	type ReqClass1 struct {
		Arg11 int
		Arg12 string
		Fn    func()
		T     time.Time
	}

	type ReqClass2 struct {
		Arg21 int
		Arg22 string
		*ReqClass1
	}

	type Req struct {
		ReqClass1
		ReqClass2
	}

	// 各测试用例，方便起见，先不写成 table-driven 的形式
	passed := IsAssigned(&Req{}) == true
	passed = passed && IsAssigned(&Req{
		ReqClass1: ReqClass1{
			Arg11: 1,
			Arg12: "123",
		},
	}) == false
	passed = passed && IsAssigned(&Req{
		ReqClass1: ReqClass1{
			Fn: nil,
		},
	}) == true
	passed = passed && IsAssigned(&Req{
		ReqClass1: ReqClass1{
			Fn: func() {
				fmt.Println("hello world")
			},
		},
	}) == false
	passed = passed && IsAssigned(&Req{
		ReqClass1: ReqClass1{
			Arg11: 1,
		},
		ReqClass2: ReqClass2{
			Arg21: 1,
		},
	}) == false
	passed = passed && IsAssigned(&Req{
		ReqClass2: ReqClass2{
			ReqClass1: &ReqClass1{
				Arg11: 2,
			},
		},
	}) == false
	passed = passed && IsAssigned(&Req{
		ReqClass1: ReqClass1{
			Arg11: 1,
		},
		ReqClass2: ReqClass2{
			Arg21: 1,
			ReqClass1: &ReqClass1{
				Arg11: 2,
			},
		},
	}) == false
	passed = passed && IsAssigned(&Req{
		ReqClass1: ReqClass1{
			T: time.Now(),
		},
	}) == true

	// 判断是否通过测试
	if passed {
		t.Logf("all tests passed :)")
	} else {
		t.Errorf("tests not passed :(")
	}
}
