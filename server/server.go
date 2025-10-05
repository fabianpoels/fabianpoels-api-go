package server

import (
	"fmt"

	"github.com/fabianpoels/fabianpoels-api-go/cache"
	"github.com/fabianpoels/fabianpoels-api-go/config"
	"github.com/fabianpoels/fabianpoels-api-go/db"
)

func Init() {
	config := config.GetConfig()
	r := NewRouter()
	db.DbConnect()
	cache.CacheConnect()
	r.Run(fmt.Sprintf("%s:%s", config.GetString("server.host"), config.GetString("server.port")))
}
