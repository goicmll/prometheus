package prometheus

type PromError struct {
	Msg string
}

// PromError.Error 实现了 error 接口的 Error 方法。
// 该方法返回错误消息，使得 PromError 类型的实例可以被用作错误处理。
//
// 返回值:
// string - 错误消息，提供了关于错误的详细信息。
func (e PromError) Error() string {
	return e.Msg
}
