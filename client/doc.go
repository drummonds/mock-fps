// Package client provides a Go client for the mock-fps server.
//
// It wraps all JSON:API endpoints (payments, submissions, admissions,
// returns, recalls, reversals, subscriptions) and admin endpoints
// (health, stand-in mode) with typed methods that handle envelope
// wrapping/unwrapping and error decoding.
//
// Usage:
//
//	c := client.New("http://localhost:8080")
//	p, err := c.CreatePayment(client.Payment{...})
package client
