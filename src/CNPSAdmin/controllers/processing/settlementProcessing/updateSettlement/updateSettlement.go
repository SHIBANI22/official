package updateSettlement

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"ominaya.com/database/sql"
	"ominaya.com/util/log"

	"html/template"

	"github.com/astaxie/beego"
)

type UpdateSettlement struct {
	beego.Controller
}

func (c *UpdateSettlement) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId - ", AdminId)
	var err error
	var Autherr error
	sessErr := false
	defer func() {

		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.Data["DisplayMessage"] = "Something went wrong.Please Contact CustomerCare."
			c.TplName = "error/error.html"
		}
		if Autherr != nil {
			c.Data["DisplayMessage"] = Autherr.Error()
			c.TplName = "error/error.html"
			return
		}
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)

			} else {
				c.Data["DisplayMessage"] = err.Error()
			}

			c.TplName = "processing/settlementProcessing/updateSettlement/updateSettlement.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "processing/settlementProcessing/updateSettlement/updateSettlement.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
		sessErr = true
		return
	}

	if err = session.ValidateSession(sess); err != nil {
		sess.SessionRelease(c.Ctx.ResponseWriter)
		session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		sessErr = true
		return
	}
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "TransactionProcessing")
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
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()

	c.Data["MenuJson"] = sess.Get("menujson")
	c.Data["language"] = sess.Get("language").(string)
	c.Data["Photo"] = sess.Get("photo").(string)
	successAmount, servicecharge, totalTransCount, totalBanks, successCount, pendingCount, declainedCount, transErr := utils.SideBarTransactionData()
	if transErr != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get sidebar data")
		return
	}
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
		c.Data["UpdateTransaction"] = beego.AppConfig.String("ENGLISH_UPDATE_TRANSACTION")
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
		c.Data["Submit"] = beego.AppConfig.String("ENGLISH_SUBMIT")
		c.Data["please_select"] = beego.AppConfig.String("ENGLISH_PLEASESELECTCOMMON")
		c.Data["English"] = beego.AppConfig.String("ENGLISH_ENGLISH")
		c.Data["French"] = beego.AppConfig.String("ENGLISH_FRENCH")
		c.Data["Update"] = beego.AppConfig.String("ENGLISH_UPDATE")
		c.Data["Back"] = beego.AppConfig.String("ENGLISH_BACK")
		c.Data["Alertmesg"] = beego.AppConfig.String("ENGLISH_ALERT_MSG")
		c.Data["BANKTransactionNumber"] = beego.AppConfig.String("ENGLISH_BANK_TRANSACTION_NUMBER")
		c.Data["Remarks"] = beego.AppConfig.String("ENGLISH_REMARKS")
		c.Data["Update"] = beego.AppConfig.String("ENGLISH_UPDATE")
		c.Data["Reporting"] = beego.AppConfig.String("ENGLISH_REPORTING_MENU")

		c.TplName = "processing/settlementProcessing/updateSettlement/updateSettlement.html"

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
		c.Data["UpdateTransaction"] = beego.AppConfig.String("FRENCH_UPDATE_TRANSACTION")
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
		c.Data["Submit"] = beego.AppConfig.String("FRENCH_SUBMIT")
		c.Data["English"] = beego.AppConfig.String("FRENCH_ENGLISH")
		c.Data["French"] = beego.AppConfig.String("FRENCH_FRENCH")
		c.Data["Update"] = beego.AppConfig.String("FRENCH_UPDATE")
		c.Data["Back"] = beego.AppConfig.String("FRENCH_BACK")
		c.Data["Alertmesg"] = beego.AppConfig.String("FRENCH_ALERT_MSG")
		c.Data["BANKTransactionNumber"] = beego.AppConfig.String("FRENCH_BANK_TRANSACTION_NUMBER")
		c.Data["Remarks"] = beego.AppConfig.String("FRENCH_REMARKS")
		c.Data["Update"] = beego.AppConfig.String("FRENCH_UPDATE")
		c.Data["Reporting"] = beego.AppConfig.String("FRENCH_REPORTING_MENU")

		c.TplName = "processing/settlementProcessing/updateSettlement/updateSettlement.html"

	}

	// uname = sess.Get("uname").(string)
	// c.Data["Uname"] = uname

	row, err := db.Db.Query(`select cnps_txn_number,transaction_time,pg_txn_number,
	entity_id,
	entity_name,
	amount,
	status,
	operator,
	channel,
	bank_txn_number,
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Transaction_start_date')),
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.BankTransactionDate')),
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Notification_url'))
	 FROM Transactions WHERE cnps_txn_number= ?`, AdminId)

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

	for i := range data {
		c.Data["Id"] = data[i][0]
		c.Data["cnps_txn_number"] = data[i][0]
		c.Data["transaction_time"] = data[i][1]
		c.Data["pg_txn_number"] = data[i][2]
		c.Data["entity_id"] = data[i][3]
		c.Data["entity_name"] = data[i][4]
		c.Data["amount"] = data[i][5]
		c.Data["status"] = data[i][6]
		c.Data["operator"] = data[i][7]
		c.Data["channel"] = data[i][8]
		c.Data["bank_txn_number"] = data[i][9]
		c.Data["cnps_txn_date"] = data[i][10]
		c.Data["bank_txn_date"] = data[i][11]
		c.Data["notification_url"] = data[i][12]

	}

	return

}

func (c *UpdateSettlement) Post() {

	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Start")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId - ", AdminId)

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	var derr error
	var Autherr error
	sessErr := false
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
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)
			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "processing/settlementProcessing/updateSettlement/updateSettlement.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin UserPage Fail")
		} else if derr != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)
			} else {
				c.Data["DisplayMessage"] = derr.Error()
			}
			c.TplName = "processing/settlementProcessing/updateSettlement/updateSettlement.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin UserPage Fail")

		} else {
			sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("System is unable to process your request.Please contact customer care")
				sessErr = true
				return
			}
			if sess.Get("language") == "English" {
				c.Data["DisplayMessage"] = beego.AppConfig.String("ENGLISH_TRANSACTION_UPDATED_SUCCESSFULLY")
				c.TplName = "processing/settlementProcessing/updateSettlement/updateSettlement.html"

				log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User  Page Success")
			} else if sess.Get("language") == "French" {
				c.Data["DisplayMessage"] = beego.AppConfig.String("FRENCH_TRANSACTION_UPDATED_SUCCESSFULLY")
				c.TplName = "processing/settlementProcessing/updateSettlement/updateSettlement.html"

				log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User  Page Success")
			}
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
		sessErr = true
		return
	}

	if err = session.ValidateSession(sess); err != nil {
		sess.SessionRelease(c.Ctx.ResponseWriter)
		session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		sessErr = true
		return
	}
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "TransactionProcessing")
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
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()
	c.Data["MenuJson"] = sess.Get("menujson")
	c.Data["language"] = sess.Get("language").(string)
	c.Data["Photo"] = sess.Get("photo").(string)
	successAmount, servicecharge, totalTransCount, totalBanks, successCount, pendingCount, declainedCount, transErr := utils.SideBarTransactionData()
	if transErr != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get sidebar data")
		return
	}
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
		c.Data["UpdateTransaction"] = beego.AppConfig.String("ENGLISH_UPDATE_TRANSACTION")
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
		c.Data["Submit"] = beego.AppConfig.String("ENGLISH_SUBMIT")
		c.Data["please_select"] = beego.AppConfig.String("ENGLISH_PLEASESELECTCOMMON")
		c.Data["English"] = beego.AppConfig.String("ENGLISH_ENGLISH")
		c.Data["French"] = beego.AppConfig.String("ENGLISH_FRENCH")
		c.Data["Update"] = beego.AppConfig.String("ENGLISH_UPDATE")
		c.Data["Back"] = beego.AppConfig.String("ENGLISH_BACK")
		c.Data["Alertmesg"] = beego.AppConfig.String("ENGLISH_ALERT_MSG")
		c.Data["BANKTransactionNumber"] = beego.AppConfig.String("ENGLISH_BANK_TRANSACTION_NUMBER")
		c.Data["Remarks"] = beego.AppConfig.String("FRENCH_REMARKS")
		c.Data["Update"] = beego.AppConfig.String("ENGLISH_UPDATE")
		c.Data["Reporting"] = beego.AppConfig.String("ENGLISH_REPORTING_MENU")

		c.TplName = "processing/settlementProcessing/updateSettlement/updateSettlement.html"

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
		c.Data["UpdateTransaction"] = beego.AppConfig.String("FRENCH_UPDATE_TRANSACTION")

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
		c.Data["Submit"] = beego.AppConfig.String("FRENCH_SUBMIT")
		c.Data["English"] = beego.AppConfig.String("FRENCH_ENGLISH")
		c.Data["French"] = beego.AppConfig.String("FRENCH_FRENCH")
		c.Data["Update"] = beego.AppConfig.String("FRENCH_UPDATE")
		c.Data["Back"] = beego.AppConfig.String("FRENCH_BACK")
		c.Data["Alertmesg"] = beego.AppConfig.String("FRENCH_ALERT_MSG")
		c.Data["BANKTransactionNumber"] = beego.AppConfig.String("FRENCH_BANK_TRANSACTION_NUMBER")
		c.Data["Remarks"] = beego.AppConfig.String("FRENCH_REMARKS")
		c.Data["Update"] = beego.AppConfig.String("FRENCH_UPDATE")
		c.Data["Reporting"] = beego.AppConfig.String("FRENCH_REPORTING_MENU")

		c.TplName = "processing/settlementProcessing/updateSettlement/updateSettlement.html"

	}

	cnpsid := c.Input().Get("input_cnps_txn_number")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cnpsid - ", cnpsid)

	txn_status := c.Input().Get("input_txn_status")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_txn_status - ", txn_status)

	remarks := c.Input().Get("input_remarks")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "remarks - ", remarks)

	bank_txn := c.Input().Get("input_bank_txn")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "bank_txn - ", bank_txn)

	creater_email := sess.Get("uname")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "creater_email - ", creater_email)

	if txn_status == "" && sess.Get("language") == "English" {
		err = errors.New(beego.AppConfig.String("ENGLISH_TRANSACTION_STATUS_NOT_EMPTY"))
		return
	} else if txn_status == "" && sess.Get("language") == "French" {
		err = errors.New(beego.AppConfig.String("FRENCH_TRANSACTION_STATUS_NOT_EMPTY"))
		return
	}
	if bank_txn == "" && sess.Get("language") == "English" {
		err = errors.New(beego.AppConfig.String("ENGLISH_TRANSACTION_STATUS_NOT_EMPTY"))
		return
	} else if bank_txn == "" && sess.Get("language") == "French" {
		err = errors.New(beego.AppConfig.String("FRENCH_TRANSACTION_STATUS_NOT_EMPTY"))
		return
	}

	if remarks == "" && sess.Get("language") == "English" {
		err = errors.New(beego.AppConfig.String("ENGLISH_TRANSACTION_STATUS_NOT_EMPTY"))
		return
	} else if remarks == "" && sess.Get("language") == "French" {
		err = errors.New(beego.AppConfig.String("FRENCH_TRANSACTION_STATUS_NOT_EMPTY"))
		return
	}

	// if txn_status == "APPROVED" {
	// 	aa := "0000"
	// } else if txn_status == "DECLINED" {
	// 	aa := "1111"
	// }

	currentTime := time.Now()

	result, err := db.Db.Exec(`UPDATE Transactions SET
	transaction_deatils= JSON_REPLACE(transaction_deatils, '$.Code',?),
	transaction_deatils= JSON_REPLACE(transaction_deatils, '$.Message',?),
	transaction_deatils= JSON_REPLACE(transaction_deatils, '$.BankTransactionDate',?),
	transaction_deatils= JSON_INSERT(transaction_deatils, '$.TransactionUpdatedby',?),
	transaction_deatils= JSON_INSERT(transaction_deatils, '$.TransactionProcess','ManualReversal'),
	 status =?,
	bank_txn_number=?
	WHERE cnps_txn_number =?`,
		txn_status, remarks, currentTime, creater_email, txn_status, bank_txn, cnpsid)
	if err != nil {
		err = errors.New("Transaction updation failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Transaction updation failed")
		return
	} else {

		row, err := db.Db.Query(`select
	cnps_txn_number,
	entity_id,
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Cnps_declaration_number')),
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Declaration_period')),
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Description')),
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Customer_mobile_number')),
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Partial_payment')),
	pg_txn_number,
	status,
	transaction_time,
	amount,
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Currency')),
	operator,
	channel,
	bank_txn_number,
	JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Notification_url'))
	FROM Transactions WHERE cnps_txn_number= ?`, cnpsid)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get the PGS info")
			return
		}
		defer sql.Close(row)
		_, data, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch data")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "pavan data", data)

		for i := range data {

			un, err := url.Parse(data[i][15])
			if err != nil {
				beego.Error(err)
				c.Abort("500")
			}

			paramsn := url.Values{}

			paramsn.Set("cnps_transaction_id", data[i][0])
			paramsn.Set("cnps_entity_id", data[i][1])
			paramsn.Set("cnps_declaration_id", data[i][2])
			paramsn.Set("declaration_period", data[i][3])
			paramsn.Set("description", data[i][4])
			paramsn.Set("customer_mobile_number", data[i][5])
			paramsn.Set("partial_payment", data[i][6])
			paramsn.Set("transaction_end_date", "") //TBD
			paramsn.Set("pgs_transaction_id", data[i][7])
			if txn_status == "APPROVED" {
				paramsn.Set("code", "0000")
			} else if txn_status == "DECLINED" {
				paramsn.Set("code", "1111")
			}

			paramsn.Set("status", txn_status)
			paramsn.Set("message", remarks)
			paramsn.Set("pgs_transaction_date", data[i][9])
			paramsn.Set("amount", data[i][10])
			paramsn.Set("currecncy", data[i][11])
			paramsn.Set("paymentMethod", data[i][12])
			paramsn.Set("paymentMode", data[i][13])
			paramsn.Set("paymentRefNo", bank_txn)

			//below block to get hash of parameters
			//signParam := l.UUID + l.Amount + l.Status + l.Cnps_transaction_id + l.Channel + beego.AppConfig.String("SKEY") + l.Cnps_entity_id
			signParam := data[i][7] + data[i][10] + txn_status + data[i][0] + data[i][13] + beego.AppConfig.String("SKEY") + data[i][1]

			paramsn.Set("sign", MakesignofParam(signParam, beego.AppConfig.String("SKEY")))
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "Sign Param", signParam)

			un.RawQuery = paramsn.Encode()
			beego.Debug(un.String())

			u, _ := url.ParseRequestURI(data[i][15])
			urlStr := u.String()

			client := &http.Client{}
			r, _ := http.NewRequest("POST", urlStr, strings.NewReader(paramsn.Encode())) // URL-encoded payload
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Add("Content-Length", strconv.Itoa(len(paramsn.Encode())))

			resp, err := client.Do(r)
			if err != nil {
				beego.Error(err)
				return
			}

			defer resp.Body.Close()

			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("response Body:", string(body))

		}

	}

	return
}

func MakesignofParam(InputString, skey string) (sign string) {

	input := InputString
	hmac512 := hmac.New(sha512.New, []byte(skey))
	hmac512.Write([]byte(input))

	//4db45e622c0ae3157bdcb53e436c96c5
	//fmt.Printf("md5:\t\t%x\n", md5.Sum(nil))

	//eb7a03c377c28da97ae97884582e6bd07fa44724af99798b42593355e39f82cb
	//fmt.Printf("sha256:\t\t%x\n", sha_256.Sum(nil))

	//5cdaf0d2f162f55ccc04a8639ee490c94f2faeab3ba57d3c50d41930a67b5fa6915a73d6c78048729772390136efed25b11858e7fc0eed1aa7a464163bd44b1c
	//fmt.Printf("sha512:\t\t%x\n", sha_512.Sum(nil))

	//34c614af69a2550a4d39138c3756e2cc50b4e5495af3657e5b726c2ac12d5e60
	//fmt.Printf("sha512_256:\t%x\n", sha_512_256)

	//GBZ7aqtVzXGdRfdXLHkb0ySp/f+vV9Zo099N+aSv+tTagUWuHrPeECDfUyd5WCoHBe7xkw2EdpyLWx+Ge4JQKg==

	fmt.Printf("hmac512:\t%s\n", base64.StdEncoding.EncodeToString(hmac512.Sum(nil)))
	sign = base64.StdEncoding.EncodeToString(hmac512.Sum(nil))
	return
}
