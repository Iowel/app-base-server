package cache

import "github.com/Iowel/app-base-server/internal/user"

type IPostCache interface {
	Set(key string, value *user.UserCache)
	Get(key string) *user.UserCache
	GetAll() []*user.UserCache
	Delete(key string)
}
