package message

import (
	"cloudProject/statics/config"
	"cloudProject/statics/translate/errCode"
)

func GetMessage(key, lang string) string {
	lang = GetValidLang(lang)
	errCodeKey := errCode.GetErrCodeKey(key)
	if errCodeKey != "" {
		key = errCodeKey
	}
	var message string
	switch lang {
	case "fa":
		message = persianMessage[key]
	case "en":
		message = englishMessage[key]
	}
	return message
}
func getDefaultMessage(lang string) string {
	switch lang {
	case "fa":
		return persianMessage["errorOccurred"]
	case "en":
		return englishMessage["errorOccurred"]
	}
	return ""
}
func GetValidLang(lang string) string {
	if !config.AcceptedLang[lang] {
		lang = config.DefaultLang
	}
	return lang
}
