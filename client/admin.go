package client

// HealthStatus is the response from the health endpoint.
type HealthStatus struct {
	Status string `json:"status"`
}

// StandInStatus is the response from the stand-in admin endpoints.
type StandInStatus struct {
	Enabled     bool `json:"enabled"`
	QueueLength int  `json:"queue_length"`
}

// Health checks the server health.
func (c *Client) Health() (*HealthStatus, error) {
	var s HealthStatus
	err := c.doPlainJSON("GET", "/health", nil, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// GetStandIn retrieves the current stand-in mode status.
func (c *Client) GetStandIn() (*StandInStatus, error) {
	var s StandInStatus
	err := c.doPlainJSON("GET", "/admin/standin", nil, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// SetStandIn enables or disables stand-in mode.
func (c *Client) SetStandIn(enabled bool) (*StandInStatus, error) {
	req := struct {
		Enabled bool `json:"enabled"`
	}{Enabled: enabled}
	var s StandInStatus
	err := c.doPlainJSON("PUT", "/admin/standin", req, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
