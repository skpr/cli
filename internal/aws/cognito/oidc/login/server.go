package login

import (
	"context"
	"embed"
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
}

// Response from the oauth2 callback.
type Response struct {
	Code             string
	State            string
	Error            string
	ErrorDescription string
}

// Embed the entire directory.
//
//go:embed tmpl
var tmpl embed.FS

// NewServer for responding to oauth2 callbacks.
func NewServer(callback string) *Server {
	return &Server{
		Callback: callback,
	}
}

// Run the server and wait for the callback.
func (s *Server) Run(ctx context.Context, ready context.CancelFunc) (Response, error) {
	router := mux.NewRouter()

	ctxShutdown, shutdown := context.WithCancel(ctx)

	router.HandleFunc("/", s.handleLoginCallback(shutdown)).Methods("GET")
	router.HandleFunc("/logout", s.handleLogoutCallback(shutdown)).Methods("GET")
	router.HandleFunc("/readyz", s.handleReadyz).Methods("GET")

	addr, err := url.Parse(s.Callback)
	if err != nil {
		return s.Response, fmt.Errorf("failed to parse callback URL: %w", err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", addr.Hostname(), addr.Port()),
		Handler: router,
	}

	group, _ := errgroup.WithContext(context.Background())

	group.Go(func() error {
		err := httputils.Wait(fmt.Sprintf("%s/readyz", s.Callback), 30*time.Second)
		if err != nil {
			return fmt.Errorf("failed to wait for server: %w", err)
		}

		ready()

		return nil
	})

	group.Go(func() error {
		<-ctxShutdown.Done()
		return srv.Shutdown(ctxShutdown)
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
