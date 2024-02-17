package merchantdashboard

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"errors"
	"html/template"
	"runtime/debug"
	"strings"

	"github.com/astaxie/beego"
	"ominaya.com/database/sql"
	"ominaya.com/util/log"
)

type MenuData struct {
	Menus string
}

type MerchantDashboard struct {
	beego.Controller
}

func (c *MerchantDashboard) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.TplName = "error/error.html"
		}
		if err != nil && err.Error() == "Session Time Out.Please Logout and Login Again." {
			c.Data["DisplayMessage"] = err.Error()
			c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Fail")
		} else if err != nil {
			c.Data["DisplayMessage"] = err.Error()
			c.TplName = "error/error.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "merchnatdashboard/merchnatdashboard.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Merchnat View Proflie  Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)

	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "UserName - ", sess.Get("uname"))
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Role - ", sess.Get("role"))

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Session ID - ", sess.SessionID())
	if err = session.ValidateSession(sess); err != nil {
		sess.SessionRelease(c.Ctx.ResponseWriter)
		session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()

	c.Data["name"] = sess.Get("name").(string)
	c.Data["role"] = sess.Get("role").(string)
	c.Data["Uname"] = sess.Get("uname").(string)
	c.Data["language"] = sess.Get("language").(string)
	c.Data["department"] = sess.Get("department").(string)

	if sess.Get("role") == "MERCHANT" && sess.Get("language") == "English" {
		c.Data["Menus1"] = template.HTML(`` + beego.AppConfig.String("ENGLISH_USER_TEMPLATE") + ``)
		headerContent := strings.Replace(beego.AppConfig.String("MERCHANT_ENGLISH_HEADER_TEMPLATE"), "{{.Uname}}", sess.Get("uname").(string), -1)

		c.Data["Header1"] = template.HTML(`` + headerContent + ``)
		c.Data["Dashboard"] = beego.AppConfig.String("ENGLISH_DASHBOARD")
		c.Data["TotalNumberAdminUsers"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_ADMIN_USERS")
		c.Data["TotalNumberMerchantUsers"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_MERCHNAT_USERS")

		c.Data["TotalNumberOfTransactions"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_TRANSACTIONS")
		c.Data["TotalAmountofTransactions"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_TRANSACTIONS_AMOUNT")

		c.Data["TotalSuccessfulTransactions"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_SUCCESSFUL_TRANSACTIONS")
		c.Data["TotalAmountofSuccessfulTransactions"] = beego.AppConfig.String("ENGLISH_TOTALAMOUNT_OF_SUCCESSFUL_TRANSACTIONS")

		c.Data["TotalDeclinedTransactions"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_DECLINED_TRANSACTIONS")
		c.Data["TotalAmountofDeclinedTransactions"] = beego.AppConfig.String("ENGLISH_TOTALAMOUNT_OF_DECLINED_TRANSACTIONS")
		c.Data["entityDashboard"] = beego.AppConfig.String("ENGLISH_ENTITY_DASHBOARD")
		c.TplName = "merchnatdashboard/merchnatdashboard.html"
	} else if sess.Get("role") == "MERCHANT" && sess.Get("language") == "French" {
		c.Data["Menus1"] = template.HTML(`` + beego.AppConfig.String("FRENCH_USER_TEMPLATE") + ``)
		headerContent := strings.Replace(beego.AppConfig.String("MERCHANT_FRENCH_HEADER_TEMPLATE"), "{{.Uname}}", sess.Get("uname").(string), -1)

		c.Data["Header1"] = template.HTML(`` + headerContent + ``)

		c.Data["Dashboard"] = beego.AppConfig.String("FRENCH_DASHBOARD")
		c.Data["TotalNumberAdminUsers"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_ADMIN_USERS")
		c.Data["TotalNumberMerchantUsers"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_MERCHNAT_USERS")

		c.Data["TotalNumberOfTransactions"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_TRANSACTIONS")
		c.Data["TotalAmountofTransactions"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_TRANSACTIONS_AMOUNT")

		c.Data["TotalSuccessfulTransactions"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_SUCCESSFUL_TRANSACTIONS")
		c.Data["TotalAmountofSuccessfulTransactions"] = beego.AppConfig.String("FRENCH_TOTALAMOUNT_OF_SUCCESSFUL_TRANSACTIONS")

		c.Data["TotalDeclinedTransactions"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_DECLINED_TRANSACTIONS")
		c.Data["TotalAmountofDeclinedTransactions"] = beego.AppConfig.String("FRENCH_TOTALAMOUNT_OF_DECLINED_TRANSACTIONS")

		c.Data["UserDetails"] = beego.AppConfig.String("FRENCH_USER_DETAILS")
		c.Data["TrandDetails"] = beego.AppConfig.String("FRENCH_TRANS_DETAILS")
		c.Data["entityDashboard"] = beego.AppConfig.String("FRENCH_ENTITY_DASHBOARD")
		c.TplName = "merchnatdashboard/merchnatdashboard.html"
	}

	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "acount  ", acount)

	row2, err := db.Db.Query(`select count(JSON_EXTRACT(transaction_deatils, "$.Status")) as TCont,SUM(JSON_EXTRACT(transaction_deatils, "$.Amount")) as TSum,
(select count(JSON_EXTRACT(transaction_deatils, "$.Status")) FROM Transactions where (JSON_EXTRACT(transaction_deatils, "$.Status") = "APPROVED") AND (JSON_EXTRACT(transaction_deatils, "$.Entity_email") = ?)) as ACont,
(select SUM(JSON_EXTRACT(transaction_deatils, "$.Amount")) FROM Transactions where (JSON_EXTRACT(transaction_deatils, "$.Status") = "APPROVED") AND (JSON_EXTRACT(transaction_deatils, "$.Entity_email") = ?)) as ASum,
(select count(JSON_EXTRACT(transaction_deatils, "$.Status")) FROM Transactions where (JSON_EXTRACT(transaction_deatils, "$.Status") = "DECLINED") AND (JSON_EXTRACT(transaction_deatils, "$.Entity_email") = ?)) as DCont,
(select SUM(JSON_EXTRACT(transaction_deatils, "$.Amount")) FROM Transactions where (JSON_EXTRACT(transaction_deatils, "$.Status") = "DECLINED") AND (JSON_EXTRACT(transaction_deatils, "$.Entity_email") = ?)) as DSum
FROM Transactions where JSON_EXTRACT(transaction_deatils, "$.Entity_email") = ?`, sess.Get("uname").(string), sess.Get("uname").(string), sess.Get("uname").(string), sess.Get("uname").(string), sess.Get("uname").(string))

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row2)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row2)
	_, data2, err := sql.Scan(row2)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data2, "\nData len - ", len(data2))
	if len(data2) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data1 - ", data2[0][0])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data2 - ", data2[0][1])

	c.Data["totaltransactioncount1"] = data2[0][0]
	c.Data["totaltransactionamount1"] = data2[0][1]
	c.Data["totaltsuccessfulransactioncount1"] = data2[0][2]
	c.Data["totaltsuccessfulransactionamount1"] = data2[0][3]
	c.Data["toatldeclinedtransactioncount1"] = data2[0][4]
	c.Data["toatldeclinedtransactionamount1"] = data2[0][5]

	return
}
