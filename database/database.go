package database

import (
  "time"
)

type DB interface{
  Connect() bool
  Disconnect() bool
  Find(key string) string
  Set(key string, val interface{}) bool
  Delete(key string) bool
  Expire(key string, val time.Duration) bool
  TTL(key string) time.Duration


}
