package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net"
	"os"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	BasePath string
)

func init() {
	BasePath, _ = os.Getwd()
}

func ParseIP(s string) (net.IP, int) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, 0
	}

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return ip, 4
		case ':':
			return ip, 6
		}
	}

	return nil, 0
}

func MD5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func structToMap(obj interface{}, lowerKey bool) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		if lowerKey {
			name = strings.ToLower(name)
		}
		data[name] = v.Field(i).Interface()
	}

	return data
}

func StructToMap(obj interface{}) map[string]interface{} {
	return structToMap(obj, false)
}

func StructToMapWithLowerKey(obj interface{}) map[string]interface{} {
	return structToMap(obj, true)
}

func StructCopyFields(a interface{}, b interface{}, fields ...string) error {
	at := reflect.TypeOf(a)
	av := reflect.ValueOf(a)
	bt := reflect.TypeOf(b)
	bv := reflect.ValueOf(b)

	if at.Kind() != reflect.Ptr {
		return errors.New("a must be a struct pointer")
	}
	av = reflect.ValueOf(av.Interface())

	_fields := make([]string, 0)
	if len(fields) > 0 {
		_fields = fields
	} else {
		for i := 0; i < bv.NumField(); i++ {
			_fields = append(_fields, bt.Field(i).Name)
		}
	}

	for i := 0; i < len(_fields); i++ {
		name := _fields[i]
		f := av.Elem().FieldByName(name)
		bValue := bv.FieldByName(name)

		// a 中有同名的字段并且类型一致才复制
		if f.IsValid() && f.Type() == bValue.Type() {
			f.Set(bValue)
		}
	}
	return nil
}

func GetStructFieldName(s interface{}) ([]string, error) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, errors.New("not Struct")
	}

	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}

	return result, nil
}

func GetStructFieldNameToSnake(s interface{}) ([]string, error) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, errors.New("not Struct")
	}

	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, CamelToSnake(t.Field(i).Name))
	}

	return result, nil
}

func GetStructTags(v interface{}, tagName string) []string {
	val := reflect.TypeOf(v)
	var tags []string

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := field.Tag.Get(tagName)
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags
}
