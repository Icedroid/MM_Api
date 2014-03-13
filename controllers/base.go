package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Icedroid/MM_Api/models"
	"github.com/Icedroid/MM_Api/services"
	"util/logs"

	"github.com/astaxie/beego"
	"labix.org/v2/mgo"
)

type NestPreparer interface {
	NestPrepare()
}

// baseRouter implemented global settings for all other routers.
type BaseController struct {
	beego.Controller
	//	i18n.Locale
	m *models.M
}

var (
	app             *models.App
	c               *models.Collection
	isNewCollection bool
	newFieldType    map[string]int
	mongoRow        map[string]interface{}
	mongoIndex      mgo.Index

	bodyJsonData = make(map[string]interface{})
	res          = make(map[string]interface{})
	sessionToken string
)

// Prepare implemented Prepare method for baseRouter.
func (bc *BaseController) Prepare() {
	appKey := strings.TrimSpace(bc.Ctx.Input.Header(models.AppConfig.String("appkeyheadername")))
	restKey := strings.TrimSpace(bc.Ctx.Input.Header(models.AppConfig.String("restapikeyheadername")))
	logs.Logger.Debugf("Header=%v, Params=%v", bc.Ctx.Request.Header, bc.Ctx.Input.Params)
	logs.Logger.Debugf("AppKey=%s, RestKey=%s", appKey, restKey)
	if "" == appKey || "" == restKey {
		bc.unauthorized()
	}
	if bc.Ctx.Input.Is("POST") || bc.Ctx.Input.Is("PUT") {
		contentType := strings.TrimSpace(bc.Ctx.Input.Header("Content-Type"))
		if "application/json" != contentType {
			bc.responseError("contentType", contentType)
		}
		var f interface{}
		var ok bool
		body := bc.Ctx.Input.Body()
		err := json.Unmarshal(body, &f)
		if err != nil {
			bc.responseError("json", string(body))
		}
		bodyJsonData, ok = f.(map[string]interface{})
		if !ok {
			bc.responseError("json", string(body))
		}
	}

	var err error
	app, err = models.NewApp()
	if err != nil {
		logs.Logger.Errorf("models.NewApp init error: %s", err)
		bc.busy()
	}
	app.AppKey = appKey
	err = app.RSet()
	logs.Logger.Debugf("App=%+v", app)
	if err != nil {
		logs.Logger.Errorf("app.Rset error: %s", err)
		bc.busy()
	}
	if !app.StatusNormal() || !app.RestKeyNormal(restKey) {
		logs.Logger.Errorf("app status '%d' or restkey '%s' incorrect", app.Status, app.RestKey)
		bc.unauthorized()
	}
	if c, ok := bc.AppController.(NestPreparer); ok {
		c.NestPrepare()
	}
}

// Finish is called once the controller method completes
func (bc *BaseController) Finish() {
	defer closeAll()
	logs.Logger.Debugf("Finish %s", bc.Ctx.Request.URL.Path)
}

func (bc *BaseController) response() {
	bc.Data["json"] = res
	bc.ServeJson()
}

func (bc *BaseController) responseError(str string, a ...interface{}) {
	defer closeAll()
	s := services.StateList[str]
	code := s.Code
	err := fmt.Sprintf(s.Msg, a...)
	bc.Ctx.Output.Header("Content-Type", "application/json;charset=UTF-8")
	bc.Ctx.Output.SetStatus(404)
	bc.Data["json"] = map[string]interface{}{"code": code, "error": err}
	bc.ServeJson()
	bc.StopRun()
}

func (bc *BaseController) unauthorized() {
	defer closeAll()
	bc.Ctx.Output.Header("Content-Type", "application/json;charset=UTF-8")
	bc.Ctx.Output.SetStatus(401)
	bc.Data["json"] = map[string]string{"error": "unauthorized"}
	bc.ServeJson()
	bc.StopRun()
}

func (bc *BaseController) busy() {
	defer closeAll()
	bc.Ctx.Output.Header("Content-Type", "application/json;charset=UTF-8")
	bc.Ctx.Output.SetStatus(500)
	bc.Data["json"] = map[string]string{"error": "It is busy...Try it later!"}
	bc.ServeJson()
	bc.StopRun()
}

func closeAll() {
	if app != nil {
		app.CloseAll()
	}
	if c != nil {
		c.CloseAll()
	}
	logs.Logger.Flush()
//	if err := recover(); err != nil {
//		logs.Logger.Errorf("controller exit get error", err)
//	}
}
