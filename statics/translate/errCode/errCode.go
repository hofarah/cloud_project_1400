package errCode

var errCodeMessages = map[string]string{
	"400": "badRequest",
	"403": "forbidden",
	"500": "errorOccurred",
}

func GetErrCodeKey(errCode string) string {
	return errCodeMessages[errCode]
}
