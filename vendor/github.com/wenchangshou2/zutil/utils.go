package zutil

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func GetFullPath(path string) (string, error) {
	fullexecpath, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir, _ := filepath.Split(fullexecpath)
	return filepath.Join(dir, path), nil
}
func stringArrayContains(needle string, haystack []string) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}

func ScanStructIntoMap(obj interface{}) (map[string]interface{}, error) {
	dataStruct := reflect.Indirect(reflect.ValueOf(obj))
	if dataStruct.Kind() != reflect.Struct {

		return nil, errors.New("expected a pointer to a struct")
	}

	dataStructType := dataStruct.Type()

	mapped := make(map[string]interface{})

	for i := 0; i < dataStructType.NumField(); i++ {
		field := dataStructType.Field(i)
		fieldv := dataStruct.Field(i)
		fieldName := field.Name
		bb := field.Tag
		sqlTag := bb.Get("sql")
		sqlTags := strings.Split(sqlTag, ",")
		var mapKey string

		inline := false

		if bb.Get("beedb") == "-" || sqlTag == "-" || reflect.ValueOf(bb).String() == "-" {
			continue
		} else if len(sqlTag) > 0 {
			//TODO: support tags that are common in json like omitempty
			if sqlTags[0] == "-" {
				continue
			}
			mapKey = sqlTags[0]
		} else {
			mapKey = fieldName
		}

		if len(sqlTags) > 1 {
			if stringArrayContains("inline", sqlTags[1:]) {
				inline = true
			}
		}

		if inline {
			// get an inner map and then put it inside the outer map
			map2, err2 := ScanStructIntoMap(fieldv.Interface())
			if err2 != nil {
				return mapped, err2
			}
			for k, v := range map2 {
				mapped[k] = v
			}
		} else {
			value := dataStruct.FieldByName(fieldName).Interface()
			mapped[mapKey] = value
		}
	}

	return mapped, nil
}
