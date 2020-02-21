package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // Register the pprof handlers
	"os"
	"os/signal"
	"syscall"
	"time"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/anthonycorbacho/chatbot/internal/bot"
	"github.com/anthonycorbacho/chatbot/internal/version"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/stats/view"
)

func main() {
	if err := run(); err != nil {
		log.Println("error :", err)
		os.Exit(1)
	}
}

func run() error {
	// Logging
	log := log.New(os.Stdout, "CHATBOT : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// Configuration
	var cfg struct {
		Web struct {
			APIHost         string        `default:"0.0.0.0:3000"`
			DebugHost       string        `default:"0.0.0.0:4000"`
			ReadTimeout     time.Duration `default:"5s"`
			WriteTimeout    time.Duration `default:"5s"`
			ShutdownTimeout time.Duration `default:"5s"`
		}
		Observability struct {
			MetricHost string `default:"0.0.0.0:5000"`
		}
	}
	err := envconfig.Process("chatboot", &cfg)
	if err != nil {
		return fmt.Errorf("paring config: %w", err)
	}

	// Start App
	log.Printf("Started: Application initializing version %v", version.Get())
	defer log.Println("main : Completed")

	// Initialize chatbot
	chatbot := initChatbot()

	// Start Debug Service
	//
	// /debug/pprof - Added to the default mux by importing the net/http/pprof package.
	// /debug/vars - Added to the default mux by importing the expvar package.
	//
	// Not concerned with shutting this down when the application is shutdown.
	log.Println("Started : Initializing debugging support")

	go func() {
		log.Printf("Debug Listening %s", cfg.Web.DebugHost)
		log.Printf("Debug Listener closed : %v", http.ListenAndServe(cfg.Web.DebugHost, http.DefaultServeMux))
	}()

	// Start Observability
	_ = view.Register(ochttp.DefaultServerViews...)

	// register prometheus
	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "chatbot",
	})
	if err != nil {
		return fmt.Errorf("Creating prometheus exporter: %w", err)
	}
	// Ensure that we register it as a stats exporter.
	view.RegisterExporter(pe)

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", pe)
		if err := http.ListenAndServe(cfg.Observability.MetricHost, mux); err != nil {
			log.Printf("Starting prometheus metric server: %e", err)
		}
	}()

	// Start the HTTP server

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/chatbot", chatbotRoute(chatbot))

	srv := http.Server{
		Addr: cfg.Web.APIHost,
		Handler: &ochttp.Handler{
			Handler:     r,
			Propagation: &tracecontext.HTTPFormat{},
		},
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	// Start starts multiplexing the listener.
	go func() {
		serverErrors <- srv.ListenAndServe()
	}()

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Printf("%v : Start shutdown", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
			err = srv.Close()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGSTOP:
			return fmt.Errorf("integrity issue caused shutdown")
		case err != nil:
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}
	return nil
}

type response struct {
	Answer   string `json:"answer"`
	Question string `json:"question"`
}

func chatbotRoute(chatbot *bot.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		question := r.URL.Query().Get("sentence")
		if question == "" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		rep := response{
			Answer:   chatbot.Sentence(r.Context(), question),
			Question: question,
		}

		js, err := json.Marshal(rep)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js) // nolint
	}
}

func initChatbot() *bot.Bot {
	chatbot := bot.New()

	// Better implementation would be to construct this set of answer from a file.
	chatbot.Respond("ping", func(msg bot.Message) string { // nolint
		return "pong!"
	})
	chatbot.Respond("what is your name?", func(msg bot.Message) string { // nolint
		return "My name is Anthony."
	})
	chatbot.Respond("Where are you located?", func(msg bot.Message) string { // nolint
		return "Seoul"
	})

	chatbot.Respond("What is your age?", func(msg bot.Message) string { // nolint
		return "That is a secret :)"
	})

	return chatbot
}
