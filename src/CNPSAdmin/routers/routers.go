package routers

import (
	"CNPSAdmin/controllers/currency/changeCurrencyConversionMethod"
	"CNPSAdmin/controllers/currency/createCurrency"
	"CNPSAdmin/controllers/currency/searchCurrency"
	"CNPSAdmin/controllers/currency/updateCurrency"
	"CNPSAdmin/controllers/currency/viewCurrency"
	"CNPSAdmin/controllers/dashboard"
	"CNPSAdmin/controllers/merchantdashboard"

	"CNPSAdmin/controllers/packages"
	"CNPSAdmin/controllers/packages/searchPackage"
	"CNPSAdmin/controllers/packages/updatePackage"
	"CNPSAdmin/controllers/packages/viewPackage"
	"CNPSAdmin/controllers/switchManagement/updateswtichstatus"
	"CNPSAdmin/controllers/switchManagement/viewswtichstatus"
	"CNPSAdmin/controllers/users/createUser"
	"CNPSAdmin/controllers/users/searchUser"
	"CNPSAdmin/controllers/users/updateUser"
	"CNPSAdmin/controllers/users/viewUser"

	"CNPSAdmin/controllers/reports/adminAduitReport"
	"CNPSAdmin/controllers/reports/channelReport"
	"CNPSAdmin/controllers/reports/refundReport"
	"CNPSAdmin/controllers/reports/settlementReport"
	transactionlogs "CNPSAdmin/controllers/reports/transactionLogs"
	"CNPSAdmin/controllers/reports/transactionReconReport"
	"CNPSAdmin/controllers/reports/transactionReport"

	"CNPSAdmin/controllers/channelManagement/searchChannel"
	"CNPSAdmin/controllers/channelManagement/updateChannel"
	"CNPSAdmin/controllers/channelManagement/viewChannel"

	"CNPSAdmin/controllers/adminChangePassword"
	"CNPSAdmin/controllers/adminViewProfile"
	forgotPassword "CNPSAdmin/controllers/forgotpassword"

	"CNPSAdmin/controllers/merchantChangePassword"
	"CNPSAdmin/controllers/merchantViewProfile"
	"CNPSAdmin/controllers/reports/merchantTransactionReport"

	"CNPSAdmin/controllers/login"
	"CNPSAdmin/controllers/logout"

	//"CNPSAdmin/controllers/ubaToken"

	"CNPSAdmin/controllers/commonChangePassword"

	"CNPSAdmin/controllers/role/createRole"
	"CNPSAdmin/controllers/role/searchRole"
	"CNPSAdmin/controllers/role/updateRole"
	"CNPSAdmin/controllers/role/viewRole"

	"CNPSAdmin/controllers/merchant/searchMerchant"
	"CNPSAdmin/controllers/merchant/updateMerchant"
	"CNPSAdmin/controllers/privilege/searchPrivilege"
	"CNPSAdmin/controllers/privilege/updatePrivilege"
	"CNPSAdmin/controllers/privilege/viewPrivilege"

	"CNPSAdmin/controllers/error"
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/recondb"
	"CNPSAdmin/session"

	"CNPSAdmin/controllers/charts"

	searchtransaction "CNPSAdmin/controllers/processing/transactionProcessing/searchTransaction"
	updatetransaction "CNPSAdmin/controllers/processing/transactionProcessing/updateTransaction"

	searchSettlement "CNPSAdmin/controllers/processing/settlementProcessing/searchSettlement"

	searchRefund "CNPSAdmin/controllers/processing/refundProcessing/searchRefund"

	"github.com/astaxie/beego"
)

func init() {
	if session.Init() != nil {
		return
	}
	if db.Init() != nil {
		return
	}
	if recondb.Init() != nil {
		return
	}
	beego.SetStaticPath("/upload", "upload")
	beego.ErrorController(&error.Error{})
	beego.Router(beego.AppConfig.String("MAIN_PATH"), &login.Login{})
	beego.Router(beego.AppConfig.String("PACKAGE_PATH"), &packages.PackagesController{})
	beego.Router(beego.AppConfig.String("CNPS_SEACRH_PACKAGE_PATH"), &searchPackage.SearchPackage{})
	beego.Router(beego.AppConfig.String("CNPS_VIEW_PACKAGE_PATH"), &viewPackage.ViewPackage{})
	beego.Router(beego.AppConfig.String("CNPS_UPDATE_PACKAGE_PATH"), &updatePackage.UpdatePackage{})
	beego.Router(beego.AppConfig.String("LOGIN_PATH"), &login.Login{})
	beego.Router(beego.AppConfig.String("LOGOUT_PATH"), &logout.Logout{})
	beego.Router(beego.AppConfig.String("DASHBOARD_PATH"), &dashboard.Dashboard{})
	beego.Router(beego.AppConfig.String("CNPS_SEACRH_USER_PATH"), &searchUser.SearchUser{})
	beego.Router(beego.AppConfig.String("CNPS_VIEW_USER_PATH"), &viewUser.ViewUser{})
	beego.Router(beego.AppConfig.String("CNPS_CREATE_USER_PATH"), &createUser.CreateUser{})
	beego.Router(beego.AppConfig.String("CNPS_UPDATE_USER_PATH"), &updateUser.UpdateUser{})

	beego.Router(beego.AppConfig.String("CNPS_VIEW_SWTICHSTATUS_PATH"), &viewswtichstatus.Viewswtichstatus{})
	beego.Router(beego.AppConfig.String("CNPS_UPDATE_SWTICHSTATUS_PATH"), &updateswtichstatus.Updateswtichstatus{})
	beego.Router(beego.AppConfig.String("CNPS_ADMIN_ADUIT_REPORT_PATH"), &adminAduitReport.AdminAduitReport{})
	beego.Router(beego.AppConfig.String("CNPS_CHANNEL_REPORT_PATH"), &channelReport.ChannelReport{})
	beego.Router(beego.AppConfig.String("CNPS_PG_TRANSACTION_REPORT_PATH"), &transactionReport.TransactionReport{})
	beego.Router(beego.AppConfig.String("CNPS_PG_SETTLEMENT_REPORT_PATH"), &settlementReport.SettlementReport{})
	beego.Router(beego.AppConfig.String("CNPS_PG_REFUND_REPORT_PATH"), &refundReport.RefundReport{})
	beego.Router(beego.AppConfig.String("CNPS_PG_TRANSACTION_RECON_REPORT_PATH"), &transactionReconReport.TransactionReconReport{})

	beego.Router(beego.AppConfig.String("CNPS_ADMIN_CHANGE_PASSWORD"), &adminChangePassword.AdminChangePassword{})
	beego.Router(beego.AppConfig.String("CNPS_ADMIN_VIEW_PROFILE"), &adminViewProfile.AdminViewProfile{})
	beego.Router(beego.AppConfig.String("CNPS_FORGOT_PASSWORD"), &forgotPassword.ForgotPassword{})

	beego.Router(beego.AppConfig.String("CNPS_SEACRH_CHANNEL_PATH"), &searchChannel.SearchChannel{})
	beego.Router(beego.AppConfig.String("CNPS_VIEW_CHANNEL_PATH"), &viewChannel.ViewChannel{})
	beego.Router(beego.AppConfig.String("CNPS_UPDATE_CHANNEL_PATH"), &updateChannel.UpdateChannel{})

	beego.Router(beego.AppConfig.String("CNPS_MERCHANT_VIEW_PROFILE"), &merchantViewProfile.MerchantViewProfile{})
	beego.Router(beego.AppConfig.String("CNPS_MERCHANT_CHANGE_PASSWORD"), &merchantChangePassword.MerchantChangePassword{})

	beego.Router(beego.AppConfig.String("CNPS_MERCHANT_TRANSACTION_REPORT_PATH"), &merchantTransactionReport.MerchantTransactionReport{})

	beego.Router(beego.AppConfig.String("CNPS_COMMON_CHANGEPASSWORD_PATH"), &commonChangePassword.CommonChangePassword{})

	beego.Router(beego.AppConfig.String("MERCHANT_DASHBOARD_PATH"), &merchantdashboard.MerchantDashboard{})

	beego.Router(beego.AppConfig.String("CNPS_ROLE_SEARCHROLE_PATH"), &searchRole.SearchRole{})
	beego.Router(beego.AppConfig.String("CNPS_ROLE_CREATEROLE_PATH"), &createRole.CreateRole{})
	beego.Router(beego.AppConfig.String("CNPS_ROLE_VIEWROLE_PATH"), &viewRole.ViewRole{})
	beego.Router(beego.AppConfig.String("CNPS_ROLE_UPDATEROLE_PATH"), &updateRole.UpdateRole{})
	beego.Router(beego.AppConfig.String("CNPS_ROLE_NAME_UNIQUE_PATH"), &createRole.CreateRole{})

	beego.Router(beego.AppConfig.String("CNPS_PRIVILEGE_SEARCHPRIVILEGE_PATH"), &searchPrivilege.SearchPrivilege{})
	beego.Router(beego.AppConfig.String("CNPS_PRIVILEGE_VIEWPRIVILEGE_PATH"), &viewPrivilege.ViewPrivilege{})
	beego.Router(beego.AppConfig.String("CNPS_PRIVILEGE_UPDATEPRIVILEGE_PATH"), &updatePrivilege.UpdatePrivilege{})
	beego.Router(beego.AppConfig.String("CNPS_TRANSACTION_AMOUNT_STATUS_PATH"), &charts.ChartController{}, "post:TransactionAmountByStatus")
	beego.Router(beego.AppConfig.String("CNPS_TRANSACTION_COUNT_STATUS_PATH"), &charts.ChartController{}, "post:TransactionCountByStatus")
	beego.Router(beego.AppConfig.String("CNPS_TRANSACTION_BY_OPERATOR_PATH"), &charts.ChartController{}, "post:TransactionByOperator")
	beego.Router(beego.AppConfig.String("CNPS_TRANSACTION_BY_CHANNEL_PATH"), &charts.ChartController{}, "post:TransactionByChannel")
	beego.Router(beego.AppConfig.String("CNPS_TRANSACTION_LIVE_PATH"), &charts.ChartController{}, "post:TransactionLiveData")
	beego.Router(beego.AppConfig.String("CNPS_REMOVE_PROFILE_PHOTO_PATH"), &adminViewProfile.AdminViewProfile{}, "get:RemoveProfilePhoto")
	beego.Router(beego.AppConfig.String("ANALYTICS_PATH"), &charts.ChartController{}, "get:Analytics")

	beego.Router(beego.AppConfig.String("PGS_MERCHANT_MERCHANT_APPROVAL"), &searchMerchant.SearchMerchant{})
	beego.Router(beego.AppConfig.String("PGS_MERCHANT_MERCHANT_APPROVAL_UPDATE"), &updateMerchant.UpdateMerchant{})

	// beego.Router(beego.AppConfig.String("UBA_TOKEN_PATH"), &ubaToken.UbaToken{})

	beego.Router(beego.AppConfig.String("CNPS_TRANSACTION_PROCESSING_PATH"), &searchtransaction.Searchtransaction{})
	beego.Router(beego.AppConfig.String("CNPS_TRANSACTION_UPDATE_PATH"), &updatetransaction.Updatetransaction{})
	beego.Router(beego.AppConfig.String("PGS_SEARCH_SETTLEMENT"), &searchSettlement.SearchSettlement{})
	beego.Router(beego.AppConfig.String("PGS_APPROVE_SETTLEMENT"), &searchSettlement.SearchSettlement{}, "get:ApproveSettlement")
	beego.Router(beego.AppConfig.String("PGS_SEARCH_REFUND"), &searchRefund.SearchRefund{})
	beego.Router(beego.AppConfig.String("PGS_APPROVE_REFUND"), &searchRefund.SearchRefund{}, "get:ApproveRefund")
	beego.Router(beego.AppConfig.String("TRANSACTION_LOGS_EXPORT_PATH"), &transactionlogs.TransactionLogs{})

	beego.Router(beego.AppConfig.String("PGS_CHANGE_CURRENCY_CONVERSION_PATH"), &changeCurrencyConversionMethod.ChangeCurrencyConversionMethod{})
	beego.Router(beego.AppConfig.String("PGS_SEACRH_CURRENCY_PATH"), &searchCurrency.SearchCurrency{})
	beego.Router(beego.AppConfig.String("PGS_CREATE_CURRENCY_PATH"), &createCurrency.CreateCurrency{})
	beego.Router(beego.AppConfig.String("PGS_VIEW_CURRENCY_PATH"), &viewCurrency.ViewCurrency{})
	beego.Router(beego.AppConfig.String("PGS_UPDATE_CURRENCY_PATH"), &updateCurrency.UpdateCurrency{})

}
