package adminAduitReport

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
	TimeStamp   string
	Adminid     string
	URL         string
	Status      string
	IP          string
	Host        string
	HTTPSMethod string
	SessionID   string
}

type AdminAduitReport struct {
	beego.Controller
}

func (c *AdminAduitReport) Get() {
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
			c.TplName = "reports/adminAduitReport/adminAduitReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/adminAduitReport/adminAduitReport.html"
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "AuditReport")
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
		c.Data["AdminAduitReport"] = beego.AppConfig.String("ENGLISH_ADUIT_REPORT")
		c.Data["SearchFilters"] = beego.AppConfig.String("ENGLISH_SEARCH_FILTERS")
		c.Data["SearchResults"] = beego.AppConfig.String("ENGLISH_USER_SEARCH_RESULTS")
		c.Data["SelectDateRange"] = beego.AppConfig.String("ENGLISH_USER_SELECT_DATARANGE")
		c.Data["AdminUserID"] = beego.AppConfig.String("ENGLISH_ADMIN_USER_ID")
		c.Data["TimeStamp"] = beego.AppConfig.String("ENGLISH_TIMESTAMP")
		c.Data["AdminID"] = beego.AppConfig.String("ENGLISH_ADMIN_ID")
		c.Data["Status"] = beego.AppConfig.String("ENGLISH_STATUS")
		c.Data["URL"] = beego.AppConfig.String("ENGLISH_URL")
		c.Data["IP"] = beego.AppConfig.String("ENGLISH_IP")
		c.Data["Host"] = beego.AppConfig.String("ENGLISH_HOST")
		c.Data["HTTPSMethod"] = beego.AppConfig.String("ENGLISH_HTTPS_METHOD")
		c.Data["SessionID"] = beego.AppConfig.String("ENGLISH_SESSION_ID")
		c.Data["ListOfAuditReports"] = beego.AppConfig.String("ENGLISH_LIST_OF_ADUIT_REPORT")
		c.Data["Search"] = beego.AppConfig.String("ENGLISH_SEARCH")
		c.Data["Reset"] = beego.AppConfig.String("ENGLISH_RESET")
		c.Data["Reports"] = beego.AppConfig.String("ENGLISH_REPORTS")
		c.Data["admin_id"] = beego.AppConfig.String("ENGLISH_ENTER_USER_NAME")
		c.Data["input_url"] = beego.AppConfig.String("ENGLISH_ENTER_PAGE_ACCESS")
		c.Data["input_ip"] = beego.AppConfig.String("ENGLISH_ENTER_IP_ADDRESS")
		c.Data["input_host"] = beego.AppConfig.String("ENGLISH_ENTER_MACHINE_NAME")
		c.Data["please_select"] = beego.AppConfig.String("ENGLISH_PLEASESELECTCOMMON")
		c.Data["Success"] = beego.AppConfig.String("ENGLISH_SUCCESS")
		c.Data["Failure"] = beego.AppConfig.String("ENGLISH_FAILURE")

		c.Data["GetRecords"] = beego.AppConfig.String("ENGLISH_GET_RECORDS")
		c.Data["List"] = beego.AppConfig.String("ENGLISH_LIST")

		c.TplName = "reports/adminAduitReport/adminAduitReport.html"
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
		c.Data["AdminAduitReport"] = beego.AppConfig.String("FRENCH_ADUIT_REPORT")
		c.Data["SearchFilters"] = beego.AppConfig.String("FRENCH_SEARCH_FILTERS")
		c.Data["SearchResults"] = beego.AppConfig.String("FRENCH_USER_SEARCH_RESULTS")
		c.Data["SelectDateRange"] = beego.AppConfig.String("FRENCH_USER_SELECT_DATARANGE")
		c.Data["AdminUserID"] = beego.AppConfig.String("FRENCH_ADMIN_USER_ID")
		c.Data["TimeStamp"] = beego.AppConfig.String("FRENCH_TIMESTAMP")
		c.Data["AdminID"] = beego.AppConfig.String("FRENCH_ADMIN_ID")
		c.Data["URL"] = beego.AppConfig.String("FRENCH_URL")
		c.Data["Status"] = beego.AppConfig.String("FRENCH_STATUS")
		c.Data["IP"] = beego.AppConfig.String("FRENCH_IP")
		c.Data["Host"] = beego.AppConfig.String("FRENCH_HOST")
		c.Data["HTTPSMethod"] = beego.AppConfig.String("FRENCH_HTTPS_METHOD")
		c.Data["SessionID"] = beego.AppConfig.String("FRENCH_SESSION_ID")
		c.Data["ListOfAuditReports"] = beego.AppConfig.String("FRENCH_LIST_OF_ADUIT_REPORT")
		c.Data["Search"] = beego.AppConfig.String("FRENCH_SEARCH")
		c.Data["Reset"] = beego.AppConfig.String("FRENCH_RESET")
		c.Data["Reports"] = beego.AppConfig.String("FRENCH_REPORTS")
		c.Data["admin_id"] = beego.AppConfig.String("FRENCH_ENTER_USER_NAME")
		c.Data["input_url"] = beego.AppConfig.String("FRENCH_ENTER_PAGE_ACCESS")
		c.Data["input_ip"] = beego.AppConfig.String("FRENCH_ENTER_IP_ADDRESS")
		c.Data["input_host"] = beego.AppConfig.String("FRENCH_ENTER_MACHINE_NAME")
		c.Data["please_select"] = beego.AppConfig.String("FRENCH_PLEASESELECTCOMMON")
		c.Data["Success"] = beego.AppConfig.String("FRENCH_SUCCESS")
		c.Data["Failure"] = beego.AppConfig.String("FRENCH_FAILURE")

		c.Data["GetRecords"] = beego.AppConfig.String("FRENCH_GET_RECORDS")
		c.Data["List"] = beego.AppConfig.String("FRENCH_LIST")

		c.TplName = "reports/adminAduitReport/adminAduitReport.html"
	}

	return
}

func (c *AdminAduitReport) Post() {
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "AuditReport")
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

	admin_id := c.Input().Get("admin_id")
	url := c.Input().Get("input_url")
	ipp := c.Input().Get("input_ip")
	status := c.Input().Get("input_status")
	host := c.Input().Get("input_host")
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

	columns := [8]string{"user_id", "created_on", "url", "ip", "host", "JSON_UNQUOTE(JSON_EXTRACT(event, '$.Method'))", "JSON_UNQUOTE(JSON_EXTRACT(event, '$.SessionID'))", "status"}

	var sqlQuery bytes.Buffer
	var sqlArgs []interface{}
	sqlQuery.WriteString("select JSON_UNQUOTE(JSON_EXTRACT(event, '$.UserName')), created_on, JSON_UNQUOTE(JSON_EXTRACT(event, '$.URL')), JSON_UNQUOTE(JSON_EXTRACT(event, '$.PIP')), JSON_UNQUOTE(JSON_EXTRACT(event, '$.Host')), JSON_UNQUOTE(JSON_EXTRACT(event, '$.Method')), JSON_UNQUOTE(JSON_EXTRACT(event, '$.SessionID')), JSON_UNQUOTE(JSON_EXTRACT(event, '$.Status')) FROM web_event WHERE user_id IS NOT NULL")
	if admin_id != "" {
		sqlQuery.WriteString(" AND user_id like ?")
		sqlArgs = append(sqlArgs, admin_id+"%")
	}
	if url != "" {
		sqlQuery.WriteString(" AND url like ?")
		sqlArgs = append(sqlArgs, url+"%")
	}
	if ipp != "" {
		sqlQuery.WriteString(" AND ip like ?")
		sqlArgs = append(sqlArgs, ipp+"%")
	}
	if host != "" {
		sqlQuery.WriteString(" AND host like ?")
		sqlArgs = append(sqlArgs, host)
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
		sqlQuery.WriteString(" AND (DATE (created_on) >= DATE(?)) AND (DATE (created_on) <= DATE(?))")
		log.Println(beego.AppConfig.String("loglevel"), "Debug", from, to)
		sqlArgs = append(sqlArgs, from, to)
	}
	if searchValue != "" {
		sqlQuery.WriteString(" AND (user_id like ? OR url like ? OR ip like ? OR host like ? OR status like ? OR JSON_UNQUOTE(JSON_EXTRACT(event, '$.SessionID')) like ? OR JSON_UNQUOTE(JSON_EXTRACT(event, '$.Method')) like ?)")
		sqlArgs = append(sqlArgs, searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%")
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
			if rowData[i][7] == "Success" {
				statusbadge = "<span class='badge badge-pill badge-success'>" + rowData[i][7] + "</span>"
			} else if rowData[i][7] == "Failure" {
				statusbadge = "<span class='badge badge-pill badge-danger'>" + rowData[i][7] + "</span>"
			} else {
				statusbadge = "--"
			}
			r = append(r, rowData[i][0], rowData[i][1], rowData[i][2], rowData[i][3], rowData[i][4], rowData[i][5], rowData[i][6], statusbadge)
			result = append(result, r)
		} else if sess.Get("role") == "ADMIN" && sess.Get("language") == "French" {
			if rowData[i][7] == "Success" {
				statusbadge = "<span class='badge badge-pill badge-success'>Succès</span>"
			} else if rowData[i][7] == "Failure" {
				statusbadge = "<span class='badge badge-pill badge-danger'>Echec</span>"
			} else {
				statusbadge = "--"
			}
			r = append(r, rowData[i][0], rowData[i][1], rowData[i][2], rowData[i][3], rowData[i][4], rowData[i][5], rowData[i][6], statusbadge)
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
