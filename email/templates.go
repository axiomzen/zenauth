package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"

	"github.com/axiomzen/zenauth/config"
	"github.com/axiomzen/zenauth/models"
)

var resetPasswordHTMLTmpl *template.Template
var resetPasswordTextTmpl *template.Template
var verifyEmailHTMLTmpl *template.Template
var verifyEmailTextTmpl *template.Template

func init() {
	conf, err := config.Get()
	if err != nil {
		panic(err)
	}
	templates, err := template.ParseGlob(filepath.Join(conf.TemplatesPath, "*.tmpl"))
	if err != nil {
		panic(err)
	}

	resetPasswordHTMLTmpl = templates.Lookup("reset_password.html.tmpl")
	if resetPasswordHTMLTmpl == nil {
		panic(fmt.Errorf("Reset password HTML template not found"))
	}
	resetPasswordTextTmpl = templates.Lookup("reset_password.txt.tmpl")
	if resetPasswordTextTmpl == nil {
		panic(fmt.Errorf("Reset password TEXT template not found"))
	}
	verifyEmailHTMLTmpl = templates.Lookup("verify_email.html.tmpl")
	if verifyEmailHTMLTmpl == nil {
		panic(fmt.Errorf("Verify email HTML template not found"))
	}
	verifyEmailTextTmpl = templates.Lookup("verify_email.txt.tmpl")
	if verifyEmailTextTmpl == nil {
		panic(fmt.Errorf("Verify email TEXT template not found"))
	}
}

// GetResetPasswordMessage returns a Message instance for the reset password action
func GetResetPasswordMessage(conf *config.ZENAUTHConfig, user *models.User) (*Message, error) {
	message := Message{}
	message.Subject = fmt.Sprintf("[%v] Reset Password", conf.AppName)
	message.From = fmt.Sprintf("%v <%v>", conf.AppName, conf.MailGunFrom)
	message.To = []string{user.Email}

	resetURL, err := url.Parse(conf.ResetPasswordURL)
	if err != nil {
		return nil, err
	}
	query := resetURL.Query()
	query.Add("token", *user.ResetToken)
	query.Add("email", user.Email)
	resetURL.RawQuery = query.Encode()

	variables := map[string]string{
		"title": message.Subject,
		"URL":   resetURL.String(),
	}

	bufHTML := &bytes.Buffer{}
	if err := resetPasswordHTMLTmpl.Execute(bufHTML, variables); err != nil {
		return nil, err
	}
	message.BodyHTML = bufHTML.String()

	bufText := &bytes.Buffer{}
	if err := resetPasswordTextTmpl.Execute(bufText, variables); err != nil {
		return nil, err
	}
	message.Body = bufText.String()

	return &message, nil
}

// GetVerifyEmailMessage returns a Message instance for the reset password action
func GetVerifyEmailMessage(conf *config.ZENAUTHConfig, user *models.User) (*Message, error) {
	message := Message{}
	message.Subject = fmt.Sprintf("[%v] Verify Email", conf.AppName)
	message.From = fmt.Sprintf("%v <%v>", conf.AppName, conf.MailGunFrom)
	message.To = []string{user.Email}

	resetURL, err := url.Parse(conf.VerifyEmailURL)
	if err != nil {
		return nil, err
	}
	query := resetURL.Query()
	query.Add("token", user.VerifyEmailToken)
	query.Add("email", user.Email)
	resetURL.RawQuery = query.Encode()

	variables := map[string]string{
		"title": message.Subject,
		"URL":   resetURL.String(),
	}

	bufHTML := &bytes.Buffer{}
	if err := verifyEmailHTMLTmpl.Execute(bufHTML, variables); err != nil {
		return nil, err
	}
	message.BodyHTML = bufHTML.String()

	bufText := &bytes.Buffer{}
	if err := verifyEmailTextTmpl.Execute(bufText, variables); err != nil {
		return nil, err
	}
	message.Body = bufText.String()

	return &message, nil
}
