package channelReport

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
	TimeStamp string
	Adminid   string
	Status    string
	Request   string
	Response  string
	Reqtype   string
}

type ChannelReport struct {
	beego.Controller
}

func (c *ChannelReport) Get() {
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
			c.TplName = "reports/channelReport/channelReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/channelReport/channelReport.html"
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "ChannelReport")
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
		c.Data["ChannelReport"] = beego.AppConfig.String("ENGLISH_CHANNEL_REPORT")
		c.Data["SearchFilters"] = beego.AppConfig.String("ENGLISH_SEARCH_FILTERS")
		c.Data["SearchResults"] = beego.AppConfig.String("ENGLISH_USER_SEARCH_RESULTS")
		c.Data["SelectDateRange"] = beego.AppConfig.String("ENGLISH_USER_SELECT_DATARANGE")
		c.Data["TimeStamp"] = beego.AppConfig.String("ENGLISH_TIMESTAMP")
		c.Data["MerchantID"] = beego.AppConfig.String("ENGLISH_MERCHANT_ID")
		c.Data["Status"] = beego.AppConfig.String("ENGLISH_STATUS")
		c.Data["Response"] = beego.AppConfig.String("ENGLISH_RESPONSE")
		c.Data["Request"] = beego.AppConfig.String("ENGLISH_REQUEST")
		c.Data["ListOfChannelReports"] = beego.AppConfig.String("ENGLISH_LIST_OF_CHANNEL_REPORT")
		c.Data["Search"] = beego.AppConfig.String("ENGLISH_SEARCH")
		c.Data["Reset"] = beego.AppConfig.String("ENGLISH_RESET")
		c.Data["Reports"] = beego.AppConfig.String("ENGLISH_REPORTS")
		c.Data["admin_id"] = beego.AppConfig.String("ENGLISH_ENTER_USER_NAME")
		c.Data["please_select"] = beego.AppConfig.String("ENGLISH_PLEASESELECTCOMMON")
		c.Data["Approved"] = beego.AppConfig.String("ENGLISH_APPROVED")
		c.Data["Declined"] = beego.AppConfig.String("ENGLISH_DECLINED")
		c.Data["requestType"] = beego.AppConfig.String("ENGLISH_REQUEST_TYPE")
		c.Data["requestPayment"] = beego.AppConfig.String("ENGLISH_REQUEST_PAYMENT")
		c.Data["requestGetstatus"] = beego.AppConfig.String("ENGLISH_REQUEST_GETSTATUS")

		c.Data["GetRecords"] = beego.AppConfig.String("ENGLISH_GET_RECORDS")
		c.Data["List"] = beego.AppConfig.String("ENGLISH_LIST")

		c.TplName = "reports/channelReport/channelReport.html"
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
		c.Data["ChannelReport"] = beego.AppConfig.String("FRENCH_CHANNEL_REPORT")
		c.Data["SearchFilters"] = beego.AppConfig.String("FRENCH_SEARCH_FILTERS")
		c.Data["SearchResults"] = beego.AppConfig.String("FRENCH_USER_SEARCH_RESULTS")
		c.Data["SelectDateRange"] = beego.AppConfig.String("FRENCH_USER_SELECT_DATARANGE")
		c.Data["TimeStamp"] = beego.AppConfig.String("FRENCH_TIMESTAMP")
		c.Data["MerchantID"] = beego.AppConfig.String("FRENCH_MERCHANT_ID")
		c.Data["Response"] = beego.AppConfig.String("FRENCH_RESPONSE")
		c.Data["Status"] = beego.AppConfig.String("FRENCH_STATUS")
		c.Data["Request"] = beego.AppConfig.String("FRENCH_REQUEST")
		c.Data["ListOfChannelReports"] = beego.AppConfig.String("FRENCH_LIST_OF_CHANNEL_REPORT")
		c.Data["Search"] = beego.AppConfig.String("FRENCH_SEARCH")
		c.Data["Reset"] = beego.AppConfig.String("FRENCH_RESET")
		c.Data["Reports"] = beego.AppConfig.String("FRENCH_REPORTS")
		c.Data["admin_id"] = beego.AppConfig.String("FRENCH_ENTER_USER_NAME")
		c.Data["please_select"] = beego.AppConfig.String("FRENCH_PLEASESELECTCOMMON")
		c.Data["Approved"] = beego.AppConfig.String("FRENCH_APPROVED")
		c.Data["Declined"] = beego.AppConfig.String("FRENCH_DECLINED")
		c.Data["requestType"] = beego.AppConfig.String("FRENCH_REQUEST_TYPE")
		c.Data["requestPayment"] = beego.AppConfig.String("FRENCH_REQUEST_PAYMENT")
		c.Data["requestGetstatus"] = beego.AppConfig.String("FRENCH_REQUEST_GETSTATUS")

		c.Data["GetRecords"] = beego.AppConfig.String("FRENCH_GET_RECORDS")
		c.Data["List"] = beego.AppConfig.String("FRENCH_LIST")

		c.TplName = "reports/channelReport/channelReport.html"
	}

	return
}

func (c *ChannelReport) Post() {
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "ChannelReport")
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

	entityId := c.Input().Get("input_entity_id")
	status := c.Input().Get("input_status")
	reqtype := c.Input().Get("input_request_type")
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

	columns := [6]string{"entity_id", "dump_time", "request", "response", "reqtype", "JSON_UNQUOTE(JSON_EXTRACT(response, '$.Status'))"}

	var sqlQuery bytes.Buffer
	var sqlArgs []interface{}
	sqlQuery.WriteString("select entity_id, dump_time, request, response,reqtype, JSON_UNQUOTE(JSON_EXTRACT(response, '$.Status')) FROM Events WHERE entity_id IS NOT NULL")
	if entityId != "" {
		sqlQuery.WriteString(" AND entity_id like ?")
		sqlArgs = append(sqlArgs, entityId+"%")
	}
	if reqtype != "" {
		sqlQuery.WriteString(" AND reqtype=?")
		sqlArgs = append(sqlArgs, reqtype)
	}

	if status != "" {
		sqlQuery.WriteString(" AND JSON_UNQUOTE(JSON_EXTRACT(response, '$.Status'))=?")
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
		sqlQuery.WriteString(" AND (DATE (dump_time) >= DATE(?)) AND (DATE (dump_time) <= DATE(?))")
		log.Println(beego.AppConfig.String("loglevel"), "Debug", from, to)
		sqlArgs = append(sqlArgs, from, to)
	}
	if searchValue != "" {
		sqlQuery.WriteString(" AND (entity_id like ? OR JSON_UNQUOTE(JSON_EXTRACT(response, '$.Status')) like ? OR reqtype like ?)")
		sqlArgs = append(sqlArgs, searchValue+"%", searchValue+"%", searchValue+"%")
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
		if sess.Get("role") == "ADMIN" && sess.Get("language") == "English" {
			if rowData[i][5] == "APPROVED" {
				statusbadge = "<span class='badge badge-pill badge-success'>" + rowData[i][5] + "</span>"
			} else if rowData[i][5] == "DECLINED" {
				statusbadge = "<span class='badge badge-pill badge-danger'>" + rowData[i][5] + "</span>"
			} else {
				statusbadge = "--"
			}
			r = append(r, rowData[i][0], rowData[i][1], rowData[i][2], rowData[i][3], rowData[i][4], statusbadge)
			result = append(result, r)
		} else if sess.Get("role") == "ADMIN" && sess.Get("language") == "French" {
			if rowData[i][5] == "APPROVED" {
				statusbadge = "<span class='badge badge-pill badge-success'>Approuvé</span>"
			} else if rowData[i][5] == "DECLINED" {
				statusbadge = "<span class='badge badge-pill badge-danger'>Décliné</span>"
			} else {
				statusbadge = "--"
			}
			r = append(r, rowData[i][0], rowData[i][1], rowData[i][2], rowData[i][3], rowData[i][4], statusbadge)
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

	return
}
