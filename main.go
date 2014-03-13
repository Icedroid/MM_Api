package main

import (
	"html/template"
	"net/http"

	_ "github.com/Icedroid/MM_Api/routers"
	"github.com/Icedroid/MM_Api/models"
	"util/logs"

	"github.com/astaxie/beego"
)

const (
	APP_VER = "1.1.0"
)

func main() {
	defer func() {
		models.RedisPool.Close()
		if err := recover(); err != nil {
			logs.Logger.Errorf("main get error", err)
		}
	}()
	beego.Errorhandler("404", page_not_found)
	beego.Run()
}

func init() {
	logs.Logger.Infof("MM_API- %s Used Beego - %s", APP_VER, beego.VERSION)
}

func page_not_found(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./" + beego.ViewsPath + "/404.html")
	//    t, _ := template.New("beegoerrortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Page Not Found"
	data["Content"] = template.HTML("<br>The Page You have requested flown the coop." +
		"<br>Perhaps you are here because:" +
		"<br><br><ul>" +
		"<br>The page has moved" +
		"<br>The page no longer exists" +
		"<br>You were looking for your puppy and got lost" +
		"<br>You like 404 pages" +
		"</ul>")
	data["Version"] = APP_VER
	rw.WriteHeader(http.StatusNotFound)
	t.Execute(rw, data)
}
