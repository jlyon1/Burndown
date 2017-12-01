package database

import (
)

type DB interface{
  Connect() bool
  Disconnect() bool
  Find(key string) string
  Set(key string, val interface{}) bool
  Delete(key string) bool
  Expire(key string, val int64) bool

}
