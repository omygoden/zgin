package requests

type TestResponse struct {
	Id int64 `json:"id" form:"id" binding:"required" label:"Id"`
}
