package util

import (
	"encoding/json"
	"fmt"
)

func Get_json_string(m interface{})  string  {
	res,err:=json.Marshal(m)
	if err!=nil{
		fmt.Println(err)
		return ""
	}
	res_str:=string(res)
	return res_str
}




