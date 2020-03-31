package zutil

import (
	"reflect"
	"strconv"
)

func MapToString(arguments map[string]interface{})string{
	str:=""
	for k,v:=range arguments{
		s:=""
		v:=reflect.ValueOf(v)
		switch v.Kind(){
		case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
			s+=strconv.FormatInt(v.Int(),10)
		case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64,reflect.Uintptr:
			s+=strconv.FormatUint(v.Uint(),10)
		case reflect.Bool:
			s+=strconv.FormatBool(v.Bool())
		case reflect.String:
			s+=strconv.Quote(v.String())
		case reflect.Float32,reflect.Float64:
			s+=strconv.Itoa(int(v.Float()))
		}
		if len(s)>0{
			str+=" -"+k+" "+s
		}
	}
	return str
}