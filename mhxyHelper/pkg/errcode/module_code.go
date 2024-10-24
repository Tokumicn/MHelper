package errcode

var (
	// Stuff 相关 error
	ErrorBuildStuffByStrFail = NewError(20020001, "构建物品信息失败")
	ErrorQueryStuffFail      = NewError(20020002, "查询物品信息失败")

	// Attribute 相关 error
	ErrorBuildAttributeByStrFail = NewError(20030001, "构建属性信息失败")
	ErrorQueryAttributeFail      = NewError(20030002, "查询属性信息失败")

	// Account 相关 error
	ErrorBuildAccountByStrFail = NewError(20040001, "构建账单信息失败")
	ErrorQueryAccountFail      = NewError(20040002, "查询账单信息失败")
)
