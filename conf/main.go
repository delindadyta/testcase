package main

import (
	. "eaciit/bankingsalesperformance/consoleapps/controllers"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	"os"
	"runtime"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	tk.Println("Starting the app..\n")

	conn, e := PrepareConnection()
	if e != nil {
		tk.Println(e)
	} else {
		base := new(BaseController)
		base.Ctx = orm.New(conn)
		defer base.Ctx.Close()
		// new(AccountPlanSummary).Generate(base)
		new(ProspectSummary).Generate(base)
		// new(LoginUserSummary).Generate(base)

	}

	tk.Println("Application Close..")
}
