package main

import (
  "github.com/stanisdev/core"
) 

func main() {
  router := core.Router{Handlers: make(map[string]map[string]core.RouterHandler)}
  router.Config = core.GetConfig()
  core.DatabaseMigrate(router.Config)
  
  router.GET("/", core.Index)
  router.GET("/login", core.Login)
  router.POST("/login", core.LoginPost)
  router.GET("/logout", core.Logout)
  router.Start()
}
