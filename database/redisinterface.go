package database

import (
  "gopkg.in/redis.v5"
  "encoding/json"
  "time"
  //"fmt"
)

type Redis struct{
  IP string
  Port string
  Password string
  DB int
  Client *redis.Client
}

func (r *Redis) Connect() (bool){
  r.Client = redis.NewClient(&redis.Options{
    Addr: r.IP + ":" +  r.Port,
    Password: r.Password,
    DB: r.DB,
  })
  _, err := r.Client.Ping().Result()
  if(err != nil){
    return false
  }else{
    return true
  }
}
func (r *Redis) Disconnect() bool{

  err := r.Client.Close()
  if(err != nil){
    return false
  }
  return true
}

func (r *Redis) Find(key string) string{
  val, _ := r.Client.Get(key).Result()
  //fmt.Printf("val " + val)
  return val
}

func (r *Redis) Expire(key string, val int64) bool{
  ret, _ := r.Client.Expire(key,time.Duration(val)*time.Second).Result()
  //fmt.Printf("val " + val)
  return ret
}

func (r *Redis) Set(key string, value interface{}) bool{
  b,_ := json.Marshal(value)
  r.Client.Set(key,b,0)
  //fmt.Printf("%v",b)
  return true
}

func (r *Redis) Delete(key string) bool{
  return false
}
