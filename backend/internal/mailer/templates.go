package mailer

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed templates
var templateFS embed.FS

var otpTmpl = template.Must(template.ParseFS(templateFS, "templates/otp.html"))

func GenerateOTPEmail(email string, code string) (string, error) {
	var body bytes.Buffer
	data := map[string]string{
		"code":  code,
		"Email": email,
	}

	if err := otpTmpl.Execute(&body, data); err != nil {
		return "", err
	}

	return body.String(), nil
}
