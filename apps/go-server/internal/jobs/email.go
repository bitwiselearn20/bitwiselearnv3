// Package jobs defines the message contracts shared by the API (publisher)
// and the worker (consumer) so both sides agree on the wire format without
// importing each other.
package jobs

// EmailQueue is the durable queue name for all outbound email jobs.
const EmailQueue = "email_jobs"

// Email job kinds.
const (
	EmailKindWelcome = "welcome"
	EmailKindOTP     = "otp"
	EmailKindContact = "contact"
)

// EmailJob is published by the API and consumed by cmd/worker, which performs
// the actual SMTP send — keeping that blocking I/O off the request path
// (ports the fire-and-forget send_welcome_email/send_otp_email/
// send_contact_email calls in the legacy routers to an async queue per the
// rewrite plan).
//
// For EmailKindContact, To/Name/Message carry the *visitor's* submitted
// name/email/message — the actual mail recipient is the site's own
// EMAIL_USER inbox (matching send_contact_email, which notifies the site
// owner rather than emailing the visitor).
type EmailJob struct {
	Kind     string `json:"kind"`
	To       string `json:"to"`
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
	OTP      string `json:"otp,omitempty"`
	Message  string `json:"message,omitempty"`
}
