package request

type RequestInterface interface {
	MakeRequest()
	handlePayload()
	handleQuery()
}
