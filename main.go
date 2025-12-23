package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type verboseResponse struct {
	Timestamp string              `json:"timestamp"`
	Method    string              `json:"method"`
	URL       string              `json:"url"`
	Query     map[string][]string `json:"query"`
	Headers   map[string][]string `json:"headers"`
	Body      string              `json:"body"`
}

var version = "dev"

func getEnvBool(key string, defaultVal bool) bool {
	val := strings.ToLower(os.Getenv(key))
	if val == "true" {
		return true
	} else if val == "false" {
		return false
	}
	return defaultVal
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	logMethodURL := getEnvBool("PONG_LOG_METHOD_URL", true)
	logHeaders := getEnvBool("PONG_LOG_HEADERS", false)
	logBody := getEnvBool("PONG_LOG_BODY", false)
	verboseResponseMode := getEnvBool("PONG_VERBOSE_RESPONSE", false)

	var logBuilder strings.Builder

	if logMethodURL {
		fmt.Fprintf(&logBuilder, "\nMethod: %s\nURL   : %s\n", r.Method, r.URL.String())
	}

	if logHeaders {
		for name, values := range r.Header {
			for _, value := range values {
				fmt.Fprintf(&logBuilder, "Header: %s = %s\n", name, value)
			}
		}
	}

	var bodyStr string
	if logBody || verboseResponseMode {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(&logBuilder, "Error reading body: %v\n", err)
		} else {
			bodyStr = string(body)
			if logBody {
				fmt.Fprintf(&logBuilder, "Body: %s\n", bodyStr)
			}
		}
		r.Body.Close()
	} else {
		defer r.Body.Close()
	}

	log.Print(logBuilder.String())

	if verboseResponseMode {
		resp := verboseResponse{
			Timestamp: time.Now().Format("2006-01-02T15:04:05.000Z07:00"),
			Method:    r.Method,
			URL:       r.URL.Path,
			Query:     r.URL.Query(),
			Headers:   r.Header,
			Body:      bodyStr,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	} else {
		// Default simple OK response
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Ok")
	}
}

func resolveVersion() string {
	if strings.TrimSpace(version) == "" {
		return "dev"
	}
	return version
}

func main() {
	shortVersionFlag := flag.Bool("v", false, "print version")
	longVersionFlag := flag.Bool("version", false, "print version")
	flag.Parse()

	if *shortVersionFlag || *longVersionFlag {
		fmt.Println(resolveVersion())
		return
	}

	// Set custom log format for timestamps
	log.SetFlags(0)
	log.SetOutput(logWriter{})

	addr := os.Getenv("PONG_ECHO_SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	http.HandleFunc("/", echoHandler)
	log.Printf("Starting echo server on %s (version %s)", addr, resolveVersion())
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

type logWriter struct{}

func (logWriter) Write(p []byte) (n int, err error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	return fmt.Printf("%s %s", timestamp, p)
}
