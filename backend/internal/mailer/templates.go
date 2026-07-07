package mailer

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed templates
var templateFS embed.FS

func GenerateOTPEmail(email string, code string) (string, error) {

	tmpl, err := template.ParseFS(templateFS, "templates/otp.html")

	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	data := map[string]string{
		"code":  code,
		"Email": email,
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}

	return body.String(), nil
}
