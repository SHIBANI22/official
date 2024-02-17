package searchRefund

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"bytes"
	"fmt"
	"strconv"
	"time"

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

type SearchRefund struct {
	beego.Controller
}

func (c *SearchRefund) Get() {
	resp := c.Input().Get("resp")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchRefund Page Start", resp)
	var err error
	sessErr := false
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
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)

			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "processing/refundProcessing/searchRefund/searchRefund.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchRefund Page Fail")
		} else {
			c.Data["DisplayMessage"] = resp
			c.TplName = "processing/refundProcessing/searchRefund/searchRefund.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchRefund Page Success")
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "SearchRefund")
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
		c.Data["SearchRefund"] = beego.AppConfig.String("ENGLISH_TRANSACTION_REPORT")
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
		c.Data["SearchRefund"] = beego.AppConfig.String("ENGLISH_REPORTS")
		c.Data["enterTXN"] = beego.AppConfig.String("ENGLISH_ENTER_CNPS_TXN_ID")
		c.Data["input_pgtxnnumnber"] = beego.AppConfig.String("ENGLISH_ENTER_PG_TXN_ID")
		c.Data["input_entityid"] = beego.AppConfig.String("ENGLISH_ENTER_ENTITY_ID")
		c.Data["input_entityname"] = beego.AppConfig.String("ENGLISH_ENTER_ENTITY_NAME")
		c.Data["input_amount"] = beego.AppConfig.String("ENGLISH_ENTER_AMOUNT")

		c.Data["GetRecords"] = beego.AppConfig.String("ENGLISH_GET_RECORDS")
		c.Data["TxnDetals"] = beego.AppConfig.String("ENGLISH_TRANSACTION_REPORTS_MENU")
		c.Data["List"] = beego.AppConfig.String("ENGLISH_LIST")
		c.Data["Initiated"] = beego.AppConfig.String("ENGLISH_INITIATED")
		c.Data["Reporting"] = beego.AppConfig.String("ENGLISH_REPORTING_MENU")
		c.Data["DeclarationTypeLabel"] = beego.AppConfig.String("ENGLISH_DECLARATION_TYPE")
		c.Data["DeclarationTypeAllLabel"] = beego.AppConfig.String("ENGLISH_ALL")

		c.Data["RefundID"] = beego.AppConfig.String("ENGLISH_REFUND_ID")
		c.Data["BankTxnID"] = beego.AppConfig.String("ENGLISH_BANK_TXN_ID")
		c.Data["RefundInitiatedDate"] = beego.AppConfig.String("ENGLISH_REFUND_INITIATED_DATE")
		c.Data["Success"] = beego.AppConfig.String("ENGLISH_SUCCESS")
		c.Data["Failed"] = beego.AppConfig.String("ENGLISH_FAILED")
		c.Data["RefundDate"] = beego.AppConfig.String("ENGLISH_REFUND_DATE")
		c.Data["Action"] = beego.AppConfig.String("ENGLISH_ACTION")
		c.Data["RefundProcessing"] = beego.AppConfig.String("ENGLISH_REFUND_PROCESSING")
		c.Data["ManualTransaction"] = beego.AppConfig.String("ENGLISH_MANUAL_TRANSACTION")

		c.TplName = "processing/refundProcessing/searchRefund/searchRefund.html"
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
		c.Data["SearchRefund"] = beego.AppConfig.String("FRENCH_TRANSACTION_REPORT")
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
		c.Data["SearchRefund"] = beego.AppConfig.String("FRENCH_REPORTS")
		c.Data["enterTXN"] = beego.AppConfig.String("FRENCH_ENTER_CNPS_TXN_ID")
		c.Data["input_pgtxnnumnber"] = beego.AppConfig.String("FRENCH_ENTER_PG_TXN_ID")
		c.Data["input_entityid"] = beego.AppConfig.String("FRENCH_ENTER_ENTITY_ID")
		c.Data["input_entityname"] = beego.AppConfig.String("FRENCH_ENTER_ENTITY_NAME")
		c.Data["input_amount"] = beego.AppConfig.String("FRENCH_ENTER_AMOUNT")
		c.Data["TxnDetals"] = beego.AppConfig.String("FRENCH_TRANSACTION_REPORTS_MENU")
		c.Data["Reporting"] = beego.AppConfig.String("FRENCH_REPORTING_MENU")

		c.Data["GetRecords"] = beego.AppConfig.String("FRENCH_GET_RECORDS")
		c.Data["List"] = beego.AppConfig.String("FRENCH_LIST")
		c.Data["Initiated"] = beego.AppConfig.String("FRENCH_INITIATED")
		c.Data["DeclarationTypeLabel"] = beego.AppConfig.String("FRENCH_DECLARATION_TYPE")
		c.Data["DeclarationTypeAllLabel"] = beego.AppConfig.String("FRENCH_ALL")

		c.Data["RefundID"] = beego.AppConfig.String("FRENCH_REFUND_ID")
		c.Data["BankTxnID"] = beego.AppConfig.String("FRENCH_BANK_TXN_ID")
		c.Data["RefundInitiatedDate"] = beego.AppConfig.String("FRENCH_REFUND_INITIATED_DATE")
		c.Data["Success"] = beego.AppConfig.String("FRENCH_SUCCESS")
		c.Data["Failed"] = beego.AppConfig.String("FRENCH_FAILED")
		c.Data["RefundDate"] = beego.AppConfig.String("FRENCH_REFUND_DATE")
		c.Data["Action"] = beego.AppConfig.String("FRENCH_ACTION")
		c.Data["RefundProcessing"] = beego.AppConfig.String("FRENCH_REFUND_PROCESSING")
		c.Data["ManualTransaction"] = beego.AppConfig.String("FRENCH_MANUAL_TRANSACTION")
		c.TplName = "processing/refundProcessing/searchRefund/searchRefund.html"
	}

	return
}

func (c *SearchRefund) Post() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Start")
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "SearchRefund")
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

	status := c.Input().Get("input_status")
	merchantId := c.Input().Get("input_merchantId")
	refundId := c.Input().Get("input_refundId")
	bankTxnId := c.Input().Get("input_bankTxnId")
	amount := c.Input().Get("input_amount")
	daterange := strings.Split(c.Input().Get("daterange"), " - ")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", c.Input().Get("daterange"))

	start := strings.TrimSpace(c.Input().Get("start"))
	length := strings.TrimSpace(c.Input().Get("length"))
	draw := strings.TrimSpace(c.Input().Get("draw"))
	searchValue := strings.TrimSpace(c.Input().Get("search[value]"))

	orderBy := strings.TrimSpace(c.Input().Get("order_by"))
	orderByColumn := strings.TrimSpace(c.Input().Get("order_by_column"))
	if c.Input().Get("order[0][column]") != "" {
		orderByColumn = strings.TrimSpace(c.Input().Get("order[0][column]"))
	}
	if c.Input().Get("order[0][dir]") != "" {
		orderBy = strings.TrimSpace(c.Input().Get("order[0][dir]"))
	}
	orderByColumnNo, err := strconv.Atoi(orderByColumn)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	}

	columns := [8]string{"merchant_id", "refund_id", "bank_txn_id", "amount", "created_at", "updated_at", "status"}

	var sqlQuery bytes.Buffer
	var sqlArgs []interface{}
	sqlQuery.WriteString("select merchant_id, refund_id, bank_txn_id, amount, created_at, updated_at, status FROM Refund WHERE status='INITIATED'")
	if merchantId != "" {
		sqlQuery.WriteString(" AND merchant_id like ?")
		sqlArgs = append(sqlArgs, merchantId+"%")
	}
	if refundId != "" {
		sqlQuery.WriteString(" AND refund_id like ?")
		sqlArgs = append(sqlArgs, refundId)
	}
	if bankTxnId != "" {
		sqlQuery.WriteString(" AND bank_txn_id like ?")
		sqlArgs = append(sqlArgs, bankTxnId+"%")
	}
	if amount != "" {
		sqlQuery.WriteString(" AND amount like ?")
		sqlArgs = append(sqlArgs, amount+"%")
	}
	if status != "" {
		sqlQuery.WriteString(" AND status=?")
		sqlArgs = append(sqlArgs, status)
	}

	if c.Input().Get("daterange") != "" && len(daterange) == 2 {
		format, dateErr := utils.CustomDateFormat(beego.AppConfig.String("DateFormat"), "")
		if dateErr != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", dateErr)
		}

		inputFromDate, _ := time.Parse(format, daterange[0])
		inputToDate, _ := time.Parse(format, daterange[1])
		from := inputFromDate.Format("2006-01-02")
		to := inputToDate.Format("2006-01-02")
		sqlQuery.WriteString(" AND (DATE (created_at) >= DATE(?)) AND (DATE (created_at) <= DATE(?))")
		log.Println(beego.AppConfig.String("loglevel"), "Debug", from, to)
		sqlArgs = append(sqlArgs, from, to)
	}
	if searchValue != "" {
		sqlQuery.WriteString(" AND (merchant_id like ? OR refund_id like ? OR bank_txn_id like ? OR amount like ? OR status like ?)")
		sqlArgs = append(sqlArgs, searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%")
	}
	if orderBy != "" {
		fmt.Fprintf(&sqlQuery, " ORDER BY %s %s", columns[orderByColumnNo], orderBy)
	}

	//Begin:: get total count befor limit condition
	totalRow, err := db.Db.Query(sqlQuery.String(), sqlArgs...)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	}
	//End:: get total count befor limit condition
	lengthInt, _ := strconv.Atoi(length)
	if lengthInt > 0 {
		sqlQuery.WriteString(" LIMIT ?,?")
		sqlArgs = append(sqlArgs, start, length)
	}
	log.Println("Debug", "Info", "Final Query", sqlQuery.String(), sqlArgs)

	row, err := db.Db.Query(sqlQuery.String(), sqlArgs...)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	}

	defer sql.Close(row)
	defer sql.Close(totalRow)
	_, rowData, err := sql.Scan(row)
	_, rowTotalData, err := sql.Scan(totalRow)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	}
	var result = make([]interface{}, 0)
	var totalRecords int = 0
	for i := range rowData {
		var r []string
		var statusbadge string
		status := rowData[i][6]
		if sess.Get("role") == "ADMIN" && sess.Get("language") == "English" {
			if status == "INITIATED" {
				statusbadge = "<span class='badge badge-pill badge-info'>" + status + "</span>"
			} else if status == "SUCCESS" {
				statusbadge = "<span class='badge badge-pill badge-success'>" + status + "</span>"
			} else if status == "FAILED" {
				statusbadge = "<span class='badge badge-pill badge-danger'>" + status + "</span>"
			} else {
				statusbadge = "--"
			}
			approveLink := "<a href='/ApproveRefund/" + rowData[i][1] + "'><h6 class='text-red' data-toggle='tooltip' data-placement='top' data-original-title='" + beego.AppConfig.String("ENGLISH_REFUND_BUTTON") + "'>" + beego.AppConfig.String("ENGLISH_REFUND_BUTTON") + "</h6></a>"
			r = append(r, rowData[i][0], rowData[i][1], rowData[i][2], rowData[i][3], rowData[i][4], rowData[i][5], statusbadge, approveLink)
			result = append(result, r)
		} else if sess.Get("role") == "ADMIN" && sess.Get("language") == "French" {
			if status == "INITIATED" {
				statusbadge = "<span class='badge badge-pill badge-info'>Initié</span>"
			} else if status == "SUCCESS" {
				statusbadge = "<span class='badge badge-pill badge-success'>SUCCÈS</span>"
			} else if status == "FAILED" {
				statusbadge = "<span class='badge badge-pill badge-danger'>ÉCHOUÉ</span>"
			} else {
				statusbadge = "--"
			}

			approveLink := "<a href='/ApproveRefund/" + rowData[i][1] + "'><h6 class='text-red' data-toggle='tooltip' data-placement='top' data-original-title='" + beego.AppConfig.String("FRENCH_REFUND_BUTTON") + "'>" + beego.AppConfig.String("FRENCH_REFUND_BUTTON") + "</h6></a>"
			r = append(r, rowData[i][0], rowData[i][1], rowData[i][2], rowData[i][3], rowData[i][4], rowData[i][5], statusbadge, approveLink)
			result = append(result, r)

		}

		totalRecords += 1
	}

	var finalData = make(map[string]interface{})
	finalData["draw"] = draw
	finalData["recordsTotal"] = totalRecords
	finalData["recordsFiltered"] = len(rowTotalData)
	finalData["data"] = result
	c.Data["json"] = finalData
	c.ServeJSON()
}

func (c *SearchRefund) ApproveRefund() {
	refundId := c.Ctx.Input.Param(":RefundID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "ApproveRefund Page Start")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "refundId - ", refundId)
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
				c.Redirect(beego.AppConfig.String("PGS_SEARCH_REFUND")+"?resp="+err.Error(), 307)
			}
		}
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "SearchRefund")
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

	user_email := sess.Get("uname")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_email - ", user_email)

	bankTxnId := callApprove()

	result, err := db.Db.Exec(`UPDATE Refund SET bank_txn_id=?,status=?,updated_by=?,updated_at=now() WHERE refund_id=?`,
		bankTxnId, "SUCCESS", user_email, refundId)
	if err != nil {
		err = errors.New("Approve Refund updation failed")
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Approve Refund error")
		return
	}

	i, rowerr := result.RowsAffected()
	if rowerr != nil || i == 0 {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Approve Refund error", err, i)
		err = errors.New("Refund updation failed")
		return
	}

	c.Redirect(beego.AppConfig.String("PGS_SEARCH_REFUND")+"?resp=Refund Success", 302)
}

// TODO : Implement Refund API call
func callApprove() (bankTxnId string) {
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Refund API not Implemented - ")
	return "1231"
}
