package client

import "fmt"

const subscriptionsPath = "/v1/notification/subscriptions"

// CreateSubscription creates a webhook subscription.
func (c *Client) CreateSubscription(sub Subscription) (*Subscription, error) {
	var env DataEnvelope[Subscription]
	err := c.doJSON("POST", subscriptionsPath,
		DataEnvelope[Subscription]{Data: sub}, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// GetSubscription retrieves a subscription by ID.
func (c *Client) GetSubscription(id string) (*Subscription, error) {
	var env DataEnvelope[Subscription]
	err := c.doJSON("GET",
		fmt.Sprintf("%s/%s", subscriptionsPath, id),
		nil, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// ListSubscriptions lists all subscriptions.
func (c *Client) ListSubscriptions() ([]Subscription, error) {
	var env ListEnvelope[Subscription]
	err := c.doJSON("GET", subscriptionsPath, nil, &env)
	if err != nil {
		return nil, err
	}
	return env.Data, nil
}

// UpdateSubscription patches a subscription.
func (c *Client) UpdateSubscription(id string, sub Subscription) (*Subscription, error) {
	var env DataEnvelope[Subscription]
	err := c.doJSON("PATCH",
		fmt.Sprintf("%s/%s", subscriptionsPath, id),
		DataEnvelope[Subscription]{Data: sub}, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// DeleteSubscription deletes a subscription.
func (c *Client) DeleteSubscription(id string) error {
	_, err := c.doJSONWithStatus("DELETE",
		fmt.Sprintf("%s/%s", subscriptionsPath, id),
		nil, nil)
	return err
}
