package event

import (
	"bytes"
	"context"
	"text/template"
	"time"

	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/platform/config"
	"github.com/mia0x75/halo/internal/platform/notify"
)

// mailRoutes returns the per-topic mail handlers. Behaviour mirrors the former
// per-event *MailSender functions: the template and recipient rules are
// preserved exactly. userRepo resolves the owner/reviewer for the execution
// events, which do not carry a User.
func mailRoutes(n notify.Notifier, userRepo domain.UserRepository) map[string]Handler {
	routes := map[string]Handler{}
	add := func(topic string, fn Handler) {
		if fn != nil {
			routes[topic] = fn
		}
	}

	add(EventTicketCreated, mailToActor(n, config.TplTicketCreated, func(ev Event) (models.User, any) {
		e := ev.(*TicketCreated)
		return e.User, e
	}))
	add(EventTicketUpdated, mailToActor(n, config.TplTicketUpdated, func(ev Event) (models.User, any) {
		e := ev.(*TicketUpdated)
		return e.User, e
	}))
	add(EventTicketRemoved, mailToActor(n, config.TplTicketRemoved, func(ev Event) (models.User, any) {
		e := ev.(*TicketRemoved)
		// Template renders the removal time.
		e.Ticket.UpdateAt = uint(time.Now().UTC().Unix())
		return e.User, e
	}))
	add(EventTicketScheduled, mailToTicketParties(n, config.TplTicketScheduled, func(ev Event) (models.Ticket, any) {
		e := ev.(*TicketScheduled)
		return e.Ticket, e
	}))
	add(EventTicketStatusPatched, ticketStatusMail(n))
	add(EventTicketExecuted, mailToExecParties(n, userRepo, config.TplTicketExecuted))
	add(EventTicketFailed, mailToExecParties(n, userRepo, config.TplTicketFailed))
	add(EventUserRegistered, mailToActor(n, config.TplUserRegistered, func(ev Event) (models.User, any) {
		e := ev.(*UserRegistered)
		return e.User, e
	}))
	add(EventPasswordUpdated, mailToActor(n, config.TplPasswordUpdated, func(ev Event) (models.User, any) {
		e := ev.(*PasswordUpdated)
		return e.User, e
	}))
	add(EventEmailUpdated, mailToActor(n, config.TplEmailUpdated, func(ev Event) (models.User, any) {
		e := ev.(*EmailUpdated)
		// The verification mail goes to the new address.
		e.User.Email = e.Email
		return e.User, e
	}))
	add(EventProfileUpdated, mailToActor(n, config.TplProfileUpdated, func(ev Event) (models.User, any) {
		e := ev.(*ProfileUpdated)
		return e.User, e
	}))
	add(EventUserCreated, mailToActor(n, config.TplUserCreated, func(ev Event) (models.User, any) {
		e := ev.(*UserCreated)
		return e.User, e
	}))
	add(EventCommentCreated, mailToTicketParties(n, config.TplCommentCreated, func(ev Event) (models.Ticket, any) {
		e := ev.(*CommentCreated)
		return e.Ticket, e
	}))
	add(EventCronCancelled, mailToTicketParties(n, config.TplCronCancelled, func(ev Event) (models.Ticket, any) {
		e := ev.(*CronCancelled)
		return e.Ticket, e
	}))

	return routes
}

// mailToActor sends to a single user carried by the event.
func mailToActor(n notify.Notifier, tpl config.Template, pick func(Event) (models.User, any)) Handler {
	return func(ctx context.Context, ev Event) error {
		user, data := pick(ev)
		subject, body := renderer(tpl, data)
		return sendTo(ctx, n, []*models.User{&user}, subject, body)
	}
}

// mailToTicketParties sends to the ticket owner and reviewer resolved from the
// user cache, matching the former behaviour.
func mailToTicketParties(n notify.Notifier, tpl config.Template, pick func(Event) (models.Ticket, any)) Handler {
	return func(ctx context.Context, ev Event) error {
		ticket, data := pick(ev)
		users := cache.UsersMap.Filter(func(elem *models.User) bool {
			return elem.UserID == ticket.UserID || elem.UserID == ticket.ReviewerID
		})
		subject, body := renderer(tpl, data)
		return sendTo(ctx, n, users, subject, body)
	}
}

// ticketStatusMail sends a status-specific template; only Closed / MrvFailure /
// Lgtm produce mail, and Closed goes to the actor alone.
func ticketStatusMail(n notify.Notifier) Handler {
	return func(ctx context.Context, ev Event) error {
		e := ev.(*TicketStatusPatched)
		users := cache.UsersMap.Filter(func(elem *models.User) bool {
			return elem.UserID == e.Ticket.UserID || elem.UserID == e.Ticket.ReviewerID
		})
		var tpl config.Template
		switch domain.TicketStatus(e.Ticket.Status) {
		case domain.TicketStatusClosed:
			users = []*models.User{&e.User}
			tpl = config.TplTicketClosed
		case domain.TicketStatusMrvFailure:
			tpl = config.TplTicketMrvFailure
		case domain.TicketStatusLgtm:
			tpl = config.TplTicketLgtm
		default:
			return nil
		}
		subject, body := renderer(tpl, e)
		return sendTo(ctx, n, users, subject, body)
	}
}

// mailToExecParties resolves the owner and reviewer through the user repository
// (not the cache), matching the former execution-mail behaviour.
func mailToExecParties(n notify.Notifier, userRepo domain.UserRepository, tpl config.Template) Handler {
	return func(ctx context.Context, ev Event) error {
		var ticket models.Ticket
		switch e := ev.(type) {
		case *TicketExecuted:
			ticket = e.Ticket
		case *TicketFailed:
			ticket = e.Ticket
		}
		users := make([]*models.User, 0, 2)
		for _, uid := range []uint{ticket.UserID, ticket.ReviewerID} {
			du, err := userRepo.GetByID(ctx, int64(uid))
			if err != nil {
				return err
			}
			users = append(users, bridge.UserFromDomain(du))
		}
		subject, body := renderer(tpl, ev)
		return sendTo(ctx, n, users, subject, body)
	}
}

// sendTo delivers the rendered message to each user.
func sendTo(ctx context.Context, n notify.Notifier, users []*models.User, subject, body string) error {
	if n == nil {
		return nil
	}
	var firstErr error
	for _, u := range users {
		if err := n.Send(ctx, notify.Message{ToName: u.Name, ToEmail: u.Email, Subject: subject, Body: body}); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// renderer looks up the template by UUID and renders subject and body.
func renderer(id config.Template, data interface{}) (string, string) {
	tpl := cache.TemplatesMap.Any(func(elem *models.Template) bool {
		return elem.UUID == string(id)
	})
	if tpl == nil {
		return "", ""
	}
	fmap := template.FuncMap{
		"formatDate": formatDate,
	}
	tplSubject := template.Must(template.New("subject").Parse(tpl.Subject))
	tplBody := template.Must(template.New("body").Funcs(fmap).Parse(tpl.Body))
	subject := new(bytes.Buffer)
	body := new(bytes.Buffer)
	_ = tplSubject.Execute(subject, data)
	_ = tplBody.Execute(body, data)
	return subject.String(), body.String()
}

// formatDate renders a unix timestamp for mail templates.
func formatDate(ts uint) string {
	return time.Unix(int64(ts), 0).Format("2006-01-02 15:04:05")
}
