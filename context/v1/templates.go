package v1

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/axiomzen/zenauth/config"
	"github.com/axiomzen/zenauth/models"
)

var changePasswordHTMLTmpl *template.Template
var generalMessageHTMLTmpl *template.Template

func init() {
	conf, err := config.Get()
	if err != nil {
		panic(err)
	}
	templates, err := template.ParseGlob(filepath.Join(conf.HTMLTemplatesPath, "*.tmpl"))
	if err != nil {
		panic(err)
	}

	changePasswordHTMLTmpl = templates.Lookup("change_password.html.tmpl")
	if changePasswordHTMLTmpl == nil {
		panic(fmt.Errorf("Change password HTML template not found"))
	}
	generalMessageHTMLTmpl = templates.Lookup("general_message.html.tmpl")
	if generalMessageHTMLTmpl == nil {
		panic(fmt.Errorf("General message HTML template not found"))
	}
}

// GetChangePasswordHTML returns a Template instance for the reset password action
func GetChangePasswordHTML(conf *config.ZENAUTHConfig, user *models.User) (string, error) {
	if user.ResetToken == nil {
		return "", fmt.Errorf("User has not requested to reset password")
	}
	variables := map[string]string{
		"title":    "Select your new password",
		"token":    *user.ResetToken,
		"email":    user.Email,
		"URL":      conf.ResetPasswordURL,
		"redirect": conf.ResetPasswordRedirectURL,
	}
	bufHTML := &bytes.Buffer{}
	if err := changePasswordHTMLTmpl.Execute(bufHTML, variables); err != nil {
		return "", err
	}
	return bufHTML.String(), nil
}

// GetGeneralMessageHTML returns a Template instance for the reset password action
func GetGeneralMessageHTML(message string) (string, error) {
	variables := map[string]string{
		"title":   message,
		"message": message,
	}
	bufHTML := &bytes.Buffer{}
	if err := generalMessageHTMLTmpl.Execute(bufHTML, variables); err != nil {
		return "", err
	}
	return bufHTML.String(), nil
}
