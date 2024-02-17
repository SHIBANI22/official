package main

import (
	"CNPSAdmin/model/utils"
	_ "CNPSAdmin/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var RedirectHttp = func(ctx *context.Context) {
	utils.SetHTTPHeader(ctx)
	if !ctx.Input.IsSecure() {
		// no need for an additional '/' between domain and uri
		url := "https://" + ctx.Input.Domain() + ":" + beego.AppConfig.String("HttpsPort") + ctx.Input.URI()
		ctx.Redirect(302, url)
	}
}

func main() {
	if beego.AppConfig.String("EnableHTTPS") == "true" {
		beego.InsertFilter("/", beego.BeforeRouter, RedirectHttp) // for http://mysite
		beego.InsertFilter("*", beego.BeforeRouter, RedirectHttp) // for http://mysite/*
	}
	beego.Run()
}
