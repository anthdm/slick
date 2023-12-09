package main

import (
	"fmt"

	"github.com/anthdm/slick"
	"github.com/anthdm/slick/app/view/dashboard"
	"github.com/anthdm/slick/app/view/profile"
	"github.com/google/uuid"
)

func main() {
	app := slick.New()

	app.Plug(WithRequestID, WithAuth)

	ph := NewProfileHandler(&NOOPSB{})

	app.Get("/profile", ph.HandleProfileIndex)
	app.Get("/dashboard", HandleDashboardIndex)

	app.Start(":3000")
}

func WithAuth(h slick.Handler) slick.Handler {
	return func(c *slick.Context) error {
		fmt.Println("auth")
		c.Set("email", "jannine@hr.com")
		return h(c)
	}
}

func WithRequestID(h slick.Handler) slick.Handler {
	return func(c *slick.Context) error {
		fmt.Println("request")
		c.Set("requestID", uuid.New())
		return h(c)
	}
}

type SupabaseClient interface {
	Auth(foo string) error
	// ...
	//
}

type NOOPSB struct{}

func (NOOPSB) Auth(foo string) error { return nil }

type ProfileHandler struct {
	sbClient SupabaseClient
	//...
}

func NewProfileHandler(sb SupabaseClient) *ProfileHandler {
	return &ProfileHandler{
		sbClient: sb,
	}
}

func (h *ProfileHandler) HandleProfileIndex(c *slick.Context) error {
	user := profile.User{
		FirstName: "A",
		LastName:  "GG",
		Email:     "gg@gg.com",
	}
	return c.Render(profile.Index(user))
}

func HandleDashboardIndex(c *slick.Context) error {
	fmt.Println(c.Get("requestID"))
	return c.Render(dashboard.Index())
}
