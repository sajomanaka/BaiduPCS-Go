// Package pcsweb web前端包
package pcsweb

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/iikira/BaiduPCS-Go/baidupcs"
	"github.com/iikira/BaiduPCS-Go/pcsconfig"
	"html/template"
	"net/http"
)

var (
	activeAPI    *baidupcs.PCSApi
	staticBox    *rice.Box // go.rice 文件盒子
	templatesBox *rice.Box // go.rice 文件盒子
)

func boxInit() (err error) {
	staticBox, err = rice.FindBox("static")
	if err != nil {
		return
	}

	templatesBox, err = rice.FindBox("template")
	if err != nil {
		return
	}

	return nil
}

func activeAPIInit() (err error) {
	// 获取当前登录的用户
	activeUser, err := pcsconfig.Config.GetActive()
	if err != nil {
		return err
	}

	activeAPI = baidupcs.NewPCS(activeUser.BDUSS)
	return nil
}

// StartServer 开启web服务
func StartServer(port uint) error {
	if port <= 0 || port > 65535 {
		return fmt.Errorf("invalid port %d", port)
	}

	err := boxInit()
	if err != nil {
		return err
	}

	http.HandleFunc("/", rootMiddleware)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticBox.HTTPBox())))
	http.HandleFunc("/about.html", middleware(aboutPage))
	http.HandleFunc("/index.html", middleware(indexPage))
	http.HandleFunc("/cgi-bin/baidu/pcs/file/list", activeAuthMiddleware(fileList))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func aboutPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(templatesBox.MustString("index.html"))
	checkErr(err)

	_, err = tmpl.Parse(templatesBox.MustString("about.html"))
	checkErr(err)

	err = tmpl.Execute(w, nil)
	checkErr(err)
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	tmpl, err := template.New("index").Parse(templatesBox.MustString("index.html"))
	checkErr(err)

	_, err = tmpl.Parse(templatesBox.MustString("baidu/userinfo.html"))
	checkErr(err)

	err = tmpl.Execute(w, r.Form.Get("path"))
	checkErr(err)
}
