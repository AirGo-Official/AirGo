package other_plugin

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// 约定前端需要匹配的类型,暂时约定4个类型，date->日期，string->字符串，num->数字，bool->布尔
var typeComparisonTable = map[string]string{
	"time.Time":  "date",
	"*time.Time": "date",

	"string":    "string",
	"uuid.UUID": "string",

	"int32":   "num",
	"int64":   "num",
	"int":     "num",
	"float64": "num",
	"float32": "num",

	"bool": "bool",
}

// 对长度不足n的数字后面补0
func Sup(i, n int64) string {
	m := fmt.Sprintf("%d", i)
	for int64(len(m)) < n {
		m = fmt.Sprintf("%s0", m)
	}
	return m
}

// struct转map
func StructToMap(data interface{}) map[string]interface{} {

	m := make(map[string]interface{})

	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem() //需要用struct的Type，不能用指针的Type
	}
	if v.Kind() != reflect.Struct {
		return m
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {

		//field.Name,            //变量名称
		//field.Offset,          //相对于结构体首地址的内存偏移量，string类型会占据16个字节
		//field.Anonymous,       //是否为匿名成员
		//field.Type,            //数据类型，reflect.Type类型
		//field.IsExported(),    //包外是否可见（即是否以大写字母开头）
		//field.Tag.Get("json")) //获取成员变量后面``里面定义的tag

		name := t.Field(i).Name
		tag := t.Field(i).Tag.Get("json")
		if tag == "-" || name == "-" {
			continue
		}
		if tag != "" {
			index := strings.Index(tag, ",")
			if index == -1 {
				name = tag
			} else {
				name = tag[:index]
			}
		}
		m[name] = v.Field(i).Interface()
	}
	return m
}

// m1=获取结构体字段英文名和中文名map，如  name:名字
// m2=获取结构体字段类型map，如  name:string
// m3=获取结构体字段数组，如  [name:aaa,age:18]
func GetStructFieldMap(data interface{}) ([]string, map[string]interface{}, map[string]interface{}) {

	m1 := make([]string, 0)
	m2 := make(map[string]interface{})
	m3 := make(map[string]interface{})
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem() //需要用struct的Type，不能用指针的Type
	}
	if v.Kind() != reflect.Struct {
		return m1, m2, m3
	}
	t := v.Type() //Value转Type
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Struct {
			//如果没有指定gorm:"embedded",则跳过递归
			comment := t.Field(i).Tag.Get("gorm")
			index := strings.Index(comment, "embedded")
			if index == -1 {
				continue
			}
			//嵌套结构体，递归
			mm1, mm2, mm3 := GetStructFieldMap(v.Field(i).Interface())
			for _, v := range mm1 {
				m1 = append(m1, v)
			}
			for k, v := range mm2 {
				if _, ok := m2[k]; !ok {
					m2[k] = v
				}
			}
			for k, v := range mm3 {
				if _, ok := m3[k]; !ok {
					m3[k] = v
				}
			}
			continue
		}
		//json:字段
		jsonName := t.Field(i).Tag.Get("json")
		if jsonName == "-" {
			continue
		}
		index := strings.Index(jsonName, ",")
		if index != -1 {
			jsonName = jsonName[:index]
		}
		//fmt.Println("jsonName:", jsonName)
		//gorm:通过gorm备注获得中文名
		comment := t.Field(i).Tag.Get("gorm")
		if comment == "-" {
			continue
		}
		index = strings.Index(comment, "comment:")
		if index == -1 {
			continue
		}
		comment = comment[index+8:]
		index = strings.Index(comment, ";")
		if index != -1 {
			comment = comment[:index]
		}
		//fmt.Println("comment:", comment)
		m2[jsonName] = comment
		if mv, ok := typeComparisonTable[v.Field(i).Type().String()]; ok {
			m3[jsonName] = mv
		} else {
			m3[jsonName] = typeComparisonTable["string"]
		}
		m1 = append(m1, jsonName)

	}
	//fmt.Println(m)
	return m1, m2, m3

}

// 数组去重
func ArrayDeduplication(slice []int64) []int64 {
	tempMap := make(map[int64]struct{}, len(slice))
	j := 0
	for _, v := range slice {
		_, ok := tempMap[v]
		if ok {
			continue
		}
		tempMap[v] = struct{}{}
		slice[j] = v
		j++
	}
	return slice[:j]
}

// 判断字符串是否在一个数组中
func In(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}

// 数组拆分
func SplitArray[T any](arr []T, num int64) [][]T {

	max := int64(len(arr))
	if max < num {
		return nil
	}
	var segmens = make([][]T, 0)
	quantity := max / num
	end := int64(0)
	for i := int64(1); i <= num; i++ {
		qu := i * quantity
		if i != num {
			segmens = append(segmens, arr[i-1+end:qu])
		} else {
			segmens = append(segmens, arr[i-1+end:])
		}
		end = qu - i
	}
	return segmens
}
