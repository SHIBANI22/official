package searchPackage

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

	"github.com/astaxie/beego"
	"ominaya.com/database/sql"
	"ominaya.com/util/log"
)

type Row struct {
	Id                string
	UserMobile        string
	UserEmail         string
	UserFirstName     string
	UserMiddleName    string
	UserLastName      string
	UserRole          string
	UserStatus        string
	UserContactNumber string
	UserDepartment    string
	UserEmployeeID    string
	UserLanguage      string
	UserCreateDate    string
}

type SearchPackage struct {
	beego.Controller
}

func (c *SearchPackage) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
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

			c.TplName = "packages/searchPackage/searchPackage.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "packages/searchPackage/searchPackage.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer  Page Success")
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

	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "Users")
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
		c.Data["SeacrhPackages"] = beego.AppConfig.String("ENGLISH_SEARCH_PACKAGES")
		c.Data["SearchFilters"] = beego.AppConfig.String("ENGLISH_SEARCH_FILTERS")

		c.Data["UserEmail"] = beego.AppConfig.String("ENGLISH_USER_EMAIL")
		c.Data["EnterUserEmail"] = beego.AppConfig.String("ENGLISH_USER_EMAIL_PLACEHOLDER")

		c.Data["name"] = beego.AppConfig.String("ENGLISH_USER_FIRST_NAME")
		c.Data["EnterUserFirstName"] = beego.AppConfig.String("ENGLISH_USER_FIRST_NAME_PLACEHOLDER")

		c.Data["SearchResults"] = beego.AppConfig.String("ENGLISH_USER_SEARCH_RESULTS")

		c.Data["UserRole"] = beego.AppConfig.String("ENGLISH_USER_ROLE")
		c.Data["PleaseSelect"] = beego.AppConfig.String("ENGLISH_USER_SELECT_ROLE")
		c.Data["UserStatus"] = beego.AppConfig.String("ENGLISH_USER_STATUS")
		c.Data["SelectDateRange"] = beego.AppConfig.String("ENGLISH_USER_SELECT_DATARANGE")
		c.Data["ID"] = beego.AppConfig.String("ENGLISH_ID")
		c.Data["Mobile"] = beego.AppConfig.String("ENGLISH_MOBILE")
		c.Data["Email"] = beego.AppConfig.String("ENGLISH_EMAIL")
		c.Data["PackageName"] = beego.AppConfig.String("ENGLISH_PACKAGE_NAME")
		c.Data["Volume"] = beego.AppConfig.String("ENGLISH_PACKAGE_VOLUME")
		c.Data["TransactionFee"] = beego.AppConfig.String("ENGLISH_PACKAGE_TXN_FEES")
		c.Data["Status"] = beego.AppConfig.String("ENGLISH_STATUS")
		c.Data["Language"] = beego.AppConfig.String("ENGLISH_LANGUAGE")
		c.Data["CreateDate"] = beego.AppConfig.String("ENGLISH_CREATE_DATE")
		c.Data["ListOfUsers"] = beego.AppConfig.String("ENGLISH_LIST_OF_PACKAGES")
		c.Data["Search"] = beego.AppConfig.String("ENGLISH_SEARCH")
		c.Data["Reset"] = beego.AppConfig.String("ENGLISH_RESET")
		c.Data["Addnew"] = beego.AppConfig.String("ENGLISH_ADDNEW")
		c.Data["input_user_email"] = beego.AppConfig.String("ENGLISH_USEREMAIL")
		c.Data["mobile"] = beego.AppConfig.String("ENGLISH_USERMOBILE")
		c.Data["please_select"] = beego.AppConfig.String("ENGLISH_PLEASESELECTCOMMON")
		c.Data["Active"] = beego.AppConfig.String("ENGLISH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("ENGLISH_INACTIVE")
		c.Data["Suspend"] = beego.AppConfig.String("ENGLISH_SUSPEND")
		c.Data["ProfileManagement"] = beego.AppConfig.String("ENGLISH_PROFILEMANAGEMENT")
		c.Data["GetRecords"] = beego.AppConfig.String("ENGLISH_GET_RECORDS")
		c.Data["List"] = beego.AppConfig.String("ENGLISH_LIST")
		c.Data["SysConfig"] = beego.AppConfig.String("ENGLISH_SYSTEM_CONFIGURATION")
		c.TplName = "packages/searchPackage/searchPackage.html"
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
		c.Data["SeacrhUsers"] = beego.AppConfig.String("FRENCH_SEARCH_USERS")
		c.Data["SearchFilters"] = beego.AppConfig.String("FRENCH_SEARCH_FILTERS")
		c.Data["SeacrhPackages"] = beego.AppConfig.String("FRENCH_SEARCH_PACKAGES")
		c.Data["SysConfig"] = beego.AppConfig.String("FRENCH_SYSTEM_CONFIGURATION")

		c.Data["UserEmail"] = beego.AppConfig.String("FRENCH_USER_EMAIL")
		c.Data["EnterUserEmail"] = beego.AppConfig.String("FRENCH_USER_EMAIL_PLACEHOLDER")

		c.Data["UserName"] = beego.AppConfig.String("FRENCH_USER_FIRST_NAME")
		c.Data["EnterUserFirstName"] = beego.AppConfig.String("FRENCH_USER_FIRST_NAME_PLACEHOLDER")

		c.Data["SearchResults"] = beego.AppConfig.String("FRENCH_USER_SEARCH_RESULTS")

		c.Data["UserRole"] = beego.AppConfig.String("FRENCH_USER_ROLE")
		c.Data["PleaseSelect"] = beego.AppConfig.String("FRENCH_USER_SELECT_ROLE")
		c.Data["UserStatus"] = beego.AppConfig.String("FRENCH_USER_STATUS")
		c.Data["SelectDateRange"] = beego.AppConfig.String("FRENCH_USER_SELECT_DATARANGE")
		c.Data["ID"] = beego.AppConfig.String("FRENCH_ID")
		c.Data["Mobile"] = beego.AppConfig.String("FRENCH_MOBILE")
		c.Data["Email"] = beego.AppConfig.String("FRENCH_EMAIL")
		c.Data["PackageName"] = beego.AppConfig.String("FRENCH_PACKAGE_NAME")
		c.Data["Volume"] = beego.AppConfig.String("FRENCH_PACKAGE_VOLUME")
		c.Data["Status"] = beego.AppConfig.String("FRENCH_STATUS")
		c.Data["Language"] = beego.AppConfig.String("FRENCH_LANGUAGE")
		c.Data["CreateDate"] = beego.AppConfig.String("FRENCH_CREATE_DATE")
		c.Data["TransactionFee"] = beego.AppConfig.String("FRENCH_TRANSACTION_FEE")
		c.Data["ListOfUsers"] = beego.AppConfig.String("FRENCH_LIST_OF_USERS")
		c.Data["Search"] = beego.AppConfig.String("FRENCH_SEARCH")
		c.Data["Reset"] = beego.AppConfig.String("FRENCH_RESET")
		c.Data["Addnew"] = beego.AppConfig.String("FRENCH_ADDNEW")
		c.Data["input_user_email"] = beego.AppConfig.String("FRENCH_USEREMAIL")
		c.Data["mobile"] = beego.AppConfig.String("FRENCH_USERMOBILE")
		c.Data["please_select"] = beego.AppConfig.String("FRENCH_PLEASESELECTCOMMON")
		c.Data["Active"] = beego.AppConfig.String("FRENCH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("FRENCH_INACTIVE")
		c.Data["Suspend"] = beego.AppConfig.String("FRENCH_SUSPEND")
		c.Data["ProfileManagement"] = beego.AppConfig.String("FRENCH_PROFILEMANAGEMENT")
		c.Data["GetRecords"] = beego.AppConfig.String("FRENCH_GET_RECORDS")
		c.Data["List"] = beego.AppConfig.String("FRENCH_LIST")

		c.TplName = "packages/searchPackage/searchPackage.html"
	}

	row, err := db.Db.Query("select id,name from Roles")

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
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("User data not found")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	c.Data["DepartArray"] = data

	return
}

func (c *SearchPackage) Post() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer API Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "Users")
	passSet := sess.Get("passwordSet").(string)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get Menus")
		return
	}
	if !auth || passSet != "YES" {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		err = errors.New("UnAuthorized")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()

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

	user_firstname := strings.TrimSpace(c.Input().Get("input_user_firstname"))
	user_email := strings.TrimSpace(c.Input().Get("input_user_email"))
	user_mobile := strings.TrimSpace(c.Input().Get("input_user_mobile"))
	//user_role := strings.TrimSpace(c.Input().Get("input_user_role"))
	user_status := strings.TrimSpace(c.Input().Get("input_user_status"))

	//start := strings.TrimSpace(c.Input().Get("start"))
	//length := strings.TrimSpace(c.Input().Get("length"))
	draw := strings.TrimSpace(c.Input().Get("draw"))
	searchValue := strings.TrimSpace(c.Input().Get("search[value]"))

	columns := [6]string{"id", "name", "volume", "txn_fees", "created_at", "status"}

	var sqlQuery bytes.Buffer
	var sqlArgs []interface{}
	sqlQuery.WriteString("select id,name,volume,txn_fees,created_at,status from Txn_Package where 1=1")
	if user_firstname != "" {
		sqlQuery.WriteString(" AND name like ?")
		sqlArgs = append(sqlArgs, user_firstname+"%")
	}
	if user_email != "" {
		sqlQuery.WriteString(" AND volume like ?")
		sqlArgs = append(sqlArgs, user_email+"%")
	}
	if user_mobile != "" {
		sqlQuery.WriteString(" AND txn_fees  like ?")
		sqlArgs = append(sqlArgs, user_mobile+"%")
	}
	/*if user_role != "" {
		sqlQuery.WriteString(" AND u.role_id=?")
		sqlArgs = append(sqlArgs, user_role)
	}*/
	if user_status != "" {
		sqlQuery.WriteString(" AND status=?")
		sqlArgs = append(sqlArgs, user_status)
	}
	/*tmp := strings.Split(c.Input().Get("daterange"), " - ")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", c.Input().Get("daterange"))
	if c.Input().Get("daterange") != "" && len(tmp) == 2 {
		format, dateErr := utils.CustomDateFormat(beego.AppConfig.String("DateFormat"), "")
		if dateErr != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", dateErr)
		}

		inputFromDate, _ := time.Parse(format, tmp[0])
		inputToDate, _ := time.Parse(format, tmp[1])
		from := inputFromDate.Format("2006-01-02")
		to := inputToDate.Format("2006-01-02")
		sqlQuery.WriteString(" AND (DATE (u.created_date) >= DATE(?)) AND (DATE (u.created_date) <= DATE(?))")
		log.Println(beego.AppConfig.String("loglevel"), "Debug", from, to)
		sqlArgs = append(sqlArgs, from, to)
	}*/
	if searchValue != "" {
		//sqlQuery.WriteString(" AND (name like ? OR u.email like ? OR u.mobile like ? OR u.status like ? OR u.last_name like ? OR u.middle_name like ? OR u.language like ? OR r.name like ?)")
		sqlQuery.WriteString(" AND (name like ? OR volume like ? OR txn_fees like ? OR status like ? )")
		//sqlArgs = append(sqlArgs, searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%", searchValue+"%")
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
	/*lengthInt, _ := strconv.Atoi(length)
	if lengthInt > 0 {
		sqlQuery.WriteString(" LIMIT ?,?")
		sqlArgs = append(sqlArgs, start, length)
	}*/
	//log.Println("Debug", "Info", "Final Query", sqlQuery.String(), sqlArgs)

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
			viewLink := "<a href='/Packages/ViewPackage/" + rowData[i][0] + "'><h    6 class='text-red' data-toggle='tooltip' data-placement='top' data-original-title='View'>" + rowData[i][0] + "</h6></a>"
			r = append(r, viewLink, rowData[i][1], rowData[i][2], rowData[i][3], rowData[i][4], statusbadge)
		} else if sess.Get("role") == "ADMIN" && sess.Get("language") == "French" {
			if rowData[i][5] == "ACTIVE" {
				statusbadge = "<span class='badge badge-pill badge-success'>Actif</span>"
			} else if rowData[i][5] == "INACTIVE" {
				statusbadge = "<span class='badge badge-pill badge-danger'>InActif</span>"
			} else {
				statusbadge = "--"
			}
			viewLink := "<a href='/Packages/ViewPackage/" + rowData[i][0] + "'><h6 class='text-red' data-toggle='tooltip' data-placement='top' data-original-title='Consulter'>" + rowData[i][0] + "</h6></a>"
			r = append(r, viewLink, rowData[i][0], rowData[i][1], rowData[i][2], rowData[i][3], statusbadge)
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
