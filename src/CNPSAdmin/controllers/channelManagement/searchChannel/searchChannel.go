package searchChannel

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"ominaya.com/database/sql"
	"ominaya.com/util/log"
)

type Row struct {
	ID            string
	GatewayName   string
	GatewayStatus string
	ChannelName   string
	ChannelStatus string
	ChannelDesc   string
	CreateDate    string
}

type SearchChannel struct {
	beego.Controller
}

func (c *SearchChannel) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Start")
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
			c.TplName = "channelmanagement/searchChannel/searchChannel.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "channelmanagement/searchChannel/searchChannel.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel  Page Success")
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "ViewChannels")
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
		c.Data["Channels"] = beego.AppConfig.String("ENGLISH_CHANNELS")
		c.Data["Channels1"] = beego.AppConfig.String("ENGLISH_CHANNELS")
		c.Data["ChannelName"] = beego.AppConfig.String("ENGLISH_CHANNEL_NAME")
		c.Data["GatewayName"] = beego.AppConfig.String("ENGLISH_GATEWAY_NAME")
		c.Data["SelectDateRange"] = beego.AppConfig.String("ENGLISH_USER_SELECT_DATARANGE")
		c.Data["GatewayStatus"] = beego.AppConfig.String("ENGLISH_GATEWAY_STATUS")
		c.Data["ChannelStatus"] = beego.AppConfig.String("ENGLISH_CHANNEL_STATUS")
		c.Data["ID"] = beego.AppConfig.String("ENGLISH_ID")
		c.Data["TimeStamp"] = beego.AppConfig.String("ENGLISH_TIMESTAMP")
		c.Data["Desc"] = beego.AppConfig.String("ENGLISH_DESC")
		c.Data["Status"] = beego.AppConfig.String("ENGLISH_STATUS")
		c.Data["SearchFilters"] = beego.AppConfig.String("ENGLISH_SEARCH_FILTERS")
		c.Data["ListOfChannels"] = beego.AppConfig.String("ENGLISH_LIST_OF_CHANNELS")
		c.Data["chenaelPlaceholder"] = beego.AppConfig.String("ENGLISH_CHANNEL_NAME")
		c.Data["chenaelPlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_CHANNEL_NAME")
		c.Data["gatewayPlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_GATEWAY_NAME")
		c.Data["PLEASESELECTPlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_PLEASESELECT")
		c.Data["Search"] = beego.AppConfig.String("ENGLISH_SEARCH")
		c.Data["Reset"] = beego.AppConfig.String("ENGLISH_RESET")
		c.Data["Active"] = beego.AppConfig.String("ENGLISH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("ENGLISH_INACTIVE")
		c.Data["SysConfig"] = beego.AppConfig.String("ENGLISH_SYSTEM_CONFIGURATION")
		c.Data["CreateDate"] = beego.AppConfig.String("ENGLISH_CREATE_DATE")

		c.Data["GetRecords"] = beego.AppConfig.String("ENGLISH_GET_RECORDS")
		c.Data["List"] = beego.AppConfig.String("ENGLISH_LIST")

		c.TplName = "channelmanagement/searchChannel/searchChannel.html"
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
		c.Data["Channels"] = beego.AppConfig.String("FRENCH_CHANNELS")
		c.Data["Channels1"] = beego.AppConfig.String("FRENCH_CHANNELS")
		c.Data["ChannelName"] = beego.AppConfig.String("FRENCH_CHANNEL_NAME")
		c.Data["GatewayName"] = beego.AppConfig.String("FRENCH_GATEWAY_NAME")
		c.Data["SelectDateRange"] = beego.AppConfig.String("FRENCH_USER_SELECT_DATARANGE")
		c.Data["GatewayStatus"] = beego.AppConfig.String("FRENCH_GATEWAY_STATUS")
		c.Data["ChannelStatus"] = beego.AppConfig.String("FRENCH_CHANNEL_STATUS")
		c.Data["ID"] = beego.AppConfig.String("FRENCH_ID")
		c.Data["TimeStamp"] = beego.AppConfig.String("FRENCH_TIMESTAMP")
		c.Data["Status"] = beego.AppConfig.String("FRENCH_STATUS")
		c.Data["Desc"] = beego.AppConfig.String("FRENCH_DESC")
		c.Data["SearchFilters"] = beego.AppConfig.String("FRENCH_SEARCH_FILTERS")
		c.Data["ListOfChannels"] = beego.AppConfig.String("FRENCH_LIST_OF_CHANNELS")
		c.Data["chenaelPlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_CHANNEL_NAME")
		c.Data["gatewayPlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_GATEWAY_NAME")
		c.Data["PLEASESELECTPlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_PLEASESELECT")
		c.Data["Search"] = beego.AppConfig.String("FRENCH_SEARCH")
		c.Data["Reset"] = beego.AppConfig.String("FRENCH_RESET")
		c.Data["Active"] = beego.AppConfig.String("FRENCH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("FRENCH_INACTIVE")
		c.Data["SysConfig"] = beego.AppConfig.String("FRENCH_SYSTEM_CONFIGURATION")
		c.Data["CreateDate"] = beego.AppConfig.String("FRENCH_CREATE_DATE")

		c.Data["GetRecords"] = beego.AppConfig.String("FRENCH_GET_RECORDS")
		c.Data["List"] = beego.AppConfig.String("FRENCH_LIST")

		c.TplName = "channelmanagement/searchChannel/searchChannel.html"
	}

	return
}

func (c *SearchChannel) Post() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Page Start")
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "ViewChannels")
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
	channel_name := c.Input().Get("channel_name")
	channel_status := c.Input().Get("input_status")
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

	columns := [6]string{"pc.uuid", "pc.channel_name", "pc.create_date", "pg.gateway_name", "pc.channel_desc", "pc.channel_status"}

	var sqlQuery bytes.Buffer
	var sqlArgs []interface{}
	sqlQuery.WriteString("select pc.uuid, pc.channel_name, pc.create_date, pg.gateway_name, pc.channel_desc, pc.channel_status from Payment_Channel as pc inner join Payment_Gateway as pg where (pc.payment_gateway_id = pg.id)")
	if channel_name != "" {
		sqlQuery.WriteString(" AND pc.channel_name like ?")
		sqlArgs = append(sqlArgs, channel_name+"%")
	}
	if channel_status != "" {
		sqlQuery.WriteString(" AND pc.channel_status=?")
		sqlArgs = append(sqlArgs, channel_status)
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
		sqlQuery.WriteString(" AND (DATE (pc.create_date) >= DATE(?)) AND (DATE (pc.create_date) <= DATE(?))")
		log.Println(beego.AppConfig.String("loglevel"), "Debug", from, to)
		sqlArgs = append(sqlArgs, from, to)
	}
	if searchValue != "" {
		sqlQuery.WriteString(" AND (pc.channel_name like ? OR pg.gateway_name like ? OR pc.channel_desc like ? OR pc.channel_status like ?)")
		sqlArgs = append(sqlArgs, searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%")
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
			if rowData[i][5] == "ACTIVE" {
				statusbadge = "<span class='badge badge-pill badge-success'>Active</span>"
			} else if rowData[i][5] == "INACTIVE" {
				statusbadge = "<span class='badge badge-pill badge-danger'>InActive</span>"
			} else {
				statusbadge = "--"
			}
			viewLink := "<a href='" + beego.URLFor("ViewChannel.Get", ":AdminID", rowData[i][0]) + "'><h6 class='text-red' data-toggle='tooltip' data-placement='top' data-original-title='View'>" + rowData[i][0] + "</h6></a>"
			r = append(r, viewLink, rowData[i][1], rowData[i][2], rowData[i][3], rowData[i][4], statusbadge)
		} else if sess.Get("role") == "ADMIN" && sess.Get("language") == "French" {
			if rowData[i][5] == "ACTIVE" {
				statusbadge = "<span class='badge badge-pill badge-success'>Actif</span>"
			} else if rowData[i][5] == "INACTIVE" {
				statusbadge = "<span class='badge badge-pill badge-danger'>InActif</span>"
			} else {
				statusbadge = "--"
			}
			viewLink := "<a href='" + beego.URLFor("ViewChannel.Get", ":AdminID", rowData[i][0]) + "'><h6 class='text-red' data-toggle='tooltip' data-placement='top' data-original-title='Consulter'>" + rowData[i][0] + "</h6></a>"
			r = append(r, viewLink, rowData[i][1], rowData[i][2], rowData[i][3], rowData[i][4], statusbadge)
		}
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
