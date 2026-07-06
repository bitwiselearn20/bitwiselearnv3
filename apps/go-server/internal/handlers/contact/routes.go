// Package contact ports apps/python-server/routers/contact.py to Go.
package contact

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/bitwiselearn/go-server/internal/jobs"
	"github.com/bitwiselearn/go-server/internal/response"
	"github.com/bitwiselearn/go-server/internal/services/queue"
)

// Deps holds the dependencies the contact handler needs.
type Deps struct {
	Publisher *queue.Publisher
}

type contactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// Register mounts the contact route on g (expected prefix "/api/v1/contact").
func Register(g *echo.Group, d Deps) {
	h := &handler{pub: d.Publisher}
	g.POST("", h.sendContact)
}

type handler struct {
	pub *queue.Publisher
}

func (h *handler) sendContact(c echo.Context) error {
	var body contactRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	if h.pub == nil {
		return response.OK(c, "Message sent successfully", nil)
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()
	if err := h.pub.Publish(ctx, jobs.EmailQueue, jobs.EmailJob{
		Kind: jobs.EmailKindContact, Name: body.Name, To: body.Email, Message: body.Message,
	}); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to send message", err.Error())
	}
	return response.OK(c, "Message sent successfully", nil)
}
