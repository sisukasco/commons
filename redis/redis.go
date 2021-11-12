package redis

import (
	"log"
	"strconv"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/common"
	redigo "github.com/gomodule/redigo/redis"
)

type Redis struct {
	pool *redigo.Pool
	common.RedisConnector
}

func New(redisURL string) (*Redis, error) {

	host, passwd, db, err := machinery.ParseRedisURL(redisURL)
	if err != nil {
		log.Fatalf("Error parsing redis URL %v", err)
		return nil, err
	}
	r := &Redis{}

	r.pool = r.NewPool("", host, passwd, db, nil, nil)

	return r, nil
}

func (r *Redis) Open() redigo.Conn {
	return r.pool.Get()
}

func (r *Redis) RegisterExpiredForms(formIDs []string) error {
	redisConn := r.Open()
	defer redisConn.Close()

	for _, fid := range formIDs {
		_, err := redisConn.Do("SADD", "expired_forms", fid)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Redis) PopExpiredForm() string {
	redisConn := r.Open()
	defer redisConn.Close()
	res, err := redigo.String(redisConn.Do("SPOP", "expired_forms"))
	if err != nil {
		return ""
	}
	return res
}
func (r *Redis) Delete(key string) error {
	redisConn := r.Open()
	defer redisConn.Close()

	_, err := redisConn.Do("DEL", key)
	return err
}

func (r *Redis) Expire(key string, seconds int64) error {
	redisConn := r.Open()
	defer redisConn.Close()

	strSecs := strconv.FormatInt(seconds, 10)

	_, err := redisConn.Do("EXPIRE", key, strSecs)
	return err
}

func (r *Redis) Exists(key string) bool {
	redisConn := r.Open()
	defer redisConn.Close()

	exists, err := redigo.Bool(redisConn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

func (r *Redis) CreateList(list string, args []int) error {
	redisConn := r.Open()
	defer redisConn.Close()
	ll := redigo.Args{}.Add(list).AddFlat(args)

	_, err := redisConn.Do("RPUSH", ll...)
	return err
}

func (r *Redis) RoundRobinFromList(list string) (int, error) {
	redisConn := r.Open()
	defer redisConn.Close()
	v, err := redigo.Int(redisConn.Do("RPOPLPUSH", list, list))
	return v, err
}

func (r *Redis) CreateHMap(key string, obj map[string]string) error {
	redisConn := r.Open()
	defer redisConn.Close()

	var args = []interface{}{key}
	for k, v := range obj {
		args = append(args, k, v)
	}

	_, err := redisConn.Do("HMSET", args...)

	return err
}

func (r *Redis) GetHMapValue(key string, subKey string) (string, error) {
	redisConn := r.Open()
	defer redisConn.Close()

	value, err := redigo.String(redisConn.Do("HGET", key, subKey))

	if err != nil {
		return "", err
	}

	return value, nil

}

func (r *Redis) SetHMapValue(key string, subKey string, value string) error {
	redisConn := r.Open()
	defer redisConn.Close()

	_, err := redisConn.Do("HSET", key, subKey, value)

	if err != nil {
		return err
	}

	return nil

}

func (r *Redis) SetKeyValue(key string, value string, expirySeconds int64) error {
	redisConn := r.Open()
	defer redisConn.Close()

	strSecs := strconv.FormatInt(expirySeconds, 10)

	_, err := redisConn.Do("SET", key, value, "EX", strSecs)

	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetKeyValue(key string) (string, error) {
	redisConn := r.Open()
	defer redisConn.Close()

	value, err := redigo.String(redisConn.Do("GET", key))

	if err != nil {
		return "", err
	}

	return value, nil
}
