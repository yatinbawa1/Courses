package mailer

import (
	"courses/internal/config"

	"github.com/resend/resend-go/v3"
)

func MailConfig() *resend.Client {
	mailClient := resend.NewClient(config.EMAIL_API_KEY)
	return mailClient
}

func SendOTPMail()
