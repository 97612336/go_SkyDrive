package index

import (
	"net/http"
	"go_SkyDrive/config"
	"io/ioutil"
	"html/template"
	"go_SkyDrive/models"
)

//主页面跳转
func Index(w http.ResponseWriter, r *http.Request) {
	//读取目录下的所有文件
	static_dir := config.Static_Path
	files, _ := ioutil.ReadDir(static_dir)
	//定义所有文件组成的数组
	var file_list []map[string]string
	for _, v := range files {
		var one_file = make(map[string]string)
		one_file["name"]=v.Name()
		one_file["href"]=config.Static_href+v.Name()
		file_list = append(file_list, one_file)
	}
	//渲染模板
	tmpl, _ := template.ParseFiles("template/index.html")
	//创建结果字典
	var Res models.Index
	Res.Files = file_list
	tmpl.Execute(w, Res)
	return
}


