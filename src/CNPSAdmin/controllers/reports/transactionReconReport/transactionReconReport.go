package transactionReconReport

import (
	"CNPSAdmin/model/recondb"
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

type TransactionReconReport struct {
	beego.Controller
}

func (c *TransactionReconReport) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Recon Reports Page Start")
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
			c.TplName = "reports/transactionReconReport/transactionReconReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/transactionReconReport/transactionReconReport.html"
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "TransactionReconReport")
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
		c.Data["TransactionReconReport"] = beego.AppConfig.String("ENGLISH_TRANSACTION_REPORT")
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

		c.TplName = "reports/transactionReconReport/transactionReconReport.html"
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
		c.Data["TransactionReconReport"] = beego.AppConfig.String("FRENCH_TRANSACTION_REPORT")
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
		c.TplName = "reports/transactionReconReport/transactionReconReport.html"
	}

	return
}

func (c *TransactionReconReport) Post() {
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "TransactionReconReport")
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
	cnpstxnnumber := c.Input().Get("input_cnpstxnnumnber")
	pgtxnnumber := c.Input().Get("input_pgtxnnumnber")
	entity_id := c.Input().Get("input_entityid")
	entity_name := c.Input().Get("input_entityname")
	bank_name := c.Input().Get("input_bank")
	channel_name := c.Input().Get("input_channel")
	amount := c.Input().Get("input_amount")
	input_declaration_type := c.Input().Get("input_declaration_type")
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

	columns := [13]string{"cnps_txn_number", "JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Transaction_start_date'))", "pg_txn_number", "transaction_time", "bank_txn_number", "JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.BankTransactionDate'))", "entity_id", "entity_name", "JSON_UNQUOTE(JSON_EXTRACT(transaction_deatils, '$.Declaration_type'))", "operator", "channel", "amount", "status"}

	var sqlQuery bytes.Buffer
	var sqlArgs []interface{}
	sqlQuery.WriteString("SELECT t.cnps_txn_number as pgs_cnps_txn_number, JSON_UNQUOTE(JSON_EXTRACT(t.transaction_deatils, '$.Transaction_start_date')) as pgs_date, t.pg_txn_number as pgs_pg_txn_number, t.transaction_time as pgs_transaction_time, t.bank_txn_number as pgs_bank_txn_number, JSON_UNQUOTE(JSON_EXTRACT(t.transaction_deatils, '$.BankTransactionDate')) as pgs_bank_txn_date, t.entity_id as pgs_entity_id, t.entity_name as pgs_entity_name, JSON_UNQUOTE(JSON_EXTRACT(t.transaction_deatils, '$.Declaration_type')) as pgs_declaration_type, t.operator as pgs_operator, t.channel as pgs_channel, t.amount as pgs_amount, t.status as pgs_status,bt.bank_txn_date as bank_bank_txn_date,bt.cnps_txn_date,bt.cnps_txn_num,bt.bank_txn_num,bt.cnps_company_num,bt.declaration_type,bt.payment_method,bt.txn_type,bt.amount,bt.currency,bt.status,bt.description,bt.reserved1,bt.reserved2,btp.recon_status,btp.recon_description,btp.recon_date_time FROM Bank_Transactions_Pro as btp left join Transactions as t on t.id=btp.txn_id inner join Bank_Transactions as bt on bt.id = btp.bank_id where btp.id is not null")
	if cnpstxnnumber != "" {
		sqlQuery.WriteString(" AND bt.cnps_txn_num=?")
		sqlArgs = append(sqlArgs, cnpstxnnumber)
	}
	if pgtxnnumber != "" {
		sqlQuery.WriteString(" AND t.pg_txn_number=?")
		sqlArgs = append(sqlArgs, pgtxnnumber)
	}
	if entity_id != "" {
		sqlQuery.WriteString(" AND t.entity_id like ?")
		sqlArgs = append(sqlArgs, entity_id+"%")
	}
	if entity_name != "" {
		sqlQuery.WriteString(" AND bt.entity_name like ?")
		sqlArgs = append(sqlArgs, entity_name+"%")
	}
	if bank_name != "" {
		sqlQuery.WriteString(" AND t.operator=?")
		sqlArgs = append(sqlArgs, bank_name)
	}
	if channel_name != "" {
		sqlQuery.WriteString(" AND bt.payment_method=?")
		sqlArgs = append(sqlArgs, channel_name)
	}
	if amount != "" {
		sqlQuery.WriteString(" AND bt.amount like ?")
		sqlArgs = append(sqlArgs, amount+"%")
	}
	if status != "" {
		sqlQuery.WriteString(" AND t.status=?")
		sqlArgs = append(sqlArgs, status)
	}
	if input_declaration_type != "" && input_declaration_type != "ALL" {
		sqlQuery.WriteString(" AND bt.declaration_type=?")
		sqlArgs = append(sqlArgs, input_declaration_type)
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
		sqlQuery.WriteString(" AND (DATE (btp.recon_date_time) >= DATE(?)) AND (DATE (btp.recon_date_time) <= DATE(?))")
		log.Println(beego.AppConfig.String("loglevel"), "Debug", from, to)
		sqlArgs = append(sqlArgs, from, to)
	}
	if searchValue != "" {
		sqlQuery.WriteString(" AND (bt.cnps_txn_num like ? OR t.pg_txn_number like ? OR t.entity_id like ? OR t.entity_name like ? OR t.operator like ? OR bt.payment_method like ? OR bt.amount like ? OR t.status like ? OR bt.bank_txn_num like ?)")
		sqlArgs = append(sqlArgs, searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%")
	}
	if orderBy != "" {
		fmt.Fprintf(&sqlQuery, " ORDER BY %s %s", columns[orderByColumnNo], orderBy)
	}

	//Begin:: get total count befor limit condition
	totalRow, err := recondb.Db.Query(sqlQuery.String(), sqlArgs...)
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

	row, err := recondb.Db.Query(sqlQuery.String(), sqlArgs...)
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
		if sess.Get("role") == "ADMIN" && sess.Get("language") == "English" {
			if rowData[i][12] == "INITIATED" {
				statusbadge = "<span class='badge badge-pill badge-info'>" + rowData[i][12] + "</span>"
			} else if rowData[i][12] == "PENDING" {
				statusbadge = "<span class='badge badge-pill badge-warning'>" + rowData[i][12] + "</span>"
			} else if rowData[i][12] == "APPROVED" {
				statusbadge = "<span class='badge badge-pill badge-success'>" + rowData[i][12] + "</span>"
			} else if rowData[i][12] == "DECLINED" {
				statusbadge = "<span class='badge badge-pill badge-danger'>" + rowData[i][12] + "</span>"
			} else {
				statusbadge = "--"
			}
		} else if sess.Get("role") == "ADMIN" && sess.Get("language") == "French" {
			if rowData[i][12] == "INITIATED" {
				statusbadge = "<span class='badge badge-pill badge-info'>Initié</span>"
			} else if rowData[i][12] == "PENDING" {
				statusbadge = "<span class='badge badge-pill badge-warning'>En attente</span>"
			} else if rowData[i][12] == "APPROVED" {
				statusbadge = "<span class='badge badge-pill badge-success'>Approuvé</span>"
			} else if rowData[i][12] == "DECLINED" {
				statusbadge = "<span class='badge badge-pill badge-danger'>Décliné</span>"
			} else {
				statusbadge = "--"
			}
		}
		// r = append(r, rowData[i][0], rowData[i][1], rowData[i][2], rowData[i][3], rowData[i][4], rowData[i][5], rowData[i][6], rowData[i][7], rowData[i][8], rowData[i][9], rowData[i][10], utils.CurrencyFormat(rowData[i][11]), statusbadge, rowData[i][13], rowData[i][14], rowData[i][15], rowData[i][16], rowData[i][17], rowData[i][18], rowData[i][19], rowData[i][20], rowData[i][21], rowData[i][22], rowData[i][23], rowData[i][24], rowData[i][25], rowData[i][26], rowData[i][27], rowData[i][28], rowData[i][29])
		r = append(r, rowData[i][0], rowData[i][1], rowData[i][2], rowData[i][3], rowData[i][4], rowData[i][5], rowData[i][6], rowData[i][7], rowData[i][8], rowData[i][9], rowData[i][10], rowData[i][11], statusbadge, rowData[i][13], rowData[i][14], rowData[i][15], rowData[i][16], rowData[i][17], rowData[i][18], rowData[i][19], rowData[i][20], rowData[i][21], rowData[i][22], rowData[i][23], rowData[i][24], rowData[i][25], rowData[i][26], rowData[i][27], rowData[i][28], rowData[i][29])
		result = append(result, r)
		totalRecords += 1
	}

	var finalData = make(map[string]interface{})
	finalData["draw"] = draw
	finalData["recordsTotal"] = totalRecords
	finalData["recordsFiltered"] = len(rowTotalData)
	finalData["data"] = result
	c.Data["json"] = finalData
	c.ServeJSON()

	return
}
