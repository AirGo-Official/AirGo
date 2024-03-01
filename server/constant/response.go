package constant

const (
	RESPONSE_ERROR   = 1  //code=1，不能正常获取请求数据
	RESPONSE_SUCCESS = 0  //code=0，能正常获取请求数据
	RESPONSE_WARNING = 10 //code=10，能正常获取请求数据，但有重要message 需要显式提醒
	// 权限 401xx
	TOKENERROR = 40101 //token过期
	// 限流 402xx
	LIMITERROR = 40201 //限流
)
