package login

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"

	httputils "github.com/skpr/cli/internal/http"
)

// Server for responding to the oauth2 callback.
type Server struct {
	Callback string
	Response Response
	// Logout renders the logout page on the root path. The identity provider
	// only redirects to sign out URLs which have been registered, so we accept
	// the callback with or without the /logout path.
	Logout bool
	// CallbackTimeout for how long we wait for the callback.
	CallbackTimeout time.Duration
}

// Response from the oauth2 callback.
type Response struct {
	Code             string
	State            string
	Error            string
	ErrorDescription string
}

const (
	// ShutdownTimeout for draining connections once a callback is received.
	ShutdownTimeout = 5 * time.Second
	// ReadyTimeout for the callback server to start responding.
	ReadyTimeout = 30 * time.Second
	// DefaultCallbackTimeout for how long we wait for the callback, which
	// includes the developer signing in with the identity provider.
	DefaultCallbackTimeout = 5 * time.Minute
)

// ErrCallbackTimeout is returned when the identity provider did not send us a
// callback. This usually means the callback URL has not been registered with
// the identity provider.
var ErrCallbackTimeout = errors.New("callback was not received")

// Embed the entire directory.
//
//go:embed tmpl
var tmpl embed.FS

// NewServer for responding to oauth2 callbacks.
func NewServer(callback string) *Server {
	return &Server{
		Callback:        callback,
		CallbackTimeout: DefaultCallbackTimeout,
	}
}

// Run the server and wait for the callback.
func (s *Server) Run(ctx context.Context, ready context.CancelFunc) (Response, error) {
	router := mux.NewRouter()

	ctxShutdown, shutdown := context.WithCancel(ctx)

	if s.Logout {
		router.HandleFunc("/", s.handleLogoutCallback(shutdown)).Methods("GET")
	} else {
		router.HandleFunc("/", s.handleLoginCallback(shutdown)).Methods("GET")
	}

	router.HandleFunc("/logout", s.handleLogoutCallback(shutdown)).Methods("GET")
	router.HandleFunc("/readyz", s.handleReadyz).Methods("GET")

	if s.CallbackTimeout == 0 {
		s.CallbackTimeout = DefaultCallbackTimeout
	}

	addr, err := url.Parse(s.Callback)
	if err != nil {
		return s.Response, fmt.Errorf("failed to parse callback URL: %w", err)
	}

	if addr.Hostname() == "" || addr.Port() == "" {
		return s.Response, fmt.Errorf("callback URL must include a host and a port: %s", s.Callback)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", addr.Hostname(), addr.Port()),
		Handler: router,
	}

	group, _ := errgroup.WithContext(context.Background())

	group.Go(func() error {
		// The callback can include a path, so we build this from the host
		// rather than the callback URL itself.
		readyz := fmt.Sprintf("%s://%s/readyz", addr.Scheme, addr.Host)

		err := httputils.Wait(ctxShutdown, readyz, ReadyTimeout)
		if err != nil {
			// Without this the server would wait for a callback which is never
			// going to arrive, because we never opened the browser session.
			shutdown()

			return fmt.Errorf("failed to wait for server: %w", err)
		}

		ready()

		return nil
	})

	group.Go(func() error {
		<-ctxShutdown.Done()

		// We shutdown with a new context because the one we waited on has been
		// cancelled, which would make the server give up on draining the
		// connection which is still delivering our response.
		ctxTimeout, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer cancel()

		return srv.Shutdown(ctxTimeout)
	})

	group.Go(func() error {
		// Without this we would wait forever for a callback which is never
		// going to arrive eg. When the identity provider rejected our callback
		// URL because it has not been registered.
		select {
		case <-ctxShutdown.Done():
			return nil
		case <-time.After(s.CallbackTimeout):
			shutdown()

			return fmt.Errorf("%w from %s within %v", ErrCallbackTimeout, s.Callback, s.CallbackTimeout)
		}
	})

	group.Go(func() error {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	return s.Response, group.Wait()
}

// Helper function to handle the callback.
func (s *Server) handleLoginCallback(shutdown context.CancelFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Shutdown the server after the callback is received.
		defer shutdown()

		s.Response = Response{
			Code:             r.URL.Query().Get("code"),
			State:            r.URL.Query().Get("state"),
			Error:            r.URL.Query().Get("error"),
			ErrorDescription: r.URL.Query().Get("error_description"),
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		data := map[string]string{
			"title": "Skpr Login",
			"code":  s.Response.Code,
		}

		tmpl, err := template.New("base.html").ParseFS(tmpl, "tmpl/login_success.html", "tmpl/base.html")
		if err != nil {
			log.Println("Failed to render template:", err)
		}

		err = tmpl.ExecuteTemplate(w, "base.html", data)
		if err != nil {
			log.Println("Failed to execute template:", err)
		}
	}
}

// Helper function to handle the callback.
func (s *Server) handleLogoutCallback(shutdown context.CancelFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Shutdown the server after the callback is received.
		defer shutdown()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		t, err := template.New("base.html").ParseFS(tmpl, "tmpl/logout_success.html", "tmpl/base.html")
		if err != nil {
			log.Println("Failed to render template:", err)
		}

		data := map[string]string{
			"title": "Skpr Logout",
		}

		err = t.ExecuteTemplate(w, "base.html", data)
		if err != nil {
			log.Println("Failed to execute template:", err)
		}
	}
}

// Helper function to handle the readyz endpoint.
func (s *Server) handleReadyz(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Ready!")
}
