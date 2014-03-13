package models

import (
	"errors"

	"github.com/garyburd/redigo/redis"
)

type App struct {
	M
	ID          int64
	Name        string
	AppKey      string
}

func NewApp() (app *App, err error) {
	app = new(App)
	app.redisDbIndex, _ = AppConfig.Int("redisappdb")
	app.redisHashPre = AppConfig.String("redisapphashpre")
	err = app.Init()
	return
}

//set app info from redis
func (app *App) RSet() error {
	if "" == app.AppKey {
		return errors.New("AppKey is empty")
	}
	hashName := app.redisHashPre + app.AppKey
	reply, err := redis.Values(app.redis.Do("HMGET", hashName, "id", "status"))
	if nil != err {
		return err
	}
	_, err = redis.Scan(reply, &app.ID, &app.Status)
	if nil != err {
		return err
	}
	return nil
}
