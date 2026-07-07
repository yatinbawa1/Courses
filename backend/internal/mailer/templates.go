package mailer

import (
	"bytes"
	"embed"
	"text/template"
)

var templateFS embed.FS

func GenerateEmailHTML(email string, code string) (string, error) {

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
