package error

var (
	SUCCESS = Template{200, "Success"}
	FAILURE = Template{500, "服务器繁忙"}
)

type Template struct {
	Code int
	Msg  string
}
