package utils

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func FinalizeHTMl(html string, replaceMap map[string]string) string {
	resultHtml := html
	for k, v := range replaceMap {
		resultHtml = strings.ReplaceAll(resultHtml, k, v)
	}

	return resultHtml
}

func GenerateUrl(host, tag, path string) string {
	if path == "" {
		return host + "/" + tag
	}
	return host + "/" + path + "/" + tag
}

func GenerateResetPasswordMail(identity, name, tag, host, templatePath string) string {
	replaceMap := make(map[string]string)
	link := GenerateUrl(host, tag, "resetpassword")
	replaceMap["@@LINK"] = link
	replaceMap["@@USER"] = identity
	replaceMap["@@NAME"] = name

	html := GetTemplateFromFile(templatePath)
	finalHtml := FinalizeHTMl(html, replaceMap)

	return finalHtml
}

func GenerateVerifyEmailMail(email, name, tag, host, templatePath string) string {
	replaceMap := make(map[string]string)
	link := GenerateUrl(host, tag, "verify")
	replaceMap["@@EMAIL"] = email
	replaceMap["@@NAME"] = name
	replaceMap["@@LINK"] = link

	html := GetTemplateFromFile(templatePath)
	finalHtml := FinalizeHTMl(html, replaceMap)

	return finalHtml
}

func GetTemplateFromFile(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return ""
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return ""
	}
	return string(data)
}
