package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nibble/mock-fps/internal/handlers"
	"github.com/nibble/mock-fps/internal/jsonapi"
	"github.com/nibble/mock-fps/internal/lifecycle"
	"github.com/nibble/mock-fps/internal/models"
	"github.com/nibble/mock-fps/internal/store"
)

func setupBenchServer() *httptest.Server {
	s := store.NewMemoryStore()
	// Use large delay so lifecycle goroutines don't interfere with benchmarks
	engine := lifecycle.NewEngine(999999, nil)
	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, s, engine)
	return httptest.NewServer(mux)
}

func BenchmarkCreatePayment(b *testing.B) {
	srv := setupBenchServer()
	defer srv.Close()
	client := srv.Client()

	b.ResetTimer()
	for i := range b.N {
		p := models.Payment{
			Resource:   models.Resource{ID: fmt.Sprintf("p-%d", i)},
			Attributes: models.PaymentAttributes{Amount: "100.00", Currency: "GBP"},
		}
		body, _ := json.Marshal(jsonapi.DataEnvelope[models.Payment]{Data: p})
		resp, err := client.Post(srv.URL+"/v1/transaction/payments", jsonapi.ContentType, bytes.NewReader(body))
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func BenchmarkGetPayment(b *testing.B) {
	srv := setupBenchServer()
	defer srv.Close()
	client := srv.Client()

	// Seed one payment
	p := models.Payment{
		Resource:   models.Resource{ID: "bench-get"},
		Attributes: models.PaymentAttributes{Amount: "100.00", Currency: "GBP"},
	}
	body, _ := json.Marshal(jsonapi.DataEnvelope[models.Payment]{Data: p})
	resp, _ := client.Post(srv.URL+"/v1/transaction/payments", jsonapi.ContentType, bytes.NewReader(body))
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	b.ResetTimer()
	for range b.N {
		resp, err := client.Get(srv.URL + "/v1/transaction/payments/bench-get")
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func BenchmarkCreatePaymentParallel(b *testing.B) {
	srv := setupBenchServer()
	defer srv.Close()
	client := srv.Client()

	var counter atomic.Int64
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := counter.Add(1)
			p := models.Payment{
				Resource:   models.Resource{ID: fmt.Sprintf("pp-%d", id)},
				Attributes: models.PaymentAttributes{Amount: "100.00", Currency: "GBP"},
			}
			body, _ := json.Marshal(jsonapi.DataEnvelope[models.Payment]{Data: p})
			resp, err := client.Post(srv.URL+"/v1/transaction/payments", jsonapi.ContentType, bytes.NewReader(body))
			if err != nil {
				b.Fatal(err)
			}
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	})
}

func BenchmarkGetPaymentParallel(b *testing.B) {
	srv := setupBenchServer()
	defer srv.Close()
	client := srv.Client()

	// Seed
	p := models.Payment{
		Resource:   models.Resource{ID: "bench-pget"},
		Attributes: models.PaymentAttributes{Amount: "100.00", Currency: "GBP"},
	}
	body, _ := json.Marshal(jsonapi.DataEnvelope[models.Payment]{Data: p})
	resp, _ := client.Post(srv.URL+"/v1/transaction/payments", jsonapi.ContentType, bytes.NewReader(body))
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Get(srv.URL + "/v1/transaction/payments/bench-pget")
			if err != nil {
				b.Fatal(err)
			}
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	})
}

// LoadTest runs a sustained load test and prints throughput + latency stats.
// Run with: go test -run TestLoadTest -v -count=1 ./internal/handlers/
func TestLoadTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping load test in short mode")
	}

	srv := setupBenchServer()
	defer srv.Close()
	client := srv.Client()

	const (
		duration    = 5 * time.Second
		concurrency = 50
	)

	var (
		ops      atomic.Int64
		errors   atomic.Int64
		totalNs  atomic.Int64
		maxNs    atomic.Int64
	)

	updateMax := func(val int64) {
		for {
			cur := maxNs.Load()
			if val <= cur || maxNs.CompareAndSwap(cur, val) {
				return
			}
		}
	}

	deadline := time.Now().Add(duration)
	var wg sync.WaitGroup

	// Mix of writes and reads
	for i := range concurrency {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			localOps := int64(0)
			for time.Now().Before(deadline) {
				start := time.Now()

				// Alternate: even workers create, odd workers read
				if workerID%2 == 0 {
					id := fmt.Sprintf("load-%d-%d", workerID, localOps)
					p := models.Payment{
						Resource:   models.Resource{ID: id},
						Attributes: models.PaymentAttributes{Amount: "100.00", Currency: "GBP"},
					}
					body, _ := json.Marshal(jsonapi.DataEnvelope[models.Payment]{Data: p})
					resp, err := client.Post(srv.URL+"/v1/transaction/payments", jsonapi.ContentType, bytes.NewReader(body))
					if err != nil {
						errors.Add(1)
						continue
					}
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
					if resp.StatusCode != http.StatusCreated {
						errors.Add(1)
					}
				} else {
					resp, err := client.Get(srv.URL + "/v1/transaction/payments")
					if err != nil {
						errors.Add(1)
						continue
					}
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
				}

				elapsed := time.Since(start).Nanoseconds()
				totalNs.Add(elapsed)
				updateMax(elapsed)
				ops.Add(1)
				localOps++
			}
		}(i)
	}

	wg.Wait()

	totalOps := ops.Load()
	totalErrs := errors.Load()
	avgNs := int64(0)
	if totalOps > 0 {
		avgNs = totalNs.Load() / totalOps
	}

	t.Logf("=== Load Test Results ===")
	t.Logf("Duration:    %s", duration)
	t.Logf("Concurrency: %d", concurrency)
	t.Logf("Total ops:   %d", totalOps)
	t.Logf("Throughput:  %.0f ops/sec", float64(totalOps)/duration.Seconds())
	t.Logf("Avg latency: %s", time.Duration(avgNs))
	t.Logf("Max latency: %s", time.Duration(maxNs.Load()))
	t.Logf("Errors:      %d", totalErrs)
}
