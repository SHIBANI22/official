package updateMerchant

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"ominaya.com/database/sql"
	"ominaya.com/util/log"

	"html/template"

	"github.com/astaxie/beego"
)

type UpdateMerchant struct {
	beego.Controller
}

func (c *UpdateMerchant) Get() {
	merchantId := c.Ctx.Input.Param(":MerchantID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Merchant Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "merchantId - ", merchantId)
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

			c.TplName = "merchant/updateMerchant/updateMerchant.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "merchant/updateMerchant/updateMerchant.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Success")
		}
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "UpdateMerchant")
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
		c.Data["MerchantApproval"] = beego.AppConfig.String("ENGLISH_MERCHANT_APPROVAL")
		c.Data["Merchant"] = beego.AppConfig.String("ENGLISH_MERCHANT")
		c.Data["SettlementCurrency"] = beego.AppConfig.String("ENGLISH_SETTLEMENT_CURRENCY")
		c.Data["ContactInfo"] = beego.AppConfig.String("ENGLISH_CONTACT_INFO")
		c.Data["Countrys"] = beego.AppConfig.String("ENGLISH_COUNTRY")
		c.Data["Email"] = beego.AppConfig.String("ENGLISH_EMAIL")
		c.Data["Mobile"] = beego.AppConfig.String("ENGLISH_MOBILE")
		c.Data["Address"] = beego.AppConfig.String("ENGLISH_Address")
		c.Data["FullName"] = beego.AppConfig.String("ENGLISH_FULL_NAME")
		c.Data["BusinessDtls"] = beego.AppConfig.String("ENGLISH_BUSINESS_DATAILS")
		c.Data["Logo"] = beego.AppConfig.String("ENGLISH_LOGO")
		c.Data["BusinessName"] = beego.AppConfig.String("ENGLISH_BUSINESS_NAME")
		c.Data["BusinessType"] = beego.AppConfig.String("ENGLISH_BUSINESS_TYPE")
		c.Data["AccountDetails"] = beego.AppConfig.String("ENGLISH_ACCOUNT_DETAILS")
		c.Data["PAN"] = beego.AppConfig.String("ENGLISH_PAN")
		c.Data["AccountName"] = beego.AppConfig.String("ENGLISH_ACCOUNT_NAME")
		c.Data["IFSCCode"] = beego.AppConfig.String("ENGLISH_IFSC_CODE")
		c.Data["BankName"] = beego.AppConfig.String("ENGLISH_BANK_NAME")

		c.Data["DOC"] = beego.AppConfig.String("ENGLISH_DOCUMENTS_VERIFICATION")
		c.Data["Id_Proof"] = beego.AppConfig.String("ENGLISH_ID_PROOF")
		c.Data["Address_Proof"] = beego.AppConfig.String("ENGLISH_ADDRESS_PROOF")
		c.TplName = "merchant/updateMerchant/updateMerchant.html"

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
		c.Data["MerchantApproval"] = beego.AppConfig.String("FRENCH_MERCHANT_APPROVAL")
		c.Data["Merchant"] = beego.AppConfig.String("FRENCH_MERCHANT")
		c.Data["SettlementCurrency"] = beego.AppConfig.String("FRENCH_SETTLEMENT_CURRENCY")
		c.Data["ContactInfo"] = beego.AppConfig.String("FRENCH_CONTACT_INFO")
		c.Data["Mobile"] = beego.AppConfig.String("FRENCH_MOBILE")
		c.Data["Email"] = beego.AppConfig.String("FRENCH_EMAIL")
		c.Data["Countrys"] = beego.AppConfig.String("FRENCH_COUNTRY")
		c.Data["Address"] = beego.AppConfig.String("FRENCH_Address")
		c.Data["FullName"] = beego.AppConfig.String("FRENCH_FULL_NAME")
		c.Data["BusinessDtls"] = beego.AppConfig.String("FRENCH_BUSINESS_DATAILS")
		c.Data["Logo"] = beego.AppConfig.String("FRENCH_LOGO")
		c.Data["BusinessName"] = beego.AppConfig.String("FRENCH_BUSINESS_NAME")
		c.Data["BusinessType"] = beego.AppConfig.String("FRENCH_BUSINESS_TYPE")
		c.Data["AccountDetails"] = beego.AppConfig.String("FRENCH_ACCOUNT_DETAILS")
		c.Data["PAN"] = beego.AppConfig.String("FRENCH_PAN")
		c.Data["AccountName"] = beego.AppConfig.String("FRENCH_ACCOUNT_NAME")
		c.Data["IFSCCode"] = beego.AppConfig.String("FRENCH_IFSC_CODE")
		c.Data["BankName"] = beego.AppConfig.String("FRENCH_BANK_NAME")

		c.Data["DOC"] = beego.AppConfig.String("FRENCH_DOCUMENTS_VERIFICATION")
		c.Data["Id_Proof"] = beego.AppConfig.String("FRENCH_ID_PROOF")
		c.Data["Address_Proof"] = beego.AppConfig.String("FRENCH_ADDRESS_PROOF")
		c.TplName = "merchant/updateMerchant/updateMerchant.html"

	}
	err = dataSet(c, merchantId)

	data, err := getCurrencies()
	c.Data["CountryNames"] = data

}
func getCurrencies() (data [][]string, err error) {

	row, err := db.Db.Query(`SELECT symbol, country FROM Currency`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get dashboard data")
		//c.Data["CountryNames"] = data
		return data, nil
	}
	return
}

// func getCurrencies1(symbol string)(data [][]string,err error){

// 	row, err := db.Db.Query(`SELECT symbol, country FROM Currency WHERE symbol=?`,symbol)

// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 		err = errors.New("Unable to get dashboard data")
// 		return
// 	}
// 	defer sql.Close(row)
// 	_, data, err = sql.Scan(row)
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 		err = errors.New("Unable to fetch data")
// 		return
// 	}
// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
// 	if len(data) <= 0 {
// 		err = errors.New("Unable to get dashboard data")

//			return data, nil
//		}
//		return
//	}
func (c *UpdateMerchant) Post() {

	merchantId := c.Ctx.Input.Param(":MerchantID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Merchant Page Start")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "merchantId - ", merchantId)

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Merchant Page Start")
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
			c.TplName = "merchant/updateMerchant/updateMerchant.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin UserPage Fail")
		} else if derr != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)
			} else {
				c.Data["DisplayMessage"] = derr.Error()
			}
			c.TplName = "merchant/updateMerchant/updateMerchant.html"

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
				c.Data["DisplayMessage"] = "Merchant Updated Successfully"
				c.TplName = "merchant/updateMerchant/updateMerchant.html"

				log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User  Page Success")
			} else if sess.Get("language") == "French" {
				c.Data["DisplayMessage"] = "Merchant Updated Successfully"
				c.TplName = "merchant/updateMerchant/updateMerchant.html"

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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "UpdateMerchant")
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
		c.Data["MerchantApproval"] = beego.AppConfig.String("ENGLISH_MERCHANT_APPROVAL")
		c.Data["Merchant"] = beego.AppConfig.String("ENGLISH_MERCHANT")
		c.Data["Reporting"] = beego.AppConfig.String("ENGLISH_REPORTING_MENU")
		c.Data["SettlementCurrency"] = beego.AppConfig.String("ENGLISH_SETTLEMENT_CURRENCY")
		c.Data["ContactInfo"] = beego.AppConfig.String("ENGLISH_CONTACT_INFO")
		c.Data["Countrys"] = beego.AppConfig.String("ENGLISH_COUNTRY")
		c.Data["Email"] = beego.AppConfig.String("ENGLISH_EMAIL")
		c.Data["Mobile"] = beego.AppConfig.String("ENGLISH_MOBILE")
		c.Data["Address"] = beego.AppConfig.String("ENGLISH_Address")
		c.Data["FullName"] = beego.AppConfig.String("ENGLISH_FULL_NAME")
		c.Data["BusinessDtls"] = beego.AppConfig.String("ENGLISH_BUSINESS_DATAILS")
		c.Data["Logo"] = beego.AppConfig.String("ENGLISH_LOGO")
		c.Data["BusinessName"] = beego.AppConfig.String("ENGLISH_BUSINESS_NAME")
		c.Data["BusinessType"] = beego.AppConfig.String("ENGLISH_BUSINESS_TYPE")
		c.Data["AccountDetails"] = beego.AppConfig.String("ENGLISH_ACCOUNT_DETAILS")
		c.Data["PAN"] = beego.AppConfig.String("ENGLISH_PAN")
		c.Data["AccountName"] = beego.AppConfig.String("ENGLISH_ACCOUNT_NAME")
		c.Data["IFSCCode"] = beego.AppConfig.String("ENGLISH_IFSC_CODE")
		c.Data["BankName"] = beego.AppConfig.String("ENGLISH_BANK_NAME")

		c.Data["DOC"] = beego.AppConfig.String("ENGLISH_DOCUMENTS_VERIFICATION")
		c.Data["Id_Proof"] = beego.AppConfig.String("ENGLISH_ID_PROOF")
		c.Data["Address_Proof"] = beego.AppConfig.String("ENGLISH_ADDRESS_PROOF")
		c.TplName = "merchant/updateMerchant/updateMerchant.html"

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
		c.Data["MerchantApproval"] = beego.AppConfig.String("FRENCH_MERCHANT_APPROVAL")
		c.Data["Merchant"] = beego.AppConfig.String("FRENCH_MERCHANT")
		c.Data["Reporting"] = beego.AppConfig.String("FRENCH_REPORTING_MENU")
		c.Data["SettlementCurrency"] = beego.AppConfig.String("FRENCH_SETTLEMENT_CURRENCY")

		c.Data["ContactInfo"] = beego.AppConfig.String("FRENCH_CONTACT_INFO")
		c.Data["Mobile"] = beego.AppConfig.String("FRENCH_MOBILE")
		c.Data["Email"] = beego.AppConfig.String("FRENCH_EMAIL")
		c.Data["Countrys"] = beego.AppConfig.String("FRENCH_COUNTRY")
		c.Data["Address"] = beego.AppConfig.String("FRENCH_Address")
		c.Data["FullName"] = beego.AppConfig.String("FRENCH_FULL_NAME")
		c.Data["BusinessDtls"] = beego.AppConfig.String("FRENCH_BUSINESS_DATAILS")
		c.Data["Logo"] = beego.AppConfig.String("FRENCH_LOGO")
		c.Data["BusinessName"] = beego.AppConfig.String("FRENCH_BUSINESS_NAME")
		c.Data["BusinessType"] = beego.AppConfig.String("FRENCH_BUSINESS_TYPE")
		c.Data["AccountDetails"] = beego.AppConfig.String("FRENCH_ACCOUNT_DETAILS")
		c.Data["PAN"] = beego.AppConfig.String("FRENCH_PAN")
		c.Data["AccountName"] = beego.AppConfig.String("FRENCH_ACCOUNT_NAME")
		c.Data["IFSCCode"] = beego.AppConfig.String("FRENCH_IFSC_CODE")
		c.Data["BankName"] = beego.AppConfig.String("FRENCH_BANK_NAME")

		c.Data["DOC"] = beego.AppConfig.String("FRENCH_DOCUMENTS_VERIFICATION")
		c.Data["Id_Proof"] = beego.AppConfig.String("FRENCH_ID_PROOF")
		c.Data["Address_Proof"] = beego.AppConfig.String("FRENCH_ADDRESS_PROOF")

		c.TplName = "merchant/updateMerchant/updateMerchant.html"

	}

	status := c.Input().Get("input_status")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_txn_status - ", status)

	if status == "" && sess.Get("language") == "English" {
		err = errors.New(beego.AppConfig.String("ENGLISH_TRANSACTION_STATUS_NOT_EMPTY"))
		return
	} else if status == "" && sess.Get("language") == "French" {
		err = errors.New(beego.AppConfig.String("FRENCH_TRANSACTION_STATUS_NOT_EMPTY"))
		return
	}

	result, err := db.Db.Exec(`UPDATE Merchants SET status =? WHERE merchantId =?`, status, merchantId)
	if err != nil {
		err = errors.New("Merchant updation failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Merchant updation failed")
		// return
	}

	symbol := c.Input().Get("settlementcurrency")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "code - ", symbol)

	data, err := getCode(symbol)
	if err != nil {
		return
	}

	Result, err := db.Db.Exec(`UPDATE Merchants SET settlement_currency =? WHERE merchantId =?`, data[0][0], merchantId)

	if err != nil {
		err = errors.New("Adding Merchant Settlement Currency failed")
		return
	}

	n, err := Result.RowsAffected()
	if err != nil || n == 0 {
		err = errors.New("Adding Merchant Settlement Currency failed")
		//return
	}
	err = dataSet(c, merchantId)
	data, err = getCurrencies()
	//data, err = getCurrencies1(symbol)
	c.Data["CountryNames"] = data
	return

}

func dataSet(c *UpdateMerchant, merchantId string) (err error) {
	row, err := db.Db.Query(`select id,merchantId,name,
	logo,
	status,
	settlement_currency,
	mobileNum,
	emailId,
	address,
	country,
	package_id,
	account_num,
	account_ifsc_code,
	bank_name,
	business_name,
	bussiness_type,
	capital,
	account_swift_code,
	id_proof_doc,
	address_proof_doc,
	id_proof_num,
	address_proof_num,
	pan
	 FROM Merchants WHERE merchantId= ?`, merchantId)

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

	for i := range data {
		idData := []byte(data[i][18])
		addData := []byte(data[i][19])
		logodata := []byte(data[i][3])
		id := base64.StdEncoding.EncodeToString(idData)
		add := base64.StdEncoding.EncodeToString(addData)
		if len(logodata) != 0 {
			logo := base64.StdEncoding.EncodeToString(logodata)
			c.Data["logo"] = logo
		} else {
			c.Data["logo"] = logodata
		}
		c.Data["id"] = data[i][0]
		c.Data["merchantId"] = data[i][1]
		c.Data["name"] = data[i][2]
		c.Data["status"] = data[i][4]
		c.Data["settlement_currency"] = data[i][5]
		c.Data["mobileNum"] = data[i][6]
		c.Data["emailId"] = data[i][7]
		c.Data["address"] = data[i][8]
		c.Data["country"] = data[i][9]
		c.Data["package_id"] = data[i][10]
		c.Data["account_num"] = data[i][11]
		c.Data["account_ifsc_code"] = data[i][12]
		c.Data["bank_name"] = data[i][13]
		c.Data["business_name"] = data[i][14]
		c.Data["bussiness_type"] = data[i][15]
		c.Data["capital"] = data[i][16]
		c.Data["account_swift_code"] = data[i][17]
		c.Data["id_proof_doc"] = id
		c.Data["address_proof_doc"] = add
		c.Data["id_proof_num"] = data[i][20]
		c.Data["address_proof_num"] = data[i][21]
		c.Data["pan"] = data[i][22]

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

func getCode(symbol string) (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT code FROM Currency where symbol=?`, symbol)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Currency data")
		return
	}
	defer sql.Close(row)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Currency data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) == 0 {
		err = errors.New("Unable to get Currency data")
	}
	return
}
