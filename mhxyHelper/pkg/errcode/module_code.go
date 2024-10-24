package errcode

var (
	// Stuff 相关 error
	ErrorBuildStuffByStrFail = NewError(20020001, "构建物品信息失败")
	ErrorQueryStuffFail      = NewError(20020002, "查询物品信息失败")

	// Attribute 相关 error
	ErrorBuildAttributeByStrFail = NewError(20020001, "构建属性信息失败")
	ErrorQueryAttributeFail      = NewError(20020002, "查询属性信息失败")
)
