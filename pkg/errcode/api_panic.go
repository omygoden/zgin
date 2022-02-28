package errcode

type ApiPanic struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
	//err  error       `json:"-"`
}
