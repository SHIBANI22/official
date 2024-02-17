package transactionlogs

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"errors"
	"fmt"
	"html/template"
	"runtime/debug"
	"time"

	"strings"

	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/jung-kurt/gofpdf"
	"ominaya.com/database/sql"
	"ominaya.com/util/log"
)

type TransactionLogs struct {
	beego.Controller
}

type CNPSRequest struct {
	Version                 string `form:"version"`
	Language                string `form:"language"`
	Cnps_entity_id          string `form:"cnps_entity_id"`
	Entity_name             string `form:"entity_name"`
	Entity_phone            string `form:"entity_phone"`
	Entity_email            string `form:"entity_email"`
	Cnps_declaration_number string `form:"cnps_declaration_number"` //TBD
	Declaration_period      string `form:"declaration_period"`      //TBD
	Amount                  string `form:"amount"`
	Currency                string `form:"currency"`
	Nature_code             string `form:"nature_code"`
	Nature_name             string `form:"nature_name"`
	Cnps_transaction_id     string `form:"cnps_transaction_id"` // CNPS_Request_ID
	Customer_mobile_number  string `form:"customer_mobile_number"`
	Customer_account_number string `form:"customer_account_number"`
	Partial_payment         string `form:"partial_payment"`
	Transaction_start_date  string `form:"transaction_start_date"`
	Notification_url        string `form:"notification_url"`
	Return_url              string `form:"return_url"`
	Cancel_url              string `form:"cancel_url"`
	Sign                    string `form:"sign"`
	Description             string `form:"description"`
	User_id                 string `form:"user_id"`
	Declaration_type        string `form:"declaration_type"`
	Payment_mode            string `form:"payment_mode"`
	Payment_method          string `form:"payment_method"`
	Payment_type            string `form:"payment_type"`

	SourceExchange             string
	SourceKey                  string
	UUID                       string
	Code                       string
	Message                    string
	TxnChannelURL              string
	RedircetURLPG              string
	ReturnURLPG                string
	NotificationURLPG          string
	CancelURLPG                string
	Status                     string
	Bank                       string
	Channel                    string
	Pgstatus                   string
	MerchantStatus             string
	Merchant                   string
	UBACardRef                 string
	UBACardResponseStatus1     int
	UBACardResponseStatus2     string
	UBACardResponseInfo        string
	UBACardResponseReferenceNo string
	BankReferenceID            string
	BankTransactionDate        string
	Bnk                        string
	ResponseCode               string
	ResponseMessage            string
	BillMapTransactionId       string
	EWPTransactionId           string
	UBADirectDebitCode         string
	TokenNumber                string
	Securepack                 string
	Timedata                   string
	DynamicQRBase64            string
	OrangeStatus               string
	OrangeMessage              string
	OrangeCode                 string
	OrangeDescription          string
	OrangePay_token            string
	OrangePayment_url          string
	OrangeNotif_token          string
	UBASecretkey               string
	MSISDN                     string
}

type MerchantRequest struct {
	Version                 string `json:"Version"`
	Language                string `json:"language"`
	Cnps_entity_id          string `json:"cnps_entity_id"`
	Entity_name             string `json:"entity_name"`
	Entity_phone            string `json:"entity_phone"`
	Entity_email            string `json:"entity_email"`
	Cnps_declaration_number string `json:"cnps_declaration_number"`
	Declaration_period      string `json:"declaration_period"`
	Amount                  string `json:"amount"`
	Currency                string `json:"currency"`
	Nature_code             string `json:",omitempty"`
	Nature_name             string `json:",omitempty"`
	Cnps_transaction_id     string `json:"Cnps_transaction_id"`
	Customer_mobile_number  string `json:"Customer_mobile_number"`
	Partial_payment         string `json:"Partial_payment"`
	Transaction_start_date  string `json:"Transaction_start_date"`
	Notification_url        string `json:"Notification_url"`
	Return_url              string `json:"Return_url"`
	Cancel_url              string `json:"Cancel_url"`
	Sign                    string `json:"Sign"`
	Description             string `json:"Description"`
	User_id                 string `json:"User_id"`
	Declaration_type        string `json:"Declaration_type"`
	Payment_method          string `json:",omitempty"`
	Payment_mode            string `json:",omitempty"`
}

type MerchantNotify struct {
	Cnps_transaction_id     string `json:"cnps_transaction_id"`
	Cnps_entity_id          string `json:"cnps_entity_id"`
	Cnps_declaration_number string `json:"cnps_declaration_id"`
	Declaration_period      string `json:"declaration_period"`
	Description             string `json:"description"`
	Customer_mobile_number  string `json:"customer_mobile_number"`
	Partial_payment         string `json:"partial_payment"`
	Pgs_transaction_id      string `json:"pgs_transaction_id"`
	Code                    string `json:"code"`
	Status                  string `json:"status"`
	Message                 string `json:"message"`
	Transaction_start_date  string `json:"pgs_transaction_date"`
	Amount                  string `json:"amount"`
	Currency                string `json:"currency"`
	paymentMethod           string `json:"paymentMethod"`
	paymentMode             string `json:"paymentMode"`
	BankReferenceID         string `json:"paymentRefNo"`
	BankTransactionDate     string `json:"transaction_end_date"`
}

type GetBillReq struct {
	Username          string `json:"Username"`
	Password          string `json:"Password"`
	CNPSTransactionID string `json:"CNPSTransactionID"`
	BANKTXNID         string `json:",omitempty"`
	ECOTXNID          string `json:",omitempty"`
}

type PaymentNotify struct {
	Username          string `json:"Username"`
	Password          string `json:"Password"`
	CNPSTransactionID string `json:"CNPSTransactionID"`
	BANKTXNID         string `json:",omitempty"`
	ECOTXNID          string `json:",omitempty"`
	PGWTxnID          string `json:"PGWTxnID"`
	CNPSEntityID      string `json:"CNPSEntityID"`
	EntityName        string `json:"EntityName"`
	Amount            string `json:"Amount"`
	Currency          string `json:"Currency"`
	PaymentRefNo      string `json:"PaymentRefNo"`
	Status            string `json:"Status"`
	DeclarationType   string `json:"DeclarationType"`
	PaymentMode       string `json:"PaymentMode"`
	PaymentMethod     string `json:"PaymentMethod"`
}

type IPNResponse struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
	Status  string `json:"Status"`
}

type MTNPaymentReq struct {
	Code      string `json:"Code"`
	Password  string `json:"Password"`
	MSISDN    string `json:"MSISDN"`
	Reference string `json:"Reference"`
	Amount    string `json:"Amount"`
	MetaData  string `json:"MetaData"`
}

type MTNPaymentRes struct {
	ResponseCode         string `json:"ResponseCode"`
	ResponseMessage      string `json:"ResponseMessage"`
	BillMapTransactionId string `json:"BillMapTransactionId"`
	EWPTransactionId     string `json:"EWPTransactionId"`
}

type OrangePaymentRes struct {
	OrangeStatus      string `json:"OrangeStatus"`
	OrangeMessage     string `json:"OrangeMessage"`
	OrangeCode        string `json:"OrangeCode"`
	OrangeDescription string `json:"OrangeDescription"`
	OrangePayment_url string `json:"OrangePayment_url"`
}

func (c *TransactionLogs) Get() {
	//AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Log Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	sessErr := false
	var Autherr error
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
			c.TplName = "reports/transactionLogs/transactionLogs.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Log Page Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "reports/transactionLogs/transactionLogs.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Log Page Success")
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "DebugLogReport")
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
		c.Data["SearchFilters"] = beego.AppConfig.String("ENGLISH_SEARCH_FILTERS")
		c.Data["CNPSTransactionNumber"] = beego.AppConfig.String("ENGLISH_CNPS_TRANSACTION_NUMBER")
		c.Data["Bank"] = beego.AppConfig.String("ENGLISH_BANK")
		c.Data["Reset"] = beego.AppConfig.String("ENGLISH_RESET")
		c.Data["please_select"] = beego.AppConfig.String("ENGLISH_PLEASESELECTCOMMON")
		c.Data["enterTXN"] = beego.AppConfig.String("ENGLISH_ENTER_CNPS_TXN_ID")
		c.Data["GetLogs"] = beego.AppConfig.String("ENGLISH_GET_LOGS")
		c.Data["TxnLogDetals"] = beego.AppConfig.String("ENGLISH_TRANSACTION_LOGS")
		c.Data["Reporting"] = beego.AppConfig.String("ENGLISH_SYSTEM_REPORTS")
		c.Data["Back"] = beego.AppConfig.String("ENGLISH_BACK_BUTTON")

		c.TplName = "reports/transactionLogs/transactionLogs.html"
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
		c.Data["SearchFilters"] = beego.AppConfig.String("FRENCH_SEARCH_FILTERS")
		c.Data["CNPSTransactionNumber"] = beego.AppConfig.String("FRENCH_CNPS_TRANSACTION_NUMBER")
		c.Data["Bank"] = beego.AppConfig.String("FRENCH_BANK")
		c.Data["Reset"] = beego.AppConfig.String("FRENCH_RESET")
		c.Data["please_select"] = beego.AppConfig.String("FRENCH_PLEASESELECTCOMMON")
		c.Data["enterTXN"] = beego.AppConfig.String("FRENCH_ENTER_CNPS_TXN_ID")
		c.Data["TxnLogDetals"] = beego.AppConfig.String("FRENCH_TRANSACTION_LOGS")
		c.Data["Reporting"] = beego.AppConfig.String("FRENCH_SYSTEM_REPORTS")
		c.Data["GetLogs"] = beego.AppConfig.String("FRENCH_GET_LOGS")
		c.Data["Back"] = beego.AppConfig.String("FRENCH_BACK_BUTTON")

		c.TplName = "reports/transactionLogs/transactionLogs.html"
	}

	return

}

func (c *TransactionLogs) Post() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Log Page Start")
	pip := c.Ctx.Input.IP()
	cnpstxnnumber := c.Input().Get("input_cnpstxnnumber")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	var req CNPSRequest
	sessErr := false
	var Autherr error
	var operatorLogs, merchantLogs, username string
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
			c.TplName = "reports/transactionLogs/transactionLogs.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Log Page Fail")
		} else {
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("System is unable to process your request.Please contact customer care")
				sessErr = true
				return
			}
			// filename := "Dubug_Reports_" + cnpstxnnumber + ".txt"
			// c.Ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename="+filename)
			// c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/text")
			// c.Ctx.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(operatorLogs)))
			// finaldata := strings.NewReader(operatorLogs)
			// //Send the file
			// io.Copy(c.Ctx.ResponseWriter, finaldata) //'Copy' the file to the client
			filename := "Dubug_Reports_" + cnpstxnnumber + ".pdf"
			c.Ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename="+filename)
			c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/pdf")
			//c.Ctx.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(operatorLogs)))
			//finaldata := strings.NewReader(operatorLogs)
			//Send the file
			//io.Copy(c.Ctx.ResponseWriter, finaldata) //'Copy' the file to the client

			pdf := gofpdf.New("P", "mm", "A4", "")
			pdf.SetHeaderFunc(func() {
				pdf.Image("static/assets/media/image/logo.png", 10, 6, 30, 0, false, "", 0, "")
				pdf.SetY(5)
				pdf.SetFont("Arial", "B", 18)
				pdf.Cell(80, 0, "")
				pdf.CellFormat(30, 10, "Transaction Logs", "0", 0, "C", false, 0, "")
				pdf.Ln(20)
			})

			pdf.SetFooterFunc(func() {
				pdf.SetY(-15)
				pdf.SetFont("Arial", "I", 11)
				pdf.CellFormat(140, 10, fmt.Sprintf("%s", username), "0", 0, "L", false, 0, "")
				pdf.CellFormat(50, 10, fmt.Sprintf("%s", time.Now().Format("2006-01-02 15:04:05")), "0", 0, "R", false, 0, "")
			})

			pdf.AddPage()
			pdf.SetFont("Arial", "", 14)
			pdf.MultiCell(190, 4, "Merchant Log (Merchant to PGS):-", "0", "L", false)
			pdf.SetFont("Arial", "", 11)
			pdf.MultiCell(190, 4, merchantLogs, "0", "L", false)
			//pdf.CellFormat(40, 7, operatorLogs, "1", 0, "L", false, 0, "")
			//pdf.Ln(-1)
			pdf.AddPage()
			pdf.SetFont("Arial", "", 14)
			pdf.MultiCell(190, 4, "Operators/Banks Log (Banks to PGS):-", "0", "L", false)
			pdf.SetFont("Arial", "", 11)
			pdf.MultiCell(190, 4, operatorLogs, "0", "L", false)

			err := pdf.Output(c.Ctx.ResponseWriter)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("System is unable to process your request.Please contact customer care")
				return
			}
			return
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
	username = sess.Get("uname").(string)
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "DebugLogReport")
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
		c.Data["Header"] = template.HTML(headerContent)

		c.Data["Dashboard"] = beego.AppConfig.String("ENGLISH_DASHBOARD")
		c.Data["SearchFilters"] = beego.AppConfig.String("ENGLISH_SEARCH_FILTERS")
		c.Data["CNPSTransactionNumber"] = beego.AppConfig.String("ENGLISH_CNPS_TRANSACTION_NUMBER")
		c.Data["Bank"] = beego.AppConfig.String("ENGLISH_BANK")
		c.Data["Reset"] = beego.AppConfig.String("ENGLISH_RESET")
		c.Data["please_select"] = beego.AppConfig.String("ENGLISH_PLEASESELECTCOMMON")
		c.Data["enterTXN"] = beego.AppConfig.String("ENGLISH_ENTER_CNPS_TXN_ID")
		c.Data["GetLogs"] = beego.AppConfig.String("ENGLISH_GET_LOGS")
		c.Data["TxnLogDetals"] = beego.AppConfig.String("ENGLISH_TRANSACTION_LOGS")
		c.Data["Reporting"] = beego.AppConfig.String("ENGLISH_SYSTEM_REPORTS")
		c.Data["Back"] = beego.AppConfig.String("ENGLISH_BACK_BUTTON")

		c.TplName = "reports/transactionLogs/transactionLogs.html"
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
		c.Data["SearchFilters"] = beego.AppConfig.String("FRENCH_SEARCH_FILTERS")
		c.Data["CNPSTransactionNumber"] = beego.AppConfig.String("FRENCH_CNPS_TRANSACTION_NUMBER")
		c.Data["Bank"] = beego.AppConfig.String("FRENCH_BANK")
		c.Data["Reset"] = beego.AppConfig.String("FRENCH_RESET")
		c.Data["please_select"] = beego.AppConfig.String("FRENCH_PLEASESELECTCOMMON")
		c.Data["enterTXN"] = beego.AppConfig.String("FRENCH_ENTER_CNPS_TXN_ID")
		c.Data["TxnLogDetals"] = beego.AppConfig.String("FRENCH_TRANSACTION_LOGS")
		c.Data["Reporting"] = beego.AppConfig.String("FRENCH_SYSTEM_REPORTS")
		c.Data["GetLogs"] = beego.AppConfig.String("FRENCH_GET_LOGS")
		c.Data["Back"] = beego.AppConfig.String("FRENCH_BACK_BUTTON")

		c.TplName = "reports/transactionLogs/transactionLogs.html"
	}

	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "cnpstxnnumber - ", cnpstxnnumber)

	mcrow, mcerr := db.Db.Query(`select transaction_time, transaction_deatils, cnps_txn_number, bank_txn_number, operator, channel, amount, status, JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.TransactionProcess')), JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.BankTransactionDate')) from Transactions where cnps_txn_number=?`, cnpstxnnumber)
	if mcerr != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", mcerr)
		err = errors.New("Unable to get log data")
		return
	}
	defer sql.Close(mcrow)
	_, mcdata, mcserr := sql.Scan(mcrow)
	if mcserr != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", mcserr)
		err = errors.New("Unable to get log data")
		return
	}

	if len(mcdata) <= 0 {
		err = errors.New("Logs not found!.")
		return
	}

	err = json.Unmarshal([]byte(mcdata[0][1]), &req)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to parse data")
		return
	}
	merchantreq := MerchantRequest{}
	err = json.Unmarshal([]byte(mcdata[0][1]), &merchantreq)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to parse data")
		return
	}
	merchantreq.Sign = "XXXXXXXXXX"
	merchantreqjson, err := json.Marshal(merchantreq)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to parse data")
		return
	}
	merchantnotify := MerchantNotify{}
	err = json.Unmarshal([]byte(mcdata[0][1]), &merchantnotify)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to parse data")
		return
	}
	merchantnotify.BankReferenceID = req.BankReferenceID
	merchantnotify.Cnps_declaration_number = req.Cnps_declaration_number
	merchantnotify.Pgs_transaction_id = req.UUID
	merchantnotify.paymentMethod = req.Bank
	merchantnotify.paymentMode = req.Channel
	merchantnotify.Transaction_start_date = req.Transaction_start_date
	merchantnotify.BankTransactionDate = req.BankTransactionDate
	merchantnotifyjson, err := json.Marshal(merchantnotify)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to parse data")
		return
	}
	merchantLogs += "\n\n"
	merchantLogs += "[Info] : " + mcdata[0][0] + "[Request Date] MerchantRequest - [" + string(merchantreqjson) + "]\n\n\n"
	merchantLogs += "[Info] : " + mcdata[0][0] + "[Response Date] MerchantNotifyRes - [" + string(merchantnotifyjson) + "]\n\n\n"
	if mcdata[0][8] != "" {
		merchantLogs += "[Info] : " + mcdata[0][0] + "[Request Date] TransactionProcess - " + mcdata[0][8] + "\n\n\n"
	}

	row, err := db.Db.Query(`select request_date, response_date, response_dump, service, request_dump,operator from Switch_Transctions where cnps_txn_id =? OR pgs_txn_id=(SELECT pg_txn_number FROM Transactions WHERE cnps_txn_number=?)`, cnpstxnnumber, cnpstxnnumber)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))

	operatorLogs += "\n\n"
	if len(data) <= 0 {
		operatorLogs += "Banks/Operator not processed request yet!."
	}
	for _, logData := range data {
		if logData[3] == "GETBILLREQUEST" {
			getbillreq := GetBillReq{}
			err = json.Unmarshal([]byte(logData[4]), &getbillreq)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to parse data")
				return
			}
			getbillreq.Username = "XXXXXXXXXXX"
			getbillreq.Password = "XXXXXXXXXXX"
			getbillreqjson, err := json.Marshal(getbillreq)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to parse data")
				return
			}
			operatorLogs += "[Info] : " + logData[0] + "[Requested Date] " + logData[1] + "[Response Date] GetBillRequest -> " + string(getbillreqjson) + "\n\n\n"
			operatorLogs += "[Info] : " + logData[0] + "[Requested Date] " + logData[1] + "[Response Date] GetBillResponse -> " + logData[2] + "\n\n\n"
		} else if logData[3] == "PAYMENTNOTIFY" {

			paymentnotifyreq := PaymentNotify{}
			err = json.Unmarshal([]byte(logData[4]), &paymentnotifyreq)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to parse data")
				return
			}
			paymentnotifyreq.Username = "XXXXXXXXXXX"
			paymentnotifyreq.Password = "XXXXXXXXXXX"
			paymentnotifyreqjson, err := json.Marshal(paymentnotifyreq)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to parse data")
				return
			}
			operatorLogs += "[Info] : " + logData[0] + "[Requested Date] " + logData[1] + "[Response Date] PaymentNotifyRequest -> " + string(paymentnotifyreqjson) + "\n\n\n"
			operatorLogs += "[Info] : " + logData[0] + "[Requested Date] " + logData[1] + "[Response Date] PaymentNotifyResponse -> " + logData[2] + "\n\n\n"
		} else if logData[3] == "CNPSNOTIFY" {
			operatorLogs += "[Info] : " + logData[1] + "[Response Date] NotifyResponse -> " + logData[2] + "\n\n\n"
		} else if logData[3] == "BILLPAYMENTT" || logData[5] == "MTNBANKBILLMAP" {
			mtnpayreq := MTNPaymentReq{}

			mtnpayreq.Code = "XXXXXXXXXX"
			mtnpayreq.Password = "XXXXXXXXXX"
			mtnpayreq.MSISDN = req.MSISDN
			mtnpayreq.Reference = req.UUID
			mtnpayreq.Amount = req.Amount
			mtnpayreq.MetaData = req.Description
			mtnpayreqjson, err := json.Marshal(&mtnpayreq)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to parse data")
				return
			}
			operatorLogs += "[Info] : " + logData[1] + "[Response Date] Payment Request -> " + string(mtnpayreqjson) + "\n\n\n"
			mtnpayres := MTNPaymentRes{}
			mtnpayreserr := json.Unmarshal([]byte(logData[2]), &mtnpayres)
			if mtnpayreserr != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", mtnpayreserr)
				err = errors.New("Unable to parse data")
				return
			}
			mtnpayresjson, err := json.Marshal(&mtnpayres)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to parse data")
				return
			}
			operatorLogs += "[Info] : " + logData[1] + "[Response Date] Payment Response -> " + string(mtnpayresjson) + "\n\n\n"

		} else if logData[3] == "ORANGEIPN" || logData[3] == "BILLMAPUSSD" || logData[3] == "BILLMAP" {
			ipnres := IPNResponse{}
			operatorLogs += "[Info] : " + logData[1] + "[Response Date] IPN Request -> " + logData[4] + "\n\n\n"
			if json.Unmarshal([]byte(logData[2]), &ipnres) == nil {
				ipnresjson, err := json.Marshal(&ipnres)
				if err != nil {
					log.Println(beego.AppConfig.String("loglevel"), "Error", err)
					err = errors.New("Unable to parse data")
					return
				}
				operatorLogs += "[Info] : " + logData[1] + "[Response Date] IPN Response -> " + string(ipnresjson) + "\n\n\n"
			} else {
				operatorLogs += "[Info] : " + logData[1] + "[Response Date] IPN Response -> " + logData[2] + "\n\n\n"
			}
		} else if logData[3] == "WEBPAYMENT" {
			orangeres := OrangePaymentRes{}
			if json.Unmarshal([]byte(logData[2]), &orangeres) == nil {
				orangeresjson, err := json.Marshal(orangeres)
				if err != nil {
					log.Println(beego.AppConfig.String("loglevel"), "Error", err)
					err = errors.New("Unable to parse data")
					return
				}
				operatorLogs += "[Info] : " + logData[0] + "[Requested Date] " + logData[1] + "[Response Date] Payment Response -> " + string(orangeresjson) + "\n\n\n"

			} else {
				operatorLogs += "[Info] : " + logData[0] + "[Requested Date] " + logData[1] + "[Response Date] Payment Response -> " + logData[2] + "\n\n\n"
			}

		} /* else {
			//operatorLogs += "[Info] : " + logData[1] + "[Response Date] Final Response -> " + logData[2] + "\n\n\n"
		}*/
	}
	operatorLogs += "\n\n"

	return
}
