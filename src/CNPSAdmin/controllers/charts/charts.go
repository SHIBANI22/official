package charts

import (
	"CNPSAdmin/model/db"
	"sort"
	"strings"

	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"errors"
	"fmt"
	"html/template"

	"runtime/debug"

	"time"

	"strconv"

	"github.com/astaxie/beego"
	"ominaya.com/database/sql"
	"ominaya.com/util/log"
)

type ChartController struct {
	beego.Controller
}

var (
	monthsNames                 = make(map[interface{}]interface{})
	uname                       string
	language                    string
	approvedLabel               string
	pendingLabel                string
	declinedLabel               string
	transactionAmountChartTitle string
	totalAmountLabel            string
	transactionCountTitle       string
	transactionByOptrTitle      string
	transactionByChannelTitle   string
)

func (c *ChartController) Prepare() {
	log.Println(beego.AppConfig.String("loglevel"), "info", "Start: Transaction Status API")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Current IP: ", pip)
	var err error
	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.TplName = "error/error.html"
		}
		if err != nil && err.Error() == "Session Time Out.Please Logout and Login Again." {
			c.Data["DisplayMessage"] = err.Error()
			c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)
			//log.Println(beego.AppConfig.String("loglevel"), "Info", "Transaction Status API Fail")
		} else if err != nil {
			c.Data["DisplayMessage"] = err.Error()
			c.TplName = "error/error.html"
			//log.Println(beego.AppConfig.String("loglevel"), "Info", "Transaction Status API Fail")
		} else {
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Transaction Status API Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Somthing went wrong. Please contact customer care")
		return
	}

	if err = session.ValidateSession(sess); err != nil {
		sess.SessionRelease(c.Ctx.ResponseWriter)
		session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()
	uname = sess.Get("uname").(string)
	language = sess.Get("language").(string)

	if language == "French" {
		monthsNames["Jan"] = "Jan"
		monthsNames["Feb"] = "Fév"
		monthsNames["Mar"] = "Mar"
		monthsNames["Apr"] = "Avr"
		monthsNames["May"] = "Mai"
		monthsNames["Jun"] = "Jui"
		monthsNames["Jul"] = "Juil"
		monthsNames["Aug"] = "Aoû"
		monthsNames["Sep"] = "Sep"
		monthsNames["Oct"] = "Oct"
		monthsNames["Nov"] = "Nov"
		monthsNames["Dec"] = "Déc"
		approvedLabel = beego.AppConfig.String("FRENCH_APPROVED")
		declinedLabel = beego.AppConfig.String("FRENCH_DECLINED")
		pendingLabel = beego.AppConfig.String("FRENCH_PENDING")
		transactionAmountChartTitle = beego.AppConfig.String("FRENCH_TRANSACTION_AMOUNT_DETAIL_LABEL")
		totalAmountLabel = beego.AppConfig.String("FRENCH_TOTALAMOUNT_OF_SUCCESSFUL_TRANSACTIONS")
		transactionCountTitle = beego.AppConfig.String("FRENCH_TRANSACTION_COUNT_TITLE")
		transactionByOptrTitle = beego.AppConfig.String("FRENCH_TRANSACTION_BY_OPERATOR_TITLE")
		transactionByChannelTitle = beego.AppConfig.String("FRENCH_TRANSACTION_BY_CHANNEL_TITLE")
	} else {
		monthsNames["Jan"] = "Jan"
		monthsNames["Feb"] = "Feb"
		monthsNames["Mar"] = "Mar"
		monthsNames["Apr"] = "Apr"
		monthsNames["May"] = "May"
		monthsNames["Jun"] = "Jun"
		monthsNames["Jul"] = "Jul"
		monthsNames["Aug"] = "Aug"
		monthsNames["Sep"] = "Sep"
		monthsNames["Oct"] = "Oct"
		monthsNames["Nov"] = "Nov"
		monthsNames["Dec"] = "Dec"
		approvedLabel = beego.AppConfig.String("ENGLISH_APPROVED")
		declinedLabel = beego.AppConfig.String("ENGLISH_DECLINED")
		pendingLabel = beego.AppConfig.String("ENGLISH_PENDING")
		transactionAmountChartTitle = beego.AppConfig.String("ENGLISH_TRANSACTION_AMOUNT_DETAIL_LABEL")
		totalAmountLabel = beego.AppConfig.String("ENGLISH_TOTALAMOUNT_OF_SUCCESSFUL_TRANSACTIONS")
		transactionCountTitle = beego.AppConfig.String("ENGLISH_TRANSACTION_COUNT_TITLE")
		transactionByOptrTitle = beego.AppConfig.String("ENGLISH_TRANSACTION_BY_OPERATOR_TITLE")
		transactionByChannelTitle = beego.AppConfig.String("ENGLISH_TRANSACTION_BY_CHANNEL_TITLE")
	}
}

func (c *ChartController) TransactionAmountByStatus() {
	current_time := time.Now().UTC()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "current_time", current_time)

	var m, _ = strconv.Atoi("01")
	from_date := time.Date(current_time.Year(), time.Month(m), 1, 0, 0, 0, 0, time.UTC)

	tmp := strings.Split(c.Input().Get("daterange"), " - ")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", c.Input().Get("daterange"))

	post_from_date := ""
	post_to_date := ""
	if len(tmp) == 2 {
		post_from_date = tmp[0]
		post_to_date = tmp[1]
	}

	if post_from_date != "" {
		from_date, _ = time.Parse("2006/01/02", post_from_date)
	}

	to_date := time.Date(current_time.Year(), current_time.Month(), current_time.Day(), 0, 0, 0, 0, time.UTC)

	if post_to_date != "" {
		to_date, _ = time.Parse("2006/01/02", post_to_date)
	}

	days := int(to_date.Sub(from_date).Hours() / 24)

	var mode string
	if days < 32 {
		mode = "D"
	} else {
		mode = "M"
	}

	if from_date.Year() != to_date.Year() && mode == "M" {
		mode = "MY"
	}
	if from_date.Month() != to_date.Month() && mode == "D" {
		mode = "DM"
	}

	start := from_date.Unix()
	end := to_date.Unix()

	var approved_map = make(map[interface{}]int, 0)
	var declined_map = make(map[interface{}]int, 0)
	var pending_map = make(map[interface{}]int, 0)
	var months []interface{}
	for start <= end {
		start_date := time.Unix(start, 0)
		var key interface{}
		if mode == "Y" {
			start_date = start_date.AddDate(1, 0, 0)
		} else if mode == "MY" {
			mn := start_date.Month().String()
			key = fmt.Sprintf("%s-%d", mn[:3], start_date.Year())
			start_date = start_date.AddDate(0, 1, 0)
		} else if mode == "M" {
			mn := start_date.Month().String()
			key = mn[:3]
			start_date = start_date.AddDate(0, 1, 0)
		} else if mode == "DM" {
			mn := start_date.Month().String()
			key = fmt.Sprintf("%d-%s", start_date.Day(), mn[:3])
			start_date = start_date.AddDate(0, 0, 1)
		} else if mode == "D" {
			key = start_date.Day()
			start_date = start_date.AddDate(0, 0, 1)
		}
		approved_map[key] = 0
		declined_map[key] = 0
		pending_map[key] = 0
		months = append(months, key)
		start = start_date.Unix()
	}

	var date_string string
	if !from_date.IsZero() && !to_date.IsZero() {
		date_string = fmt.Sprintf(" DATE(transaction_time) >= '%s' AND DATE(transaction_time) <= '%s'", from_date.Format("2006-01-02"), to_date.Format("2006-01-02"))
	} else {
		date_string = fmt.Sprintf(" YEAR(transaction_time) = %s", time.Now().Year())
	}

	var group_by string
	if mode == "Y" {
		group_by = " GROUP BY y"
	} else if mode == "D" {
		group_by = " GROUP BY df"
	} else {
		group_by = " GROUP BY ym"
	}

	rows, err := db.Db.Query(`SELECT SUM(CASE WHEN status = "APPROVED" THEN amount ELSE 0 END) AS sta, SUM(CASE WHEN status = "PENDING" THEN amount ELSE 0 END) AS pta, SUM(CASE WHEN status="DECLINED" THEN amount ELSE 0 END) AS dta, DATE(transaction_time) as df, MONTH(transaction_time) as mn, YEAR(transaction_time) AS y, date_format(transaction_time, '%b') AS month, date_format(transaction_time, '%b-%Y') AS ym, date_format(transaction_time, '%Y-%b') AS yM, date_format(transaction_time, '%e-%b') AS dM, date_format(transaction_time, '%e') AS d FROM Transactions WHERE ` + date_string + group_by)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Db selection error")
		return
	}

	_, data, err := sql.Scan(rows)
	defer sql.Close(rows)
	if err != nil {
		err = errors.New("Columns not found")
		return
	}

	for _, v := range data {
		approvedAmount, _ := strconv.Atoi(v[0])
		pendingAmount, _ := strconv.Atoi(v[1])
		declinedAmount, _ := strconv.Atoi(v[2])
		if mode == "Y" {
			year := v[5]
			approved_map[year] = approvedAmount
			declined_map[year] = declinedAmount
			pending_map[year] = pendingAmount
		} else if mode == "MY" {
			month_year := v[7]
			approved_map[month_year] = approvedAmount
			declined_map[month_year] = declinedAmount
			pending_map[month_year] = pendingAmount
		} else if mode == "M" {
			month := v[6]
			approved_map[month] = approvedAmount
			declined_map[month] = declinedAmount
			pending_map[month] = pendingAmount
		} else if mode == "DM" {
			day_month := v[9]
			approved_map[day_month] = approvedAmount
			declined_map[day_month] = declinedAmount
			pending_map[day_month] = pendingAmount
		} else if mode == "D" {
			day, _ := strconv.Atoi(v[10])
			approved_map[day] = approvedAmount
			declined_map[day] = declinedAmount
			pending_map[day] = pendingAmount
		}
	}

	var approved_transaction []int
	var declined_transaction []int
	var pending_transaction []int
	for _, v := range months {
		approved_transaction = append(approved_transaction, approved_map[v])
		declined_transaction = append(declined_transaction, declined_map[v])
		pending_transaction = append(pending_transaction, pending_map[v])
	}
	finalMonths := translateMonth(mode, months)

	finalData := make(map[string]interface{})
	finalData["approved_transaction"] = approved_transaction
	finalData["declined_transaction"] = declined_transaction
	finalData["pending_transaction"] = pending_transaction
	finalData["months"] = finalMonths
	finalData["approved_title"] = approvedLabel
	finalData["declined_title"] = declinedLabel
	finalData["pending_title"] = pendingLabel
	finalData["transaction_amount_detail_label"] = transactionAmountChartTitle
	finalData["total_amount_label"] = totalAmountLabel
	//fmt.Println("Final data: ", finalData)
	c.Data["json"] = finalData
	c.ServeJSON()

}

func (this *ChartController) TransactionCountByStatus() {

	current_time := time.Now()

	var m, _ = strconv.Atoi("01")
	from_date := time.Date(current_time.Year(), time.Month(m), 1, 0, 0, 0, 0, time.UTC)

	tmp := strings.Split(this.Input().Get("daterange"), " - ")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", this.Input().Get("daterange"))

	post_from_date := ""
	post_to_date := ""
	if len(tmp) == 2 {
		post_from_date = tmp[0]
		post_to_date = tmp[1]
	}

	if post_from_date != "" {
		from_date, _ = time.Parse("2006/01/02", post_from_date)
	}

	to_date := time.Date(current_time.Year(), current_time.Month(), current_time.Day(), 0, 0, 0, 0, time.UTC)

	if post_to_date != "" {
		to_date, _ = time.Parse("2006/01/02", post_to_date)
	}

	days := int(to_date.Sub(from_date).Hours() / 24)

	var mode string
	if days < 32 {
		mode = "D"
	} else {
		mode = "M"
	}

	if from_date.Year() != to_date.Year() && mode == "M" {
		mode = "MY"
	}
	if from_date.Month() != to_date.Month() && mode == "D" {
		mode = "DM"
	}

	start := from_date.Unix()
	end := to_date.Unix()

	var approved_map = make(map[interface{}]int, 0)
	var declined_map = make(map[interface{}]int, 0)
	var pending_map = make(map[interface{}]int, 0)
	var months []interface{}
	for start <= end {
		start_date := time.Unix(start, 0)
		var key interface{}
		if mode == "Y" {
			start_date = start_date.AddDate(1, 0, 0)
		} else if mode == "MY" {
			mn := start_date.Month().String()
			key = fmt.Sprintf("%s-%d", mn[:3], start_date.Year())
			start_date = start_date.AddDate(0, 1, 0)
		} else if mode == "M" {
			mn := start_date.Month().String()
			key = mn[:3]
			start_date = start_date.AddDate(0, 1, 0)
		} else if mode == "DM" {
			mn := start_date.Month().String()
			key = fmt.Sprintf("%d-%s", start_date.Day(), mn[:3])
			start_date = start_date.AddDate(0, 0, 1)
		} else if mode == "D" {
			key = start_date.Day()
			start_date = start_date.AddDate(0, 0, 1)
		}
		approved_map[key] = 0
		declined_map[key] = 0
		pending_map[key] = 0
		months = append(months, key)
		start = start_date.Unix()
	}

	var date_string string
	if !from_date.IsZero() && !to_date.IsZero() {
		date_string = fmt.Sprintf(" DATE(transaction_time) >= '%s' AND DATE(transaction_time) <= '%s'", from_date.Format("2006-01-02"), to_date.Format("2006-01-02"))
	} else {
		date_string = fmt.Sprintf(" YEAR(transaction_time) = %s", time.Now().Year())
	}

	var group_by string
	if mode == "Y" {
		group_by = " GROUP BY y"
	} else if mode == "D" {
		group_by = " GROUP BY df"
	} else {
		group_by = " GROUP BY ym"
	}

	rows, err := db.Db.Query(`SELECT COUNT(CASE WHEN status = "APPROVED" THEN 1 END) AS sta, COUNT(CASE WHEN status = "PENDING" THEN 1 END) AS pta, COUNT(CASE WHEN status="DECLINED" THEN 1 END) AS dta, DATE(transaction_time) as df, MONTH(transaction_time) as mn, YEAR(transaction_time) AS y, date_format(transaction_time, '%b') AS month, date_format(transaction_time, '%b-%Y') AS ym, date_format(transaction_time, '%Y-%b') AS yM, date_format(transaction_time, '%e-%b') AS dM, date_format(transaction_time, '%e') AS d FROM Transactions WHERE ` + date_string + group_by)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Db selection error")
		return
	}

	_, data, err := sql.Scan(rows)
	defer sql.Close(rows)
	if err != nil {
		err = errors.New("Columns not found")
		return
	}

	for _, v := range data {
		approvedCount, _ := strconv.Atoi(v[0])
		pendingCount, _ := strconv.Atoi(v[1])
		declinedCount, _ := strconv.Atoi(v[2])
		if mode == "Y" {
			year := v[5]
			approved_map[year] = approvedCount
			declined_map[year] = declinedCount
			pending_map[year] = pendingCount
		} else if mode == "MY" {
			month_year := v[7]
			approved_map[month_year] = approvedCount
			declined_map[month_year] = declinedCount
			pending_map[month_year] = pendingCount
		} else if mode == "M" {
			month := v[6]
			approved_map[month] = approvedCount
			declined_map[month] = declinedCount
			pending_map[month] = pendingCount
		} else if mode == "DM" {
			day_month := v[9]
			approved_map[day_month] = approvedCount
			declined_map[day_month] = declinedCount
			pending_map[day_month] = pendingCount
		} else if mode == "D" {
			day, _ := strconv.Atoi(v[10])
			approved_map[day] = approvedCount
			declined_map[day] = declinedCount
			pending_map[day] = pendingCount
		}
	}

	var approved_transaction []int
	var declined_transaction []int
	var pending_transaction []int
	for _, v := range months {
		approved_transaction = append(approved_transaction, approved_map[v])
		declined_transaction = append(declined_transaction, declined_map[v])
		pending_transaction = append(pending_transaction, pending_map[v])
	}

	finalMonths := translateMonth(mode, months)

	finalData := make(map[string]interface{})
	finalData["approved_transaction"] = approved_transaction
	finalData["declined_transaction"] = declined_transaction
	finalData["pending_transaction"] = pending_transaction
	finalData["title"] = transactionCountTitle
	finalData["months"] = finalMonths
	finalData["approved_title"] = approvedLabel
	finalData["declined_title"] = declinedLabel
	finalData["pending_title"] = pendingLabel
	// fmt.Println("Final data: ", finalData)
	this.Data["json"] = finalData
	this.ServeJSON()

}
func (this *ChartController) TransactionByOperator() {

	status := this.GetString("status")
	if status == "" {
		status = "ALL"
	}

	CNPS_OPERATORS := beego.AppConfig.String("CNPS_OPERATORS")

	var operateQuery []string
	CNPS_OPERATORS_LIST := strings.Split(CNPS_OPERATORS, "|")
	var amount []string
	for _, operator := range CNPS_OPERATORS_LIST {
		operateQuery = append(operateQuery, "SUM(CASE WHEN operator = '"+operator+"' THEN amount ELSE 0 END)")
	}

	operateQueryFinal := strings.Join(operateQuery, ", ")

	if len(CNPS_OPERATORS_LIST) > 0 && CNPS_OPERATORS != "" {
		rows, err := db.Db.Query("SELECT "+operateQueryFinal+" FROM Transactions WHERE operator IS NOT NULL AND status=? OR (1=(CASE WHEN ('ALL'=? AND status != 'INITIATED') THEN 1 ELSE 0 END))", status, status)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Db selection error")
			return
		}
		_, data, err := sql.Scan(rows)
		defer sql.Close(rows)
		if err != nil {
			err = errors.New("Columns not found")
			return
		}
		amount = data[0]
	}

	finalData := make(map[string]interface{})
	finalData["operator_list"] = CNPS_OPERATORS_LIST
	finalData["amount"] = amount
	finalData["title"] = transactionByOptrTitle
	finalData["chart_label"] = totalAmountLabel + " (XOF)"
	// fmt.Println("Final data: ", finalData)
	this.Data["json"] = finalData
	this.ServeJSON()

}

func (this *ChartController) TransactionByChannel() {

	status := this.GetString("status")
	if status == "" {
		status = "ALL"
	}

	CNPS_CHANNELS := beego.AppConfig.String("CNPS_CHANNELS")

	var channelQuery []string
	CNPS_CHANNELS_LIST := strings.Split(CNPS_CHANNELS, "|")
	var amount []string
	for _, channel := range CNPS_CHANNELS_LIST {
		channelQuery = append(channelQuery, "SUM(CASE WHEN channel = '"+channel+"' THEN amount ELSE 0 END)")
	}
	channelQueryFinal := strings.Join(channelQuery, ", ")
	if len(CNPS_CHANNELS_LIST) > 0 && CNPS_CHANNELS != "" {
		rows, err := db.Db.Query("SELECT "+channelQueryFinal+" FROM Transactions WHERE operator IS NOT NULL AND status=? OR (1=(CASE WHEN ('ALL'=? AND status != 'INITIATED') THEN 1 ELSE 0 END))", status, status)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Db selection error")
			return
		}
		_, data, err := sql.Scan(rows)
		defer sql.Close(rows)
		if err != nil {
			err = errors.New("Columns not found")
			return
		}
		amount = data[0]
	}

	finalData := make(map[string]interface{})
	finalData["channel_list"] = CNPS_CHANNELS_LIST
	finalData["amount"] = amount
	finalData["title"] = transactionByChannelTitle
	finalData["chart_label"] = totalAmountLabel + " (XOF)"
	// fmt.Println("Final data: ", finalData)
	this.Data["json"] = finalData
	this.ServeJSON()

}

func translateMonth(mode string, months []interface{}) (finalMonths []interface{}) {
	// fmt.Println("months", months, "Mode ", mode)
	for _, m := range months {
		if mode == "M" {
			finalMonths = append(finalMonths, monthsNames[m])
		} else if mode == "DM" {
			dm := strings.Split(m.(string), "-")
			mn := monthsNames[dm[1]]
			fm := dm[0] + "-" + mn.(string)
			finalMonths = append(finalMonths, fm)
		} else if mode == "MY" {
			my := strings.Split(m.(string), "-")
			mn := monthsNames[my[0]]
			fm := mn.(string) + "-" + my[1]
			finalMonths = append(finalMonths, fm)
		} else {
			finalMonths = append(finalMonths, m)
		}
	}
	return
}

func (c *ChartController) TransactionLiveData() {
	current_time := time.Now()
	from_date := time.Date(current_time.Year(), current_time.Month(), current_time.Day(), 8, 0, 0, 0, time.UTC)

	to_date := time.Date(current_time.Year(), current_time.Month(), current_time.Day(), 23, 0, 0, 0, time.UTC)

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "from_to_date", from_date, to_date)
	start := from_date.Unix()
	end := to_date.Unix()

	var date_map = make(map[interface{}]int, 0)
	for start <= end {
		start_date := time.Unix(start, 0).UTC()
		var key interface{}
		key = start_date.Format("2006-01-02 15:04")
		start_date = start_date.Add(time.Hour * 1)

		date_map[key] = 0
		start = start_date.Unix()
	}

	rows, err := db.Db.Query(`SELECT SUM(amount) AS a, DATE_FORMAT(transaction_time, "%Y-%m-%d %H:%i") as hm FROM Transactions WHERE status="APPROVED" AND DATE(transaction_time)=CURDATE() GROUP BY hm`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Db selection error")
		return
	}

	_, data, err := sql.Scan(rows)
	defer sql.Close(rows)
	if err != nil {
		err = errors.New("Columns not found")
		return
	}
	// log.Println("Debug", "Info", "data", data)
	for _, v := range data {
		date_map[v[1]] = 0
	}

	var kyes []string
	for key, _ := range date_map {
		kyes = append(kyes, key.(string))
	}
	sort.Strings(kyes)
	var dateSorted = make(map[string]int)
	for _, v := range kyes {
		dateSorted[v] = 0
	}
	for _, v := range data {
		amount, _ := strconv.Atoi(v[0])
		day := v[1]
		dateSorted[day] = amount
	}

	finalData := make(map[string]interface{})
	finalData["transaction"] = dateSorted
	c.Data["json"] = finalData
	c.ServeJSON()
}

func (c *ChartController) TransactionLiveDataStaticUpdateBackup() {
	rows, err := db.Db.Query(`SELECT SUM(CASE WHEN status = "APPROVED" THEN amount ELSE 0 END) AS sta, DATE_FORMAT(transaction_time, "%Y-%m-%d %H:%i") as hm FROM Transactions WHERE status="APPROVED" AND DATE(transaction_time) = CURRENT_DATE() GROUP BY hm`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Db selection error")
		return
	}

	_, data, err := sql.Scan(rows)
	defer sql.Close(rows)
	if err != nil {
		err = errors.New("Columns not found")
		return
	}
	log.Println("Debug", "Info", "data", data)

	var xAxis []interface{}
	var yAxis []interface{}
	for _, v := range data {
		log.Println("Debug", "Info", "v", v[0], v[1])
		xAxis = append(xAxis, v[1])
		yAxis = append(yAxis, v[0])
	}

	finalData := make(map[string]interface{})
	finalData["xAxis"] = xAxis
	finalData["yAxis"] = yAxis
	// fmt.Println("Final data: ", finalData)
	c.Data["json"] = finalData
	c.ServeJSON()

}

func (c *ChartController) TransactionLiveMovingDataBackup() {
	current_time := time.Now()
	fmt.Println("current_time", current_time)
	// var m, _ = strconv.Atoi("01")
	from_date := time.Date(current_time.Year(), current_time.Month(), current_time.Day(), current_time.Hour(), current_time.Minute(), 0, 0, time.UTC).Add(time.Hour * -4)

	// tmp := strings.Split(c.Input().Get("daterange"), " - ")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "from_date", from_date)

	// post_from_date := ""
	// post_to_date := ""
	// if len(tmp) == 2 {
	// 	post_from_date = tmp[0]
	// 	post_to_date = tmp[1]
	// }

	// if post_from_date != "" {
	// 	from_date, _ = time.Parse("2006/01/02", post_from_date)
	// }

	to_date := time.Date(current_time.Year(), current_time.Month(), current_time.Day(), current_time.Hour(), current_time.Minute(), 0, 0, time.UTC)

	// if post_to_date != "" {
	// 	to_date, _ = time.Parse("2006/01/02", post_to_date)
	// }

	days := int(to_date.Sub(from_date).Hours() / 24)

	var mode string
	if days < 32 {
		mode = "D"
	} else {
		mode = "M"
	}

	if from_date.Year() != to_date.Year() && mode == "M" {
		mode = "MY"
	}
	if from_date.Month() != to_date.Month() && mode == "D" {
		mode = "DM"
	}
	mode = "HM"
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "from_to_date", from_date, to_date)

	start := from_date.Unix()
	end := to_date.Unix()

	var approved_map = make(map[interface{}]int, 0)
	var months []interface{}
	for start <= end {
		start_date := time.Unix(start, 0).UTC()
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "start_date_new", start_date.Format("2006-01-02 15:04"))
		var key interface{}
		if mode == "Y" {
			start_date = start_date.AddDate(1, 0, 0)
		} else if mode == "MY" {
			mn := start_date.Month().String()
			key = fmt.Sprintf("%s-%d", mn[:3], start_date.Year())
			start_date = start_date.AddDate(0, 1, 0)
		} else if mode == "M" {
			mn := start_date.Month().String()
			key = mn[:3]
			start_date = start_date.AddDate(0, 1, 0)
		} else if mode == "DM" {
			mn := start_date.Month().String()
			key = fmt.Sprintf("%d-%s", start_date.Day(), mn[:3])
			start_date = start_date.AddDate(0, 0, 1)
		} else if mode == "D" {
			key = start_date.Day()
			start_date = start_date.AddDate(0, 0, 1)
		} else if mode == "HM" {
			key = start_date.Format("2006-01-02 15:04")
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "time_date", start_date.Hour(), start_date.Minute(), key)
			start_date = start_date.Add(time.Minute * 1)
		}
		approved_map[key] = 0
		months = append(months, key)
		start = start_date.Unix()
	}
	log.Println("Debug", "Info", "approved_map", approved_map)

	// var date_string string
	// if !from_date.IsZero() && !to_date.IsZero() {
	// 	date_string = fmt.Sprintf(" DATE(transaction_time) >= '%s' AND DATE(transaction_time) <= '%s'", from_date.Format("2006-01-02"), to_date.Format("2006-01-02"))
	// } else {
	// 	date_string = fmt.Sprintf(" YEAR(transaction_time) = %s", time.Now().Year())
	// }

	var group_by string
	if mode == "Y" {
		group_by = " GROUP BY y"
	} else if mode == "D" {
		group_by = " GROUP BY df"
	} else {
		group_by = " GROUP BY hm"
	}
	group_by = " GROUP BY hm"
	rows, err := db.Db.Query(`SELECT SUM(CASE WHEN status = "APPROVED" THEN amount ELSE 0 END) AS sta, SUM(CASE WHEN status = "PENDING" THEN amount ELSE 0 END) AS pta, SUM(CASE WHEN status="DECLINED" THEN amount ELSE 0 END) AS dta, DATE(transaction_time) as df, MONTH(transaction_time) as mn, YEAR(transaction_time) AS y, date_format(transaction_time, '%b') AS month, date_format(transaction_time, '%b-%Y') AS ym, date_format(transaction_time, '%Y-%b') AS yM, date_format(transaction_time, '%e-%b') AS dM, date_format(transaction_time, '%e') AS d, DATE_FORMAT(transaction_time, "%Y-%m-%d %H:%i") as hm FROM Transactions WHERE status="APPROVED" ` + group_by)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Db selection error")
		return
	}

	_, data, err := sql.Scan(rows)
	defer sql.Close(rows)
	if err != nil {
		err = errors.New("Columns not found")
		return
	}
	log.Println("Debug", "Info", "data", data)
	for _, v := range data {
		approvedAmount, _ := strconv.Atoi(v[0])
		if mode == "Y" {
			year := v[5]
			approved_map[year] = approvedAmount
		} else if mode == "MY" {
			month_year := v[7]
			approved_map[month_year] = approvedAmount
		} else if mode == "M" {
			month := v[6]
			approved_map[month] = approvedAmount
		} else if mode == "DM" {
			day_month := v[9]
			approved_map[day_month] = approvedAmount
		} else if mode == "D" {
			day, _ := strconv.Atoi(v[10])
			approved_map[day] = approvedAmount
		} else if mode == "HM" {
			day := v[11]
			approved_map[day] = approvedAmount
		}
	}

	var approved_transaction []int
	for _, v := range months {
		approved_transaction = append(approved_transaction, approved_map[v])
	}
	finalMonths := translateMonth(mode, months)

	finalData := make(map[string]interface{})
	finalData["approved_transaction"] = approved_transaction
	finalData["months"] = finalMonths
	finalData["approved_title"] = approvedLabel
	finalData["transaction_amount_detail_label"] = transactionAmountChartTitle
	finalData["total_amount_label"] = totalAmountLabel
	fmt.Println("Final data: ", finalData)
	c.Data["json"] = finalData
	c.ServeJSON()

}

func (c *ChartController) Analytics() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
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
		if err != nil && err.Error() == "Session Time Out.Please Logout and Login Again." {
			c.Data["DisplayMessage"] = err.Error()
			c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Fail")
		} else if err != nil {
			c.Data["DisplayMessage"] = err.Error()
			c.TplName = "error/error.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "dashboard/analytics.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Success")
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
		return
	}
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "Dashboard")
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

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "UserName - ", sess.Get("uname"))
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Role - ", sess.Get("role"))

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Session ID - ", sess.SessionID())
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()

	c.Data["name"] = sess.Get("name").(string)
	c.Data["role"] = sess.Get("role").(string)
	c.Data["Uname"] = sess.Get("uname").(string)
	c.Data["language"] = sess.Get("language").(string)
	c.Data["department"] = sess.Get("department").(string)

	c.Data["MenuJson"] = sess.Get("menujson")

	c.Data["operatorsCount"] = len(strings.Split(beego.AppConfig.String("CNPS_OPERATORS"), "|"))
	c.Data["Photo"] = sess.Get("photo").(string)
	c.Data["ChannelsCount"] = len(strings.Split(beego.AppConfig.String("CNPS_CHANNELS"), "|"))
	successAmount, servicecharge, totalTransCount, totalBanks, successCount, pendingCount, declainedCount, transErr := utils.SideBarTransactionData()
	if transErr != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get sidebar data")
		return
	}
	currentTime := time.Now()
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
		c.Data["TotalNumberAdminUsers"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_ADMIN_USERS")
		c.Data["TotalNumberMerchantUsers"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_MERCHNAT_USERS")

		c.Data["TotalNumberOfTransactions"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_TRANSACTIONS")
		c.Data["TotalAmountofTransactions"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_TRANSACTIONS_AMOUNT")
		c.Data["TotalSuccessfulTransactions"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_SUCCESSFUL_TRANSACTIONS")
		c.Data["TotalAmountofSuccessfulTransactions"] = beego.AppConfig.String("ENGLISH_TOTALAMOUNT_OF_SUCCESSFUL_TRANSACTIONS")
		c.Data["TotalDeclinedTransactions"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_DECLINED_TRANSACTIONS")
		c.Data["TotalAmountofDeclinedTransactions"] = beego.AppConfig.String("ENGLISH_TOTALAMOUNT_OF_DECLINED_TRANSACTIONS")
		c.Data["UserDetails"] = beego.AppConfig.String("ENGLISH_USER_DETAILS")
		c.Data["TrandDetails"] = beego.AppConfig.String("ENGLISH_TRANS_DETAILS")
		c.Data["Dashboard2"] = beego.AppConfig.String("ENGLISH_DASHBOARD")
		c.Data["TotalPendingTransactions"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_PENDING_TRANSACTIONS")
		c.Data["TotalAmountofPendingTransactions"] = beego.AppConfig.String("ENGLISH_TOTALAMOUNT_OF_PENDING_TRANSACTIONS")
		c.Data["TotalInitiatedTransactions"] = beego.AppConfig.String("ENGLISH_TOTALNUMBER_OF_INITIATED_TRANSACTIONS")
		c.Data["TotalAmountofInitiatedTransactions"] = beego.AppConfig.String("ENGLISH_TOTALAMOUNT_OF_INITIATED_TRANSACTIONS")

		//New template words
		c.Data["UserTransactionDetails"] = strings.Replace(beego.AppConfig.String("ENGLISH_USER_TRANSACTION_DETAILS"), "{{.curdate}}", currentTime.Format("2006-Jan-02"), -1)
		c.Data["TransactionDetails"] = beego.AppConfig.String("ENGLISH_TRANSACTION_DETAILS")
		c.Data["TransactionAmountDetails"] = beego.AppConfig.String("ENGLISH_TRANSACTION_AMOUNT_DETAILS")
		c.Data["Aproved"] = beego.AppConfig.String("ENGLISH_APPROVED")
		c.Data["Pending"] = beego.AppConfig.String("ENGLISH_PENDING")
		c.Data["Declained"] = beego.AppConfig.String("ENGLISH_DECLINED")
		c.Data["TotalTransactions"] = beego.AppConfig.String("ENGLISH_TOTAL_TRANSACTIONS")
		c.Data["SuccessfulTransactions"] = beego.AppConfig.String("ENGLISH_SUCCESSFUL_TRANSACTIONS")
		c.Data["DeclainedTransactions"] = beego.AppConfig.String("ENGLISH_DECLAINED_TRANSACTIONS")
		c.Data["CnpsPortalAdminUsers"] = beego.AppConfig.String("ENGLISH_CNPS_PORTAL_ADMIN_USERS")
		c.Data["CnpsPortalPaymentOrganizations"] = beego.AppConfig.String("ENGLISH_CNPS_PORTAL_PAYMENT_ORGANIZATIONS")
		c.Data["TransactioByOperatorDetails"] = beego.AppConfig.String("ENGLISH_TRANSACTION_BY_OPERATOR")
		c.Data["TransactionByChannelDetails"] = beego.AppConfig.String("ENGLISH_TRANSACTION_BY_CHANNEL")
		c.Data["TransactionStatusDetails"] = beego.AppConfig.String("ENGLISH_TRANSACTION_STATUS_DETAILS")
		c.Data["TransactionCountDetails"] = beego.AppConfig.String("ENGLISH_TRANSACTION_COUNT_DETAILS")
		c.Data["Search"] = beego.AppConfig.String("ENGLISH_SEARCH")
		c.Data["Reset"] = beego.AppConfig.String("ENGLISH_RESET")
		c.Data["DateRange"] = beego.AppConfig.String("ENGLISH_USER_SELECT_DATARANGE")
		c.Data["Apply"] = beego.AppConfig.String("ENGLISH_APPLY")
		c.Data["PleaseSelect"] = beego.AppConfig.String("ENGLISH_USER_SELECT_ROLE")
		c.Data["All"] = beego.AppConfig.String("ENGLISH_ALL")
		c.Data["SelectStatus"] = beego.AppConfig.String("ENGLISH_SELECT_STATUS")
		c.Data["TotalNumberChannels"] = beego.AppConfig.String("ENGLISH_TOTAL_NUMBER_OF_CHANNELS")
		c.Data["CnpsPortalPaymentChannels"] = beego.AppConfig.String("ENGLISH_CNPS_PORTAL_PAYMENT_CHANNELS")
		c.Data["ANALYTICS"] = beego.AppConfig.String("ENGLISH_ANALYTICS")

		//sess.Set("Menus", beego.AppConfig.String("MENU_TEMPLATE"))
		//sess.Set("Header", headerContent)
		c.TplName = "dashboard/analytics.html"
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
		c.Data["TotalNumberAdminUsers"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_ADMIN_USERS")
		c.Data["TotalNumberMerchantUsers"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_MERCHNAT_USERS")
		c.Data["TotalNumberOfTransactions"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_TRANSACTIONS")
		c.Data["TotalAmountofTransactions"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_TRANSACTIONS_AMOUNT")
		c.Data["TotalSuccessfulTransactions"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_SUCCESSFUL_TRANSACTIONS")
		c.Data["TotalAmountofSuccessfulTransactions"] = beego.AppConfig.String("FRENCH_TOTALAMOUNT_OF_SUCCESSFUL_TRANSACTIONS")
		c.Data["TrandDetails"] = beego.AppConfig.String("ENGLISH_TRANS_DETAILS")
		c.Data["TotalDeclinedTransactions"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_DECLINED_TRANSACTIONS")
		c.Data["TotalAmountofDeclinedTransactions"] = beego.AppConfig.String("FRENCH_TOTALAMOUNT_OF_DECLINED_TRANSACTIONS")
		c.Data["UserDetails"] = beego.AppConfig.String("FRENCH_USER_DETAILS")
		c.Data["TrandDetails"] = beego.AppConfig.String("FRENCH_TRANS_DETAILS")
		sess.Set("Menus", beego.AppConfig.String("FRENCH_MENU_TEMPLATE"))
		c.Data["Dashboard2"] = beego.AppConfig.String("FRENCH_DASHBOARD")
		c.Data["TotalPendingTransactions"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_PENDING_TRANSACTIONS")
		c.Data["TotalAmountofPendingTransactions"] = beego.AppConfig.String("FRENCH_TOTALAMOUNT_OF_PENDING_TRANSACTIONS")
		c.Data["TotalInitiatedTransactions"] = beego.AppConfig.String("FRENCH_TOTALNUMBER_OF_INITIATED_TRANSACTIONS")
		c.Data["TotalAmountofInitiatedTransactions"] = beego.AppConfig.String("FRENCH_TOTALAMOUNT_OF_INITIATED_TRANSACTIONS")

		//New Template words
		c.Data["UserTransactionDetails"] = strings.Replace(beego.AppConfig.String("FRENCH_USER_TRANSACTION_DETAILS"), "{{.curdate}}", currentTime.Format("2006-Jan-02"), -1)
		c.Data["TransactionDetails"] = beego.AppConfig.String("FRENCH_TRANSACTION_DETAILS")
		c.Data["TransactionAmountDetails"] = beego.AppConfig.String("FRENCH_TRANSACTION_AMOUNT_DETAILS")
		c.Data["Aproved"] = beego.AppConfig.String("FRENCH_APPROVED")
		c.Data["Pending"] = beego.AppConfig.String("FRENCH_PENDING")
		c.Data["Declained"] = beego.AppConfig.String("FRENCH_DECLINED")
		c.Data["TotalTransactions"] = beego.AppConfig.String("FRENCH_TOTAL_TRANSACTIONS")
		c.Data["SuccessfulTransactions"] = beego.AppConfig.String("FRENCH_SUCCESSFUL_TRANSACTIONS")
		c.Data["DeclainedTransactions"] = beego.AppConfig.String("FRENCH_DECLAINED_TRANSACTIONS")
		c.Data["CnpsPortalAdminUsers"] = beego.AppConfig.String("FRENCH_CNPS_PORTAL_ADMIN_USERS")
		c.Data["CnpsPortalPaymentOrganizations"] = beego.AppConfig.String("FRENCH_CNPS_PORTAL_PAYMENT_ORGANIZATIONS")
		c.Data["TransactioByOperatorDetails"] = beego.AppConfig.String("FRENCH_TRANSACTION_BY_OPERATOR")
		c.Data["TransactionByChannelDetails"] = beego.AppConfig.String("FRENCH_TRANSACTION_BY_CHANNEL")
		c.Data["TransactionStatusDetails"] = beego.AppConfig.String("FRENCH_TRANSACTION_STATUS_DETAILS")
		c.Data["TransactionCountDetails"] = beego.AppConfig.String("FRENCH_TRANSACTION_COUNT_DETAILS")
		c.Data["Search"] = beego.AppConfig.String("FRENCH_SEARCH")
		c.Data["Reset"] = beego.AppConfig.String("FRENCH_RESET")
		c.Data["DateRange"] = beego.AppConfig.String("FRENCH_USER_SELECT_DATARANGE")
		c.Data["Apply"] = beego.AppConfig.String("FRENCH_APPLY")
		c.Data["PleaseSelect"] = beego.AppConfig.String("FRENCH_USER_SELECT_ROLE")
		c.Data["All"] = beego.AppConfig.String("FRENCH_ALL")
		c.Data["SelectStatus"] = beego.AppConfig.String("FRENCH_SELECT_STATUS")
		c.Data["TotalNumberChannels"] = beego.AppConfig.String("FRENCH_TOTAL_NUMBER_OF_CHANNELS")
		c.Data["CnpsPortalPaymentChannels"] = beego.AppConfig.String("FRENCH_CNPS_PORTAL_PAYMENT_CHANNELS")
		c.Data["ANALYTICS"] = beego.AppConfig.String("FRENCH_ANALYTICS")

		//sess.Set("Header", headerContent)

		c.TplName = "dashboard/analytics.html"
	}
	row, err := db.Db.Query(`SELECT count(*)as acount,(SELECT count(*) FROM Users where role = 'MERCHANT') as mcount FROM Users where role = 'ADMIN'`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data1 - ", data[0][0])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data2 - ", data[0][1])

	c.Data["admincount"] = data[0][0]
	c.Data["merchantcount"] = data[0][1]

	// row1, err := db.Db.Query(`select count(JSON_EXTRACT(transaction_deatils, "$.Status")) as TCont,SUM(JSON_EXTRACT(transaction_deatils, "$.Amount")) as TSum,
	// (select count(JSON_EXTRACT(transaction_deatils, "$.Status")) FROM Transactions where JSON_EXTRACT(transaction_deatils, "$.Status") = "APPROVED") as ACont,
	// (select SUM(JSON_EXTRACT(transaction_deatils, "$.Amount")) FROM Transactions where JSON_EXTRACT(transaction_deatils, "$.Status") = "APPROVED") as ASum,
	// (select count(JSON_EXTRACT(transaction_deatils, "$.Status")) FROM Transactions where JSON_EXTRACT(transaction_deatils, "$.Status") = "DECLINED") as DCont,
	// (select SUM(JSON_EXTRACT(transaction_deatils, "$.Amount")) FROM Transactions where JSON_EXTRACT(transaction_deatils, "$.Status") = "DECLINED") as DSum
	// FROM Transactions`)

	row1, err := db.Db.Query(`select (select count(status) FROM Transactions where status = "APPROVED" OR status="PENDING" OR status="DECLINED") as TCont,sum(amount) as TSum,
	(select count(status) FROM Transactions where status = "APPROVED") as ACont,
	(select sum(amount) FROM Transactions where status = "APPROVED") as ASum,
	(select count(status) FROM Transactions where status = "DECLINED") as DCont,
	(select sum(amount) FROM Transactions where status = "DECLINED") as DSum,
	(select count(status) FROM Transactions where status = "PENDING") as PCont,
	(select sum(amount) FROM Transactions where status = "PENDING") as PSum,
	(select count(status) FROM Transactions where status = "INITIATED") as ICont,
	(select sum(amount) FROM Transactions where status = "INITIATED") as ISum
 	from Transactions`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row1)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row1)
	_, data1, err := sql.Scan(row1)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data1, "\nData len - ", len(data1))
	if len(data1) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data1 - ", data1[0][0])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data2 - ", data1[0][1])

	c.Data["totaltransactioncount"] = data1[0][0]
	c.Data["totaltransactionamount"] = data1[0][1]
	c.Data["totaltsuccessfulransactioncount"] = data1[0][2]
	c.Data["totaltsuccessfulransactionamount"] = data1[0][3]
	c.Data["toatldeclinedtransactioncount"] = data1[0][4]
	c.Data["toatldeclinedtransactionamount"] = data1[0][5]
	c.Data["toatlpendingtransactioncount"] = data1[0][6]
	c.Data["toatlpendingtransactionamount"] = data1[0][7]
	c.Data["toatlinitiatedtransactioncount"] = data1[0][8]
	c.Data["toatlinitiatedtransactionamount"] = data1[0][9]

	return
}
