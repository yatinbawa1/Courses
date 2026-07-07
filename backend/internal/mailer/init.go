package mailer

import (
	"context"
	"courses/internal/config"
	"courses/internal/models"
	"fmt"

	"github.com/resend/resend-go/v3"
)

type MailSender interface {
	SendOTPMail(ctx context.Context, user *models.User, otp string) error
}

type ResendMailer struct {
	mailClient *resend.Client
}

func NewResendMailer() *ResendMailer {
	client := resend.NewClient(config.EMAIL_API_KEY)
	return &ResendMailer{
		mailClient: client,
	}
}

func (m *ResendMailer) SendOTPMail(ctx context.Context, user *models.User, otp string) error {

	hmx, err := GenerateEmailHTML(user.Email, otp)
	if err != nil {
		return fmt.Errorf("Unable to Generate HTML for Sending OTP %w", err)
	}

	params := &resend.SendEmailRequest{
		From:    "yatinabc3@gmail.com",
		To:      []string{user.Email},
		Subject: "Your Verification Code",
		Html:    hmx,
	}

	_, err = m.mailClient.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	return nil
}
