package searchMerchant

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

type SearchMerchant struct {
	beego.Controller
}

func (c *SearchMerchant) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Start")
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
			c.TplName = "merchant/searchMerchant/searchMerchant.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "merchant/searchMerchant/searchMerchant.html"
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "SearchMerchant")
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
		c.Data["SearchMerchant"] = beego.AppConfig.String("ENGLISH_TRANSACTION_REPORT")
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

		c.Data["MerchantList"] = beego.AppConfig.String("ENGLISH_MERCHANT_LIST")
		c.Data["Merchant"] = beego.AppConfig.String("ENGLISH_MERCHANT")
		c.Data["MerchantApproval"] = beego.AppConfig.String("ENGLISH_MERCHANT_APPROVAL")
		c.Data["MerchantID"] = beego.AppConfig.String("ENGLISH_MERCHANT_ID")
		c.Data["Name"] = beego.AppConfig.String("ENGLISH_NAME")
		c.Data["MobileNumber"] = beego.AppConfig.String("ENGLISH_MOBILE_NUMBER")
		c.Data["EmailID"] = beego.AppConfig.String("ENGLISH_EMAIL_ID")
		c.Data["Country"] = beego.AppConfig.String("ENGLISH_COUNTRY")
		c.Data["RegistrationDate"] = beego.AppConfig.String("ENGLISH_REGESTRATION_DATE")
		c.Data["Action"] = beego.AppConfig.String("ENGLISH_ACTION")
		c.Data["Active"] = beego.AppConfig.String("ENGLISH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("ENGLISH_INACTIVE")
		c.Data["Approved"] = beego.AppConfig.String("ENGLISH_APPROVED")
		c.Data["TempBlock"] = beego.AppConfig.String("ENGLISH_TEMP_BLOCK")
		c.TplName = "merchant/searchMerchant/searchMerchant.html"
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
		c.Data["SearchMerchant"] = beego.AppConfig.String("FRENCH_TRANSACTION_REPORT")
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

		c.Data["MerchantList"] = beego.AppConfig.String("FRENCH_MERCHANT_LIST")
		c.Data["Merchant"] = beego.AppConfig.String("FRENCH_MERCHANT")
		c.Data["MerchantApproval"] = beego.AppConfig.String("FRENCH_MERCHANT_APPROVAL")
		c.Data["MerchantID"] = beego.AppConfig.String("FRENCH_MERCHANT_ID")
		c.Data["Name"] = beego.AppConfig.String("FRENCH_NAME")
		c.Data["MobileNumber"] = beego.AppConfig.String("FRENCH_MOBILE_NUMBER")
		c.Data["EmailID"] = beego.AppConfig.String("FRENCH_EMAIL_ID")
		c.Data["Country"] = beego.AppConfig.String("FRENCH_COUNTRY")
		c.Data["RegistrationDate"] = beego.AppConfig.String("FRENCH_REGESTRATION_DATE")
		c.Data["Action"] = beego.AppConfig.String("FRENCH_ACTION")
		c.Data["Active"] = beego.AppConfig.String("FRENCH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("FRENCH_INACTIVE")
		c.Data["Approved"] = beego.AppConfig.String("FRENCH_APPROVED")
		c.Data["TempBlock"] = beego.AppConfig.String("FRENCH_TEMP_BLOCK")
		c.TplName = "merchant/searchMerchant/searchMerchant.html"
	}

	return
}

func (c *SearchMerchant) Post() {
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "SearchMerchant")
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

	merchantId := c.Input().Get("input_merchantId")
	name := c.Input().Get("input_name")
	mobileNum := c.Input().Get("input_mobileNumber")
	emailId := c.Input().Get("input_emailId")
	country := c.Input().Get("input_country")
	status := c.Input().Get("input_status")
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

	columns := [13]string{"merchantId", "name", "mobileNum", "emailId", "country", "status", "created_at"}

	var sqlQuery bytes.Buffer
	var sqlArgs []interface{}
	sqlQuery.WriteString("select merchantId, name, mobileNum, emailId, country, status,created_at FROM Merchants where status IS NOT NULL ")
	if merchantId != "" {
		sqlQuery.WriteString(" AND merchantId=?")
		sqlArgs = append(sqlArgs, merchantId)
	}
	if name != "" {
		sqlQuery.WriteString(" AND name like ?")
		sqlArgs = append(sqlArgs, name+"%")
	}
	if mobileNum != "" {
		sqlQuery.WriteString(" AND mobileNum like ?")
		sqlArgs = append(sqlArgs, mobileNum+"%")
	}
	if emailId != "" {
		sqlQuery.WriteString(" AND emailId like ?")
		sqlArgs = append(sqlArgs, emailId+"%")
	}
	if country != "" {
		sqlQuery.WriteString(" AND country like ?")
		sqlArgs = append(sqlArgs, country+"%")
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
		sqlQuery.WriteString(" AND (merchantId like ? OR name like ? OR mobileNum like ? OR emailId like ? OR country like ? OR status like ?)")
		sqlArgs = append(sqlArgs, searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%")
	}
	if orderBy != "" {
		fmt.Fprintf(&sqlQuery, " ORDER BY %s %s", columns[orderByColumnNo], orderBy)
	}

	//Begin:: get total count befor limit condition
	// a := sqlQuery.String()
	// fmt.Println(a)
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
		merchantStatus := rowData[i][5]
		if sess.Get("role") == "ADMIN" && sess.Get("language") == "English" {
			if merchantStatus == "PENDING" {
				statusbadge = "<span class='badge badge-pill badge-primary'>" + merchantStatus + "</span>"
			} else if merchantStatus == "APPROVED" {
				statusbadge = "<span class='badge badge-pill badge-secondary'>" + merchantStatus + "</span>"
			} else if merchantStatus == "ACTIVE" {
				statusbadge = "<span class='badge badge-pill badge-success'>" + merchantStatus + "</span>"
			} else if merchantStatus == "TEMP_BLOCK" {
				statusbadge = "<span class='badge badge-pill badge-warning'>" + merchantStatus + "</span>"
			} else if merchantStatus == "INACTIVE" {
				statusbadge = "<span class='badge badge-pill badge-danger'>" + merchantStatus + "</span>"
			} else {
				statusbadge = "--"
			}
			updateLink := "<a href='/Merchant/UpdateMerchant/" + rowData[i][0] + "'><h6 class='text-red' data-toggle='tooltip' data-placement='top' data-original-title='" + beego.AppConfig.String("ENGLISH_UPDATE_BUTTON") + "'>" + beego.AppConfig.String("ENGLISH_UPDATE_BUTTON") + "</h6></a>"
			r = append(r, rowData[i][0], rowData[i][1], rowData[i][2], rowData[i][3], rowData[i][4], statusbadge, rowData[i][6], updateLink)
			result = append(result, r)
		} else if sess.Get("role") == "ADMIN" && sess.Get("language") == "French" {
			if merchantStatus == "PENDING" {
				statusbadge = "<span class='badge badge-pill badge-primary'>En attente</span>"
			} else if merchantStatus == "APPROVED" {
				statusbadge = "<span class='badge badge-pill badge-secondary'>Approuvé</span>"
			} else if merchantStatus == "ACTIVE" {
				statusbadge = "<span class='badge badge-pill badge-success'>ACTIVE</span>"
			} else if merchantStatus == "TEMP_BLOCK" {
				statusbadge = "<span class='badge badge-pill badge-warning'>blocage temporaire</span>"
			} else if merchantStatus == "INACTIVE" {
				statusbadge = "<span class='badge badge-pill badge-danger'>INACTIVE</span>"
			} else {
				statusbadge = "--"
			}
			updateLink := "<a href='/Merchant/UpdateMerchant/" + rowData[i][0] + "'><h6 class='text-red' data-toggle='tooltip' data-placement='top' data-original-title='" + beego.AppConfig.String("FRENCH_UPDATE_BUTTON") + "'>" + beego.AppConfig.String("FRENCH_UPDATE_BUTTON") + "</h6></a>"
			r = append(r, rowData[i][0], rowData[i][1], rowData[i][2], rowData[i][3], rowData[i][4], statusbadge, rowData[i][6], updateLink)
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
