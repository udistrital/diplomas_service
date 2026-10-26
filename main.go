package main

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/udistrital/diplomas_mid/routers"
	"github.com/udistrital/diplomas_mid/services"
)

func main() {
	if err := services.LoadRuntimeConfig(); err != nil {
		logs.Error("error cargando configuracion: %v", err)
		return
	}

	beego.Run()
}
