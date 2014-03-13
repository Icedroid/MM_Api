package models

import (
	"database/sql"
	"errors"
	"time"

	"util/logs"
	"github.com/astaxie/beego/config"
	"github.com/garyburd/redigo/redis"
)

var (
	RedisPool *redis.Pool
)

type M struct {
	redis           redis.Conn
	redisHashPre   string
	redisDbIndex   int
}

var AppConfig, _ = config.NewConfig("ini", "conf/app.conf")

func init() {
	server := AppConfig.String("redisserver")
	RedisPool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 10 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", server, time.Second*5, time.Second*10, time.Second*10)
			if err != nil {
				logs.Logger.Errorf("Dail master redis server %s %v", server, err)
				return nil, err
			}
			if _, err := c.Do("PING"); err != nil {
				c.Close()
				return nil, err
			}
			//if password != "" {
			//	if _, err := c.Do("AUTH", password); err != nil {
			//		c.Close()
			//		return nil, err
			//	}
			//}
			return c, err
		},
	}
}

func (m *M) Init() error {
	err := m.connectRedis()
	if nil != err {
		return err
	}

	return nil
}

//connect redis
func (m *M) connectRedis() error {
	if nil != m.redis {
		_, err := m.redis.Do("PING")
		if nil == err {
			return nil
		}
	}
	m.redis = RedisPool.Get()
	if nil == m.redis {
		return errors.New("Connect redis failed")
	}
	_, err := m.redis.Do("SELECT", m.redisDbIndex)
	logs.Logger.Debugf("SELECT redis db index=%d", m.redisDbIndex)
	if nil != err {
		logs.Logger.Errorf("SELECT RedisDbIndex:%s", err)
		return err
	}
	return nil
}

func (m *M) CloseAll() {
	m.closeRedis()
}

func (m *M) closeRedis() {
	if nil != m && nil != m.redis {
		err := m.redis.Close()
		if err != nil {
			logs.Logger.Errorf("close redis err %s.", err)
		}else{
			logs.Logger.Debug("close redis success.")
		}
	}
}
