package regexpx

import (
	"github.com/dlclark/regexp2"
	"regexp"
)

const (
	regPhone    = "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	regEmail    = "^([A-Za-z0-9_\\-\\.])+\\@([A-Za-z0-9_\\-\\.])+\\.([A-Za-z]{2,4})$"
	regUserName = "^([A-Za-z_])+\\w"
	regPassword = `^(?![0-9a-zA-Z]+$)(?![a-zA-Z~!@#$%^&*()_+]+$)(?![0-9~!@#$%^&*()_+]+$)[0-9A-Za-z~!@#$%^&*()_+]{6,20}$`
)

var (
	phoneRegex    = regexp.MustCompile(regPhone)
	usernameRegex = regexp.MustCompile(regUserName)
	emailRegex    = regexp.MustCompile(regEmail)
)

func RegPhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

func RegUserName(username string) bool {
	return usernameRegex.MatchString(username)
}

func RegEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func RegPassword(password string) bool {
	reg, _ := regexp2.Compile(regPassword, 0)
	m, _ := reg.FindStringMatch(password)
	if m != nil {
		return true
	}
	return false
}
