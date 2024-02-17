package merchantTransactionReport

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"

	"strings"

	"html/template"
	"runtime/debug"

	"errors"

	"github.com/astaxie/beego"

	"ominaya.com/database/sql"
	"ominaya.com/util/log"
)

type Row struct {
	CNPSTxnDate           string
	TimeStamp             string
	Merchantid            string
	MerchantName          string
	Status                string
	Amount                string
	BankName              string
	ChannelName           string
	CNPSTransactionNumber string
	PGTransactionNumber   string
	BankTransactionNumber string
	BankTransactionDate   string
}

type MerchantTransactionReport struct {
	beego.Controller
}

func (c *MerchantTransactionReport) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Start")
	var err error
	sessErr := false
	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.TplName = "error/error.html"
		}
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)

			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "reports/merchanttransactionReport/merchanttransactionReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/merchanttransactionReport/merchanttransactionReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Success")
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
		sessErr = true

		return
	}
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()

	c.Data["language"] = sess.Get("language").(string)
	if sess.Get("role") == "MERCHANT" && sess.Get("language") == "English" {
		c.Data["Menus1"] = template.HTML(`` + beego.AppConfig.String("ENGLISH_USER_TEMPLATE") + ``)
		headerContent := strings.Replace(beego.AppConfig.String("MERCHANT_ENGLISH_HEADER_TEMPLATE"), "{{.Uname}}", sess.Get("uname").(string), -1)

		c.Data["Header1"] = template.HTML(`` + headerContent + ``)
		c.Data["EntityTxnrep"] = beego.AppConfig.String("ENGLISH_ENTITY_TXNREPORT")

		c.Data["Dashboard"] = beego.AppConfig.String("ENGLISH_DASHBOARD")
		c.Data["TransactionReport"] = beego.AppConfig.String("ENGLISH_TRANSACTION_REPORT")
		c.Data["SearchFilters"] = beego.AppConfig.String("ENGLISH_SEARCH_FILTERS")

		c.Data["SearchResults"] = beego.AppConfig.String("ENGLISH_USER_SEARCH_RESULTS")
		c.Data["SelectDateRange"] = beego.AppConfig.String("ENGLISH_USER_SELECT_DATARANGE")
		c.Data["CNPSTxnDate"] = beego.AppConfig.String("ENGLISH_CNPS_TXN_DATE")
		c.Data["PGTxnDate"] = beego.AppConfig.String("ENGLISH_PG_TXN_DATE")
		c.Data["BankTxnDate"] = beego.AppConfig.String("ENGLISH_BANK_TXN_DATE")

		c.Data["TimeStamp"] = beego.AppConfig.String("ENGLISH_TIMESTAMP")
		c.Data["MerchantID"] = beego.AppConfig.String("ENGLISH_MERCHANT_ID")
		c.Data["CNPSTransactionNumber"] = beego.AppConfig.String("ENGLISH_CNPS_TRANSACTION_NUMBER")
		c.Data["PGTransactionNumber"] = beego.AppConfig.String("ENGLISH_PG_TRANSACTION_NUMBER")
		c.Data["BANKTransactionNumber"] = beego.AppConfig.String("ENGLISH_BANK_TRANSACTION_NUMBER")
		c.Data["MerchantName"] = beego.AppConfig.String("ENGLISH_MERCHNAT_NAME")
		c.Data["Amount"] = beego.AppConfig.String("ENGLISH_AMOUNT")
		c.Data["Bank"] = beego.AppConfig.String("ENGLISH_BANK")
		c.Data["Channel"] = beego.AppConfig.String("ENGLISH_CHANNEL")

		c.Data["Status"] = beego.AppConfig.String("ENGLISH_STATUS")
		c.Data["ListOfTransactionReports"] = beego.AppConfig.String("ENGLISH_LIST_OF_TRANSACTION_REPORT")
		c.Data["Search"] = beego.AppConfig.String("ENGLISH_SEARCH")
		c.Data["Reset"] = beego.AppConfig.String("ENGLISH_RESET")
		c.Data["please_select"] = beego.AppConfig.String("ENGLISH_PLEASESELECTCOMMON")
		c.Data["Approved"] = beego.AppConfig.String("ENGLISH_APPROVED")
		c.Data["Declined"] = beego.AppConfig.String("ENGLISH_DECLINED")
		c.Data["Pending"] = beego.AppConfig.String("ENGLISH_PENDING")
		c.Data["Reports"] = beego.AppConfig.String("ENGLISH_REPORTS")
		c.Data["input_cnpstxnnumnber"] = beego.AppConfig.String("ENGLISH_ENTER_CNPS_TXN_ID")
		c.Data["input_pgtxnnumnber"] = beego.AppConfig.String("ENGLISH_ENTER_PG_TXN_ID")
		c.Data["input_entityid"] = beego.AppConfig.String("ENGLISH_ENTER_ENTITY_ID")
		c.Data["input_entityname"] = beego.AppConfig.String("ENGLISH_ENTER_ENTITY_NAME")
		c.Data["input_amount"] = beego.AppConfig.String("ENGLISH_ENTER_AMOUNT")

		c.TplName = "reports/merchanttransactionReport/merchanttransactionReport.html"
	} else if sess.Get("role") == "MERCHANT" && sess.Get("language") == "French" {
		c.Data["Menus1"] = template.HTML(`` + beego.AppConfig.String("FRENCH_USER_TEMPLATE") + ``)
		headerContent := strings.Replace(beego.AppConfig.String("MERCHANT_FRENCH_HEADER_TEMPLATE"), "{{.Uname}}", sess.Get("uname").(string), -1)

		c.Data["Header1"] = template.HTML(`` + headerContent + ``)

		c.Data["Dashboard"] = beego.AppConfig.String("FRENCH_DASHBOARD")
		c.Data["EntityTxnrep"] = beego.AppConfig.String("FRENCH_ENTITY_TXNREPORT")
		c.Data["TransactionReport"] = beego.AppConfig.String("FRENCH_TRANSACTION_REPORT")
		c.Data["SearchFilters"] = beego.AppConfig.String("FRENCH_SEARCH_FILTERS")

		c.Data["SearchResults"] = beego.AppConfig.String("FRENCH_USER_SEARCH_RESULTS")
		c.Data["SelectDateRange"] = beego.AppConfig.String("FRENCH_USER_SELECT_DATARANGE")

		c.Data["CNPSTxnDate"] = beego.AppConfig.String("FRENCH_CNPS_TXN_DATE")
		c.Data["PGTxnDate"] = beego.AppConfig.String("FRENCH_PG_TXN_DATE")
		c.Data["MerchantID"] = beego.AppConfig.String("FRENCH_MERCHANT_ID")
		c.Data["BankTxnDate"] = beego.AppConfig.String("FRENCH_BANK_TXN_DATE")

		c.Data["Status"] = beego.AppConfig.String("FRENCH_STATUS")
		c.Data["CNPSTransactionNumber"] = beego.AppConfig.String("FRENCH_CNPS_TRANSACTION_NUMBER")
		c.Data["PGTransactionNumber"] = beego.AppConfig.String("FRENCH_PG_TRANSACTION_NUMBER")
		c.Data["BANKTransactionNumber"] = beego.AppConfig.String("FRENCH_BANK_TRANSACTION_NUMBER")
		c.Data["MerchantName"] = beego.AppConfig.String("FRENCH_MERCHANT_NAME")
		c.Data["Amount"] = beego.AppConfig.String("FRENCH_AMOUNT")
		c.Data["Bank"] = beego.AppConfig.String("FRENCH_BANK")
		c.Data["Channel"] = beego.AppConfig.String("FRENCH_CHANNEL")
		c.Data["ListOfTransactionReports"] = beego.AppConfig.String("FRENCH_LIST_OF_TRANSACTION_REPORT")
		c.Data["Search"] = beego.AppConfig.String("FRENCH_SEARCH")
		c.Data["Reset"] = beego.AppConfig.String("FRENCH_RESET")
		c.Data["please_select"] = beego.AppConfig.String("FRENCH_PLEASESELECTCOMMON")
		c.Data["Approved"] = beego.AppConfig.String("FRENCH_APPROVED")
		c.Data["Declined"] = beego.AppConfig.String("FRENCH_DECLINED")
		c.Data["Pending"] = beego.AppConfig.String("FRENCH_PENDING")
		c.Data["Reports"] = beego.AppConfig.String("FRENCH_REPORTS")
		c.Data["input_cnpstxnnumnber"] = beego.AppConfig.String("FRENCH_ENTER_CNPS_TXN_ID")
		c.Data["input_pgtxnnumnber"] = beego.AppConfig.String("FRENCH_ENTER_PG_TXN_ID")
		c.Data["input_entityid"] = beego.AppConfig.String("FRENCH_ENTER_ENTITY_ID")
		c.Data["input_entityname"] = beego.AppConfig.String("FRENCH_ENTER_ENTITY_NAME")
		c.Data["input_amount"] = beego.AppConfig.String("FRENCH_ENTER_AMOUNT")

		c.TplName = "reports/merchanttransactionReport/merchanttransactionReport.html"
	}

	return
}

func (c *MerchantTransactionReport) Post() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	sessErr := false
	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.TplName = "error/error.html"
		}
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)
			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "reports/merchanttransactionReport/merchanttransactionReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/merchanttransactionReport/merchanttransactionReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder  Page Success")
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
		sessErr = true

		return
	}
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()

	c.Data["language"] = sess.Get("language").(string)
	if sess.Get("role") == "MERCHANT" && sess.Get("language") == "English" {
		c.Data["Menus1"] = template.HTML(`` + beego.AppConfig.String("ENGLISH_USER_TEMPLATE") + ``)
		headerContent := strings.Replace(beego.AppConfig.String("MERCHANT_ENGLISH_HEADER_TEMPLATE"), "{{.Uname}}", sess.Get("uname").(string), -1)

		c.Data["Header1"] = template.HTML(`` + headerContent + ``)
		c.Data["EntityTxnrep"] = beego.AppConfig.String("ENGLISH_ENTITY_TXNREPORT")

		c.Data["Dashboard"] = beego.AppConfig.String("ENGLISH_DASHBOARD")
		c.Data["TransactionReport"] = beego.AppConfig.String("ENGLISH_TRANSACTION_REPORT")
		c.Data["SearchFilters"] = beego.AppConfig.String("ENGLISH_SEARCH_FILTERS")

		c.Data["SearchResults"] = beego.AppConfig.String("ENGLISH_USER_SEARCH_RESULTS")
		c.Data["SelectDateRange"] = beego.AppConfig.String("ENGLISH_USER_SELECT_DATARANGE")
		c.Data["CNPSTxnDate"] = beego.AppConfig.String("ENGLISH_CNPS_TXN_DATE")
		c.Data["PGTxnDate"] = beego.AppConfig.String("ENGLISH_PG_TXN_DATE")
		c.Data["BankTxnDate"] = beego.AppConfig.String("ENGLISH_BANK_TXN_DATE")

		c.Data["TimeStamp"] = beego.AppConfig.String("ENGLISH_TIMESTAMP")
		c.Data["MerchantID"] = beego.AppConfig.String("ENGLISH_MERCHANT_ID")
		c.Data["CNPSTransactionNumber"] = beego.AppConfig.String("ENGLISH_CNPS_TRANSACTION_NUMBER")
		c.Data["PGTransactionNumber"] = beego.AppConfig.String("ENGLISH_PG_TRANSACTION_NUMBER")
		c.Data["BANKTransactionNumber"] = beego.AppConfig.String("ENGLISH_BANK_TRANSACTION_NUMBER")
		c.Data["MerchantName"] = beego.AppConfig.String("ENGLISH_MERCHNAT_NAME")
		c.Data["Amount"] = beego.AppConfig.String("ENGLISH_AMOUNT")
		c.Data["Bank"] = beego.AppConfig.String("ENGLISH_BANK")
		c.Data["Channel"] = beego.AppConfig.String("ENGLISH_CHANNEL")

		c.Data["Status"] = beego.AppConfig.String("ENGLISH_STATUS")
		c.Data["ListOfTransactionReports"] = beego.AppConfig.String("ENGLISH_LIST_OF_TRANSACTION_REPORT")
		c.Data["Search"] = beego.AppConfig.String("ENGLISH_SEARCH")
		c.Data["Reset"] = beego.AppConfig.String("ENGLISH_RESET")
		c.Data["please_select"] = beego.AppConfig.String("ENGLISH_PLEASESELECTCOMMON")
		c.Data["Approved"] = beego.AppConfig.String("ENGLISH_APPROVED")
		c.Data["Declined"] = beego.AppConfig.String("ENGLISH_DECLINED")
		c.Data["Pending"] = beego.AppConfig.String("ENGLISH_PENDING")
		c.Data["Reports"] = beego.AppConfig.String("ENGLISH_REPORTS")
		c.Data["input_cnpstxnnumnber"] = beego.AppConfig.String("ENGLISH_ENTER_CNPS_TXN_ID")
		c.Data["input_pgtxnnumnber"] = beego.AppConfig.String("ENGLISH_ENTER_PG_TXN_ID")
		c.Data["input_entityid"] = beego.AppConfig.String("ENGLISH_ENTER_ENTITY_ID")
		c.Data["input_entityname"] = beego.AppConfig.String("ENGLISH_ENTER_ENTITY_NAME")
		c.Data["input_amount"] = beego.AppConfig.String("ENGLISH_ENTER_AMOUNT")

		c.TplName = "reports/merchanttransactionReport/merchanttransactionReport.html"
	} else if sess.Get("role") == "MERCHANT" && sess.Get("language") == "French" {
		c.Data["Menus1"] = template.HTML(`` + beego.AppConfig.String("FRENCH_USER_TEMPLATE") + ``)
		headerContent := strings.Replace(beego.AppConfig.String("MERCHANT_FRENCH_HEADER_TEMPLATE"), "{{.Uname}}", sess.Get("uname").(string), -1)

		c.Data["Header1"] = template.HTML(`` + headerContent + ``)

		c.Data["Dashboard"] = beego.AppConfig.String("FRENCH_DASHBOARD")
		c.Data["EntityTxnrep"] = beego.AppConfig.String("FRENCH_ENTITY_TXNREPORT")
		c.Data["TransactionReport"] = beego.AppConfig.String("FRENCH_TRANSACTION_REPORT")
		c.Data["SearchFilters"] = beego.AppConfig.String("FRENCH_SEARCH_FILTERS")

		c.Data["SearchResults"] = beego.AppConfig.String("FRENCH_USER_SEARCH_RESULTS")
		c.Data["SelectDateRange"] = beego.AppConfig.String("FRENCH_USER_SELECT_DATARANGE")

		c.Data["CNPSTxnDate"] = beego.AppConfig.String("FRENCH_CNPS_TXN_DATE")
		c.Data["PGTxnDate"] = beego.AppConfig.String("FRENCH_PG_TXN_DATE")
		c.Data["MerchantID"] = beego.AppConfig.String("FRENCH_MERCHANT_ID")
		c.Data["BankTxnDate"] = beego.AppConfig.String("FRENCH_BANK_TXN_DATE")

		c.Data["Status"] = beego.AppConfig.String("FRENCH_STATUS")
		c.Data["CNPSTransactionNumber"] = beego.AppConfig.String("FRENCH_CNPS_TRANSACTION_NUMBER")
		c.Data["PGTransactionNumber"] = beego.AppConfig.String("FRENCH_PG_TRANSACTION_NUMBER")
		c.Data["BANKTransactionNumber"] = beego.AppConfig.String("FRENCH_BANK_TRANSACTION_NUMBER")
		c.Data["MerchantName"] = beego.AppConfig.String("FRENCH_MERCHANT_NAME")
		c.Data["Amount"] = beego.AppConfig.String("FRENCH_AMOUNT")
		c.Data["Bank"] = beego.AppConfig.String("FRENCH_BANK")
		c.Data["Channel"] = beego.AppConfig.String("FRENCH_CHANNEL")
		c.Data["ListOfTransactionReports"] = beego.AppConfig.String("FRENCH_LIST_OF_TRANSACTION_REPORT")
		c.Data["Search"] = beego.AppConfig.String("FRENCH_SEARCH")
		c.Data["Reset"] = beego.AppConfig.String("FRENCH_RESET")
		c.Data["please_select"] = beego.AppConfig.String("FRENCH_PLEASESELECTCOMMON")
		c.Data["Approved"] = beego.AppConfig.String("FRENCH_APPROVED")
		c.Data["Declined"] = beego.AppConfig.String("FRENCH_DECLINED")
		c.Data["Pending"] = beego.AppConfig.String("FRENCH_PENDING")
		c.Data["Reports"] = beego.AppConfig.String("FRENCH_REPORTS")
		c.Data["input_cnpstxnnumnber"] = beego.AppConfig.String("FRENCH_ENTER_CNPS_TXN_ID")
		c.Data["input_pgtxnnumnber"] = beego.AppConfig.String("FRENCH_ENTER_PG_TXN_ID")
		c.Data["input_entityid"] = beego.AppConfig.String("FRENCH_ENTER_ENTITY_ID")
		c.Data["input_entityname"] = beego.AppConfig.String("FRENCH_ENTER_ENTITY_NAME")
		c.Data["input_amount"] = beego.AppConfig.String("FRENCH_ENTER_AMOUNT")

		c.TplName = "reports/merchanttransactionReport/merchanttransactionReport.html"
	}

	status := c.Input().Get("input_status")
	cnpstxnnumber := c.Input().Get("input_cnpstxnnumnber")
	pgtxnnumber := c.Input().Get("input_pgtxnnumnber")
	bank_name := c.Input().Get("input_bank")
	channel_name := c.Input().Get("input_channel")
	amount := c.Input().Get("input_amount")

	// status := c.Input().Get("status")
	// mask := c.Input().Get("mask")

	// c.Data["RRN"] = rrn
	// c.Data["Status"] = status
	// c.Data["Mask"] = mask
	c.Data["DateRange"] = c.Input().Get("daterange")
	tmp := strings.Split(c.Input().Get("daterange"), " - ")

	from := ""
	to := ""
	if len(tmp) == 2 {
		from = tmp[0]
		to = tmp[1]
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", status)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", from)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", to)

	uname := sess.Get("uname")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "User Name - ", uname)

	row, err := db.Db.Query(`select transaction_time,
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.UUID')),
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Amount')),
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Status')),
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Bank')),
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Channel')),
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Cnps_transaction_id')),
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.BankReferenceID')),
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Transaction_start_date')),
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.BankTransactionDate')) FROM Transactions WHERE (JSON_EXTRACT(transaction_deatils, '$.Entity_email') = ?) AND 
	(JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Status'))='' OR (JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Status'))) like ?) AND 
	(JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Cnps_transaction_id'))='' OR (JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Cnps_transaction_id'))) like ?) AND
	 (JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.UUID'))='' OR (JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.UUID'))) like ?) AND
	(JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Bank'))='' OR (JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Bank'))) like ?) AND
	(JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Channel'))='' OR (JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Channel'))) like ?) AND
	(JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Amount'))='' OR (JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Amount'))) like ?) AND
	DATE(transaction_time) >= STR_TO_DATE(?, '%m/%d/%Y') AND DATE(transaction_time) <= STR_TO_DATE(?, '%m/%d/%Y') ORDER BY transaction_time DESC`, sess.Get("uname").(string), status+"%", cnpstxnnumber+"%", pgtxnnumber+"%", bank_name+"%", channel_name+"%", amount+"%", from, to)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch data")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	var result []Row

	for i := range data {
		var r Row
		r.TimeStamp = data[i][0]
		r.PGTransactionNumber = data[i][1]
		r.Amount = data[i][2]
		r.Status = data[i][3]
		r.BankName = data[i][4]
		r.ChannelName = data[i][5]
		r.CNPSTransactionNumber = data[i][6]
		r.BankTransactionNumber = data[i][7]
		r.CNPSTxnDate = data[i][8]
		r.BankTransactionDate = data[i][9]

		result = append(result, r)
	}
	c.Data["CustomerData"] = result

	return
}
