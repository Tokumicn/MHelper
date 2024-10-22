package errcode

var (
	// Stuff 相关 error

	ErrorBuildHelperFail = NewError(20020001, "构建帮助信息失败")
	ErrorGetHelperFail   = NewError(20020002, "查询帮助信息失败")
)
