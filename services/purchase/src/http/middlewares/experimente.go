package middlewares

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequestLogger(c *fiber.Ctx) error {
	// URL untuk permintaan yang akan dilacak
	url := "https://jsonplaceholder.typicode.com/posts"

	// Variabel untuk mencatat waktu setiap tahap
	var (
		dnsStart, dnsEnd         time.Time
		connectStart, connectEnd time.Time
		tlsStart, tlsEnd         time.Time
		gotFirstResponseByte     time.Time
		totalRequestStart        time.Time
	)

	// Buat permintaan HTTP dengan trace
	req, _ := http.NewRequest("GET", url, nil)

	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			dnsStart = time.Now()
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			dnsEnd = time.Now()
		},
		ConnectStart: func(network, addr string) {
			connectStart = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			connectEnd = time.Now()
		},
		TLSHandshakeStart: func() {
			tlsStart = time.Now()
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			tlsEnd = time.Now()
		},
		GotFirstResponseByte: func() {
			gotFirstResponseByte = time.Now()
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	// Catat waktu total permintaan
	totalRequestStart = time.Now()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Request failed: %v", err))
	}
	defer resp.Body.Close()

	// Log ke terminal
	fmt.Println("Request Details:")
	fmt.Printf("  DNS Lookup:           %v\n", dnsEnd.Sub(dnsStart))
	fmt.Printf("  TCP Handshake:        %v\n", connectEnd.Sub(connectStart))
	fmt.Printf("  TLS Handshake:        %v\n", tlsEnd.Sub(tlsStart))
	fmt.Printf("  Time to First Byte:   %v\n", gotFirstResponseByte.Sub(totalRequestStart))
	fmt.Printf("  Total Request Time:   %v\n", time.Since(totalRequestStart))

	// Lanjutkan ke handler berikutnya
	return c.Next()
}
