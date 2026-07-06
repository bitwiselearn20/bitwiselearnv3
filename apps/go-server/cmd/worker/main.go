// Command worker consumes the email job queue so no SMTP I/O ever runs on the
// API request path (scale-to-zero on Container Apps: min replicas 0, scaled
// by queue depth).
package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bitwiselearn/go-server/internal/config"
	"github.com/bitwiselearn/go-server/internal/jobs"
	"github.com/bitwiselearn/go-server/internal/services/mail"
	"github.com/bitwiselearn/go-server/internal/services/queue"
)

func main() {
	cfg := config.Load()
	sender := mail.New(cfg.EmailUser, cfg.EmailPass)

	msgs, conn, err := queue.Consume(cfg.MQClient, jobs.EmailQueue)
	if err != nil {
		log.Fatalf("queue consume failed: %v", err)
	}
	defer conn.Close()
	log.Printf("worker listening on queue %q", jobs.EmailQueue)

	go func() {
		for d := range msgs {
			var job jobs.EmailJob
			if err := json.Unmarshal(d.Body, &job); err != nil {
				log.Printf("bad email job payload: %v", err)
				continue
			}
			if err := dispatch(sender, job); err != nil {
				// Matches the legacy routers' try/except-pass around email
				// sends: a delivery failure doesn't retry-loop or crash the
				// worker, it's just logged.
				log.Printf("email job failed (kind=%s to=%s): %v", job.Kind, job.To, err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("worker shutting down")
}

func dispatch(sender *mail.Sender, job jobs.EmailJob) error {
	switch job.Kind {
	case jobs.EmailKindWelcome:
		return sender.SendWelcome(job.To, job.Name, job.Password, job.Role)
	case jobs.EmailKindOTP:
		return sender.SendOTP(job.To, job.OTP)
	case jobs.EmailKindContact:
		return sender.SendContact(job.Name, job.To, job.Message)
	default:
		log.Printf("unknown email job kind: %s", job.Kind)
		return nil
	}
}
