// Package mail ports services/email.py's SMTP sending (Gmail over implicit
// TLS on 465). Only the worker imports this — the API only ever publishes an
// EmailJob, so a slow/blocked SMTP call can never stall a request.
package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

const smtpHost = "smtp.gmail.com"
const smtpAddr = "smtp.gmail.com:465"

// Sender holds the SMTP credentials (EMAIL_USER / EMAIL_PASS).
type Sender struct {
	User string
	Pass string
}

// New builds a Sender. Send calls no-op with an error if user/pass are empty,
// matching local-dev behaviour where email isn't configured.
func New(user, pass string) *Sender {
	return &Sender{User: user, Pass: pass}
}

func (s *Sender) send(to, subject, htmlBody string) error {
	if s.User == "" || s.Pass == "" {
		return fmt.Errorf("email not configured (EMAIL_USER/EMAIL_PASS unset)")
	}

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		s.User, to, subject, htmlBody,
	)

	conn, err := tls.Dial("tcp", smtpAddr, &tls.Config{ServerName: smtpHost})
	if err != nil {
		return fmt.Errorf("smtp dial: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Close()

	auth := smtp.PlainAuth("", s.User, s.Pass, smtpHost)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}
	if err := client.Mail(s.User); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = w.Write([]byte(msg))
	return err
}

// SendWelcome ports send_welcome_email.
func (s *Sender) SendWelcome(to, name, password, role string) error {
	html := fmt.Sprintf(`<h2>Welcome to BitwiseLearn!</h2>
<p>Hello %s,</p>
<p>Your %s account has been created successfully.</p>
<p><strong>Email:</strong> %s</p>
<p><strong>Password:</strong> %s</p>
<p>Please change your password after your first login.</p>
<p>Best regards,<br>BitwiseLearn Team</p>`, name, role, to, password)
	return s.send(to, fmt.Sprintf("Welcome to BitwiseLearn - %s Account", role), html)
}

// SendOTP ports send_otp_email.
func (s *Sender) SendOTP(to, otp string) error {
	html := fmt.Sprintf(`<h2>Password Reset OTP</h2>
<p>Your OTP for password reset is:</p>
<h1 style="color: #4F46E5; letter-spacing: 4px;">%s</h1>
<p>This OTP is valid for 10 minutes.</p>
<p>If you didn't request this, please ignore this email.</p>
<p>Best regards,<br>BitwiseLearn Team</p>`, otp)
	return s.send(to, "BitwiseLearn - Password Reset OTP", html)
}

// SendContact ports send_contact_email: notifies the site's own inbox
// (EMAIL_USER) about a contact-form submission; name/email/message describe
// the visitor, not the recipient.
func (s *Sender) SendContact(visitorName, visitorEmail, message string) error {
	html := fmt.Sprintf(`<h2>New Contact Form Submission</h2>
<p><strong>Name:</strong> %s</p>
<p><strong>Email:</strong> %s</p>
<p><strong>Message:</strong> %s</p>`, visitorName, visitorEmail, message)
	return s.send(s.User, fmt.Sprintf("Contact Form: %s", visitorName), html)
}
