package dashboard

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"errors"
	"html/template"
	"runtime/debug"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"ominaya.com/database/sql"
	"ominaya.com/util/log"
)

type MenuData struct {
	Menus string
}

type Dashboard struct {
	beego.Controller
}

func (c *Dashboard) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	var Autherr error
	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.TplName = "error/error.html"
		}
		if Autherr != nil {
			c.Data["DisplayMessage"] = Autherr.Error()
			c.TplName = "error/error.html"
			return
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
			c.TplName = "dashboard/dashboard.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Success")
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
	if err = session.ValidateSession(sess); err != nil {
		sess.SessionRelease(c.Ctx.ResponseWriter)
		session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "Dashboard")
	passSet := sess.Get("passwordSet").(string)
	if err != nil {
		beego.Error(err)
		Autherr = errors.New("Unable to get Menus")
		return
	}
	if !auth || passSet != "YES" {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "UserName - ", sess.Get("uname"))
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Role - ", sess.Get("role"))

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Session ID - ", sess.SessionID())
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()

	c.Data["name"] = sess.Get("name").(string)
	c.Data["role"] = sess.Get("role").(string)
	c.Data["Uname"] = sess.Get("uname").(string)
	c.Data["language"] = sess.Get("language").(string)
	c.Data["department"] = sess.Get("department").(string)

	c.Data["MenuJson"] = sess.Get("menujson")

	c.Data["operatorsCount"] = len(strings.Split(beego.AppConfig.String("CNPS_OPERATORS"), "|"))
	c.Data["Photo"] = sess.Get("photo").(string)
	c.Data["ChannelsCount"] = len(strings.Split(beego.AppConfig.String("CNPS_CHANNELS"), "|"))
	successAmount, servicecharge, totalTransCount, totalBanks, successCount, pendingCount, declainedCount, transErr := utils.SideBarTransactionData()
	if transErr != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get sidebar data")
		return
	}
	currentTime := time.Now()
	if sess.Get("role") == "ADMIN" && sess.Get("language") == "English" {
		menuContent := strings.Replace(beego.AppConfig.String("MENU_TEMPLATE"), "{{.SuccessAmount}}", beego.AppConfig.String("BaseCurrency")+" "+successAmount, -1)
		menuContent = strings.Replace(menuContent, "{{.TotalTransCount}}", totalTransCount, -1)
		menuContent = strings.Replace(menuContent, "{{.servicecharge}}", beego.AppConfig.String("BaseCurrency")+" "+servicecharge, -1)
		menuContent = strings.Replace(menuContent, "{{.BanksCount}}", totalBanks, -1)
		menuContent = strings.Replace(menuContent, "{{.SuccessCount}}", successCount, -1)
		menuContent = strings.Replace(menuContent, "{{.PendingCount}}", pendingCount, -1)
		menuContent = strings.Replace(menuContent, "{{.DeclainedCount}}", declainedCount, -1)
		c.Data["Menus"] = template.HTML(`` + menuContent + ``)
		headerContent := strings.Replace(beego.AppConfig.String("ENGLISH_HEADER_TEMPLATE"), "{{.Fullname}}", sess.Get("fullname").(string), -1)
		headerContent = strings.Replace(headerContent, "{{.Rolename}}", sess.Get("rolename").(string), -1)
		c.Data["Header"] = template.HTML(`` + headerContent + ``)
		c.Data["Dashboard"] = beego.AppConfig.String("ENGLISH_DASHBOARD")
		c.Data["TotalNumberAdminUsers"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_ADMIN_USERS")
		c.Data["TotalNumberMerchantUsers"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_MERCHNAT_USERS")

		c.Data["TotalNumberOfTransactions"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_TRANSACTIONS")
		c.Data["TotalAmountofTransactions"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_TRANSACTIONS_AMOUNT")
		c.Data["TotalSuccessfulTransactions"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_SUCCESSFUL_TRANSACTIONS")
		c.Data["TotalAmountofSuccessfulTransactions"] = beego.AppConfig.String("ENGLISH_TOTALAMOUNT_OF_SUCCESSFUL_TRANSACTIONS")
		c.Data["TotalDeclinedTransactions"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_DECLINED_TRANSACTIONS")
		c.Data["TotalAmountofDeclinedTransactions"] = beego.AppConfig.String("ENGLISH_TOTALAMOUNT_OF_DECLINED_TRANSACTIONS")
		c.Data["UserDetails"] = beego.AppConfig.String("ENGLISH_USER_DETAILS")
		c.Data["TrandDetails"] = beego.AppConfig.String("ENGLISH_TRANS_DETAILS")
		c.Data["Dashboard2"] = beego.AppConfig.String("ENGLISH_DASHBOARD")
		c.Data["TotalPendingTransactions"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_PENDING_TRANSACTIONS")
		c.Data["TotalAmountofPendingTransactions"] = beego.AppConfig.String("ENGLISH_TOTALAMOUNT_OF_PENDING_TRANSACTIONS")
		c.Data["TotalInitiatedTransactions"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_INITIATED_TRANSACTIONS")
		c.Data["TotalAmountofInitiatedTransactions"] = beego.AppConfig.String("ENGLISH_TOTALAMOUNT_OF_INITIATED_TRANSACTIONS")

		//New template words
		c.Data["UserTransactionDetails"] = strings.Replace(beego.AppConfig.String("ENGLISH_USER_TRANSACTION_DETAILS"), "{{.curdate}}", currentTime.Format("2006-Jan-02"), -1)
		c.Data["TransactionDetails"] = beego.AppConfig.String("ENGLISH_TRANSACTION_DETAILS")
		c.Data["TransactionAmountDetails"] = beego.AppConfig.String("ENGLISH_TRANSACTION_AMOUNT_DETAILS")
		c.Data["Aproved"] = beego.AppConfig.String("ENGLISH_APPROVED")
		c.Data["Pending"] = beego.AppConfig.String("ENGLISH_PENDING")
		c.Data["Declained"] = beego.AppConfig.String("ENGLISH_DECLINED")
		c.Data["TotalTransactions"] = beego.AppConfig.String("ENGLISH_TOTAL_TRANSACTIONS")
		c.Data["SuccessfulTransactions"] = beego.AppConfig.String("ENGLISH_SUCCESSFUL_TRANSACTIONS")
		c.Data["DeclainedTransactions"] = beego.AppConfig.String("ENGLISH_DECLAINED_TRANSACTIONS")
		c.Data["CnpsPortalAdminUsers"] = beego.AppConfig.String("ENGLISH_CNPS_PORTAL_ADMIN_USERS")
		c.Data["CnpsPortalPaymentOrganizations"] = beego.AppConfig.String("ENGLISH_CNPS_PORTAL_PAYMENT_ORGANIZATIONS")
		c.Data["TransactioByOperatorDetails"] = beego.AppConfig.String("ENGLISH_TRANSACTION_BY_OPERATOR")
		c.Data["TransactionByChannelDetails"] = beego.AppConfig.String("ENGLISH_TRANSACTION_BY_CHANNEL")
		c.Data["TransactionStatusDetails"] = beego.AppConfig.String("ENGLISH_TRANSACTION_STATUS_DETAILS")
		c.Data["TransactionCountDetails"] = beego.AppConfig.String("ENGLISH_TRANSACTION_COUNT_DETAILS")
		c.Data["Search"] = beego.AppConfig.String("ENGLISH_SEARCH")
		c.Data["Reset"] = beego.AppConfig.String("ENGLISH_RESET")
		c.Data["DateRange"] = beego.AppConfig.String("ENGLISH_USER_SELECT_DATARANGE")
		c.Data["Apply"] = beego.AppConfig.String("ENGLISH_APPLY")
		c.Data["PleaseSelect"] = beego.AppConfig.String("ENGLISH_USER_SELECT_ROLE")
		c.Data["All"] = beego.AppConfig.String("ENGLISH_ALL")
		c.Data["SelectStatus"] = beego.AppConfig.String("ENGLISH_SELECT_STATUS")
		c.Data["TotalNumberChannels"] = beego.AppConfig.String("ENGLISH_TOTAL_NUMBER_OF_CHANNELS")
		c.Data["CnpsPortalPaymentChannels"] = beego.AppConfig.String("ENGLISH_CNPS_PORTAL_PAYMENT_CHANNELS")
		c.Data["Msg"] = beego.AppConfig.String("ENGLISH_MSG")
		c.Data["From"] = beego.AppConfig.String("FROM")
		c.Data["To"] = beego.AppConfig.String("TO")

		//sess.Set("Menus", beego.AppConfig.String("MENU_TEMPLATE"))
		//sess.Set("Header", headerContent)
		c.TplName = "dashboard/dashboard.html"
	} else if sess.Get("role") == "ADMIN" && sess.Get("language") == "French" {
		menuContent := strings.Replace(beego.AppConfig.String("FRENCH_MENU_TEMPLATE"), "{{.SuccessAmount}}", beego.AppConfig.String("BaseCurrency")+" "+successAmount, -1)
		menuContent = strings.Replace(menuContent, "{{.TotalTransCount}}", totalTransCount, -1)
		menuContent = strings.Replace(menuContent, "{{.servicecharge}}", beego.AppConfig.String("BaseCurrency")+" "+servicecharge, -1)
		menuContent = strings.Replace(menuContent, "{{.BanksCount}}", totalBanks, -1)
		menuContent = strings.Replace(menuContent, "{{.SuccessCount}}", successCount, -1)
		menuContent = strings.Replace(menuContent, "{{.PendingCount}}", pendingCount, -1)
		menuContent = strings.Replace(menuContent, "{{.DeclainedCount}}", declainedCount, -1)
		c.Data["Menus"] = template.HTML(`` + menuContent + ``)
		headerContent := strings.Replace(beego.AppConfig.String("FRENCH_HEADER_TEMPLATE"), "{{.Fullname}}", sess.Get("fullname").(string), -1)
		headerContent = strings.Replace(headerContent, "{{.Rolename}}", sess.Get("rolename").(string), -1)
		c.Data["Header"] = template.HTML(`` + headerContent + ``)

		c.Data["Dashboard"] = beego.AppConfig.String("FRENCH_DASHBOARD")
		c.Data["TotalNumberAdminUsers"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_ADMIN_USERS")
		c.Data["TotalNumberMerchantUsers"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_MERCHNAT_USERS")
		c.Data["TotalNumberOfTransactions"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_TRANSACTIONS")
		c.Data["TotalAmountofTransactions"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_TRANSACTIONS_AMOUNT")
		c.Data["TotalSuccessfulTransactions"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_SUCCESSFUL_TRANSACTIONS")
		c.Data["TotalAmountofSuccessfulTransactions"] = beego.AppConfig.String("FRENCH_TOTALAMOUNT_OF_SUCCESSFUL_TRANSACTIONS")
		c.Data["TrandDetails"] = beego.AppConfig.String("ENGLISH_TRANS_DETAILS")
		c.Data["TotalDeclinedTransactions"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_DECLINED_TRANSACTIONS")
		c.Data["TotalAmountofDeclinedTransactions"] = beego.AppConfig.String("FRENCH_TOTALAMOUNT_OF_DECLINED_TRANSACTIONS")
		c.Data["UserDetails"] = beego.AppConfig.String("FRENCH_USER_DETAILS")
		c.Data["TrandDetails"] = beego.AppConfig.String("FRENCH_TRANS_DETAILS")
		sess.Set("Menus", beego.AppConfig.String("FRENCH_MENU_TEMPLATE"))
		c.Data["Dashboard2"] = beego.AppConfig.String("FRENCH_DASHBOARD")
		c.Data["TotalPendingTransactions"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_PENDING_TRANSACTIONS")
		c.Data["TotalAmountofPendingTransactions"] = beego.AppConfig.String("FRENCH_TOTALAMOUNT_OF_PENDING_TRANSACTIONS")
		c.Data["TotalInitiatedTransactions"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_INITIATED_TRANSACTIONS")
		c.Data["TotalAmountofInitiatedTransactions"] = beego.AppConfig.String("FRENCH_TOTALAMOUNT_OF_INITIATED_TRANSACTIONS")

		//New Template words
		c.Data["UserTransactionDetails"] = strings.Replace(beego.AppConfig.String("FRENCH_USER_TRANSACTION_DETAILS"), "{{.curdate}}", currentTime.Format("2006-Jan-02"), -1)
		c.Data["TransactionDetails"] = beego.AppConfig.String("FRENCH_TRANSACTION_DETAILS")
		c.Data["TransactionAmountDetails"] = beego.AppConfig.String("FRENCH_TRANSACTION_AMOUNT_DETAILS")
		c.Data["Aproved"] = beego.AppConfig.String("FRENCH_APPROVED")
		c.Data["Pending"] = beego.AppConfig.String("FRENCH_PENDING")
		c.Data["Declained"] = beego.AppConfig.String("FRENCH_DECLINED")
		c.Data["TotalTransactions"] = beego.AppConfig.String("FRENCH_TOTAL_TRANSACTIONS")
		c.Data["SuccessfulTransactions"] = beego.AppConfig.String("FRENCH_SUCCESSFUL_TRANSACTIONS")
		c.Data["DeclainedTransactions"] = beego.AppConfig.String("FRENCH_DECLAINED_TRANSACTIONS")
		c.Data["CnpsPortalAdminUsers"] = beego.AppConfig.String("FRENCH_CNPS_PORTAL_ADMIN_USERS")
		c.Data["CnpsPortalPaymentOrganizations"] = beego.AppConfig.String("FRENCH_CNPS_PORTAL_PAYMENT_ORGANIZATIONS")
		c.Data["TransactioByOperatorDetails"] = beego.AppConfig.String("FRENCH_TRANSACTION_BY_OPERATOR")
		c.Data["TransactionByChannelDetails"] = beego.AppConfig.String("FRENCH_TRANSACTION_BY_CHANNEL")
		c.Data["TransactionStatusDetails"] = beego.AppConfig.String("FRENCH_TRANSACTION_STATUS_DETAILS")
		c.Data["TransactionCountDetails"] = beego.AppConfig.String("FRENCH_TRANSACTION_COUNT_DETAILS")
		c.Data["Search"] = beego.AppConfig.String("FRENCH_SEARCH")
		c.Data["Reset"] = beego.AppConfig.String("FRENCH_RESET")
		c.Data["DateRange"] = beego.AppConfig.String("FRENCH_USER_SELECT_DATARANGE")
		c.Data["Apply"] = beego.AppConfig.String("FRENCH_APPLY")
		c.Data["PleaseSelect"] = beego.AppConfig.String("FRENCH_USER_SELECT_ROLE")
		c.Data["All"] = beego.AppConfig.String("FRENCH_ALL")
		c.Data["SelectStatus"] = beego.AppConfig.String("FRENCH_SELECT_STATUS")
		c.Data["TotalNumberChannels"] = beego.AppConfig.String("FRENCH_TOTAL_NUMBER_OF_CHANNELS")
		c.Data["CnpsPortalPaymentChannels"] = beego.AppConfig.String("FRENCH_CNPS_PORTAL_PAYMENT_CHANNELS")
		c.Data["Msg"] = beego.AppConfig.String("FRENCH_MSG")
		c.Data["From"] = beego.AppConfig.String("FROM")
		c.Data["To"] = beego.AppConfig.String("TO")

		//sess.Set("Header", headerContent)

		c.TplName = "dashboard/dashboard.html"
	}
	row, err := db.Db.Query(`SELECT count(*)as acount,(SELECT count(*) FROM Users where role = 'MERCHANT') as mcount FROM Users where role = 'ADMIN'`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data1 - ", data[0][0])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data2 - ", data[0][1])

	c.Data["admincount"] = data[0][0]
	c.Data["merchantcount"] = data[0][1]

	// row1, err := db.Db.Query(`select count(JSON_EXTRACT(transaction_deatils, "$.Status")) as TCont,SUM(JSON_EXTRACT(transaction_deatils, "$.Amount")) as TSum,
	// (select count(JSON_EXTRACT(transaction_deatils, "$.Status")) FROM Transactions where JSON_EXTRACT(transaction_deatils, "$.Status") = "APPROVED") as ACont,
	// (select SUM(JSON_EXTRACT(transaction_deatils, "$.Amount")) FROM Transactions where JSON_EXTRACT(transaction_deatils, "$.Status") = "APPROVED") as ASum,
	// (select count(JSON_EXTRACT(transaction_deatils, "$.Status")) FROM Transactions where JSON_EXTRACT(transaction_deatils, "$.Status") = "DECLINED") as DCont,
	// (select SUM(JSON_EXTRACT(transaction_deatils, "$.Amount")) FROM Transactions where JSON_EXTRACT(transaction_deatils, "$.Status") = "DECLINED") as DSum
	// FROM Transactions`)

	row1, err := db.Db.Query(`select (select count(status) FROM Transactions where status = "APPROVED" OR status="PENDING" OR status="DECLINED") as TCont,sum(amount) as TSum,
	(select count(status) FROM Transactions where status = "APPROVED") as ACont,
	(select sum(amount) FROM Transactions where status = "APPROVED") as ASum,
	(select count(status) FROM Transactions where status = "DECLINED") as DCont,
	(select sum(amount) FROM Transactions where status = "DECLINED") as DSum,
	(select count(status) FROM Transactions where status = "PENDING") as PCont,
	(select sum(amount) FROM Transactions where status = "PENDING") as PSum,
	(select count(status) FROM Transactions where status = "INITIATED") as ICont,
	(select sum(amount) FROM Transactions where status = "INITIATED") as ISum
 	from Transactions`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row1)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row1)
	_, data1, err := sql.Scan(row1)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data1, "\nData len - ", len(data1))
	if len(data1) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data1 - ", data1[0][0])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data2 - ", data1[0][1])

	c.Data["totaltransactioncount"] = data1[0][0]
	c.Data["totaltransactionamount"] = data1[0][1]
	c.Data["totaltsuccessfulransactioncount"] = data1[0][2]
	c.Data["totaltsuccessfulransactionamount"] = data1[0][3]
	c.Data["toatldeclinedtransactioncount"] = data1[0][4]
	c.Data["toatldeclinedtransactionamount"] = data1[0][5]
	c.Data["toatlpendingtransactioncount"] = data1[0][6]
	c.Data["toatlpendingtransactionamount"] = data1[0][7]
	c.Data["toatlinitiatedtransactioncount"] = data1[0][8]
	c.Data["toatlinitiatedtransactionamount"] = data1[0][9]

	return
}
