package zutil

import "encoding/base64"

func Base64Encode(str string) string{
	data := []byte(str)
	return base64.StdEncoding.EncodeToString(data)
}
func Base64Decode(str string) (string,error){
	data,err:=base64.StdEncoding.DecodeString(str)
	if err!=nil{
		return "",err
	}
	return string(data),nil
}