package response

type QueryParams struct {
	Username  string `form:"username"`
	Major     string `form:"major"`
	Grade     string `form:"grade"`
	Direction string `form:"direction"`
	Role      string `form:"role"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"pageSize,default=6"`
}
