package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/fafacms/core/util"
	"github.com/hunterhug/fafacms/core/util/kv"
	"strconv"
	"strings"
)

var (
	redisToken = "ff_tokens"
	redisUser  = "ff_users"
)

// diy user redis
type TokenManage interface {
	CheckToken(token string) (user *model.User, err error)                 // 检查令牌是否存在, 返回redis用户，redis用户不存在缓存击穿到mysql，mysql不存在删除该用户所有令牌
	SetToken(user *model.User, validTimes int64) (token string, err error) // 设置令牌，登录或者激活用户的时候，有效期7天，每次都是强覆盖
	RefreshToken(token string) error                                       // 刷新令牌，每次浏览器启动时，自己保持的cookie请求延长令牌时间
	DeleteToken(token string) error                                        // 删除令牌，在退出登录的时候
	RefreshUser(id []int) error                                            // 刷新用户缓存信息
	DeleteUserToken(id int) error                                          // 删除一个用户下面所有的临时令牌，在用户修改密码的情况
	DeleteUser(id int) error                                               // 删除缓存中的用户信息，当用户被删除的时候应该删除
	AddUser(id int) (user *model.User, err error)                          // 增加缓存mysql用户到redis，有效期1天
}

type RedisSession struct {
	Pool *redis.Pool
}

func (s *RedisSession) Set(key string, value []byte, expireSecond int64) (err error) {
	conn := s.Pool.Get()
	if conn.Err() != nil {
		err = conn.Err()
		return
	}

	defer conn.Close()
	err = conn.Send("MULTI")
	if err != nil {
		return err
	}
	err = conn.Send("SET", key, value)
	if err != nil {
		return err
	}
	err = conn.Send("EXPIRE", key, expireSecond)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXEC")
	return
}

func (s *RedisSession) Delete(key string) (err error) {
	conn := s.Pool.Get()
	if conn.Err() != nil {
		err = conn.Err()
		return
	}

	defer conn.Close()
	_, err = conn.Do("DEL", key)
	return err
}

func (s *RedisSession) EXPIRE(key string, expireSecond int) (err error) {
	conn := s.Pool.Get()
	if conn.Err() != nil {
		err = conn.Err()
		return
	}

	defer conn.Close()
	_, err = conn.Do("EXPIRE", key, expireSecond)
	return err
}

func (s *RedisSession) Keys(pattern string) (result []string, exist bool, err error) {
	conn := s.Pool.Get()
	if conn.Err() != nil {
		err = conn.Err()
		return
	}

	defer conn.Close()
	keys, err := redis.ByteSlices(conn.Do("KEYS", pattern))
	if err == redis.ErrNil {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	result = make([]string, len(keys))
	for k, v := range keys {
		result[k] = string(v)
	}
	return result, true, nil
}

func (s *RedisSession) Get(key string) (value []byte, exist bool, err error) {
	conn := s.Pool.Get()
	if conn.Err() != nil {
		err = conn.Err()
		return
	}

	value, err = redis.Bytes(conn.Do("GET", key))
	if err == redis.ErrNil {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return value, true, nil
}

func HashTokenKey(token string) string {
	return fmt.Sprintf("%s_%s", redisToken, token)
}

func GenToken(id int) string {
	return fmt.Sprintf("%d_%s", id, util.GetGUID())
}

func HashUserKey(id int, name string) string {
	return fmt.Sprintf("%s_%d_%s", redisUser, id, name)
}

func UserKeys(id int) string {
	return fmt.Sprintf("%s_%d_*", redisUser, id)
}

func UserTokenKeys(id int) string {
	return fmt.Sprintf("%s_%d_*", redisToken, id)
}
func (s *RedisSession) CheckToken(token string) (user *model.User, err error) {
	if token == "" {
		err = errors.New("token nil")
		return
	}

	value, exist, err := s.Get(HashTokenKey(token))
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("token not found")
	}

	userKey := string(value)
	value, exist, err = s.Get(userKey)
	if err != nil {
		return nil, err
	}

	if exist {
		user = new(model.User)
		json.Unmarshal(value, user)
		return
	}

	temp := strings.Split(userKey, "_")
	if len(temp) != 3 || temp[0] != redisUser {
		return nil, errors.New("token invalid")
	}

	id, err := strconv.Atoi(temp[1])
	if err != nil {
		return nil, errors.New("token invalid")
	}
	user, err = s.AddUser(id)
	return
}

func (s *RedisSession) RefreshToken(token string) (err error) {
	return s.EXPIRE(HashTokenKey(token), 24*3600*7)
}

func (s *RedisSession) DeleteToken(token string) (err error) {
	return s.Delete(HashTokenKey(token))
}

func (s *RedisSession) DeleteUserToken(id int) (err error) {
	result, exist, err := s.Keys(UserTokenKeys(id))
	if err == nil && exist {
		for _, v := range result {
			s.Delete(v)
		}
	}
	return
}

func (s *RedisSession) DeleteUser(id int) (err error) {
	result, exist, err := s.Keys(UserKeys(id))
	if err == nil && exist {
		for _, v := range result {
			return s.Delete(v)
		}
	}
	return
}

func (s *RedisSession) AddUser(id int) (user *model.User, err error) {
	user = new(model.User)
	user.Id = id
	exist, err := user.GetRaw()
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("user not exist in db")
	}

	user.Password = ""
	user.ActivateCode = ""
	user.ActivateCodeExpired = 0
	user.ResetCode = ""
	user.ResetCodeExpired = 0
	userKey := HashUserKey(user.Id, user.Name)
	raw, _ := json.Marshal(user)
	err = s.Set(userKey, raw, 48*3600)
	if err != nil {
		return nil, err
	}

	return
}

func (s *RedisSession) RefreshUser(ids []int) (err error) {
	for _, id := range ids {
		s.AddUser(id)
	}
	return
}

func (s *RedisSession) SetToken(user *model.User, validTimes int64) (token string, err error) {
	if user == nil || user.Id == 0 {
		err = errors.New("user nil")
		return
	}

	user.Password = ""
	user.ActivateCode = ""
	user.ActivateCodeExpired = 0
	user.ResetCode = ""
	user.ResetCodeExpired = 0

	token = GenToken(user.Id)
	userKey := HashUserKey(user.Id, user.Name)
	err = s.Set(HashTokenKey(token), []byte(userKey), validTimes)
	if err != nil {
		return
	}

	raw, _ := json.Marshal(user)
	s.Set(userKey, raw, 24*3600*7)
	return
}

var FafaSessionMgr TokenManage

func InitSession(redisConf kv.MyRedisConf) error {
	pool, err := kv.NewRedis(&redisConf)
	if err != nil {
		return err
	}
	FafaSessionMgr = &RedisSession{Pool: pool}
	return nil
}
