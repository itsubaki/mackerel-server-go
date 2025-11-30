package domain

type Channels struct {
	Channels []any `json:"channels"`
}

type Channel struct {
	OrgID             string            `json:"-"`
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Type              string            `json:"type"`
	Events            []string          `json:"events,omitempty"`            // email, slack. webhook
	Emails            []string          `json:"emails,omitempty"`            // email
	UserIDs           []string          `json:"userIds,omitempty"`           // email
	URL               string            `json:"url,omitempty"`               // slack, webhook
	Mentions          map[string]string `json:"mentions,omitempty"`          // slack
	EnabledGraphImage bool              `json:"enabledGraphImage,omitempty"` // slack
}

func (c *Channel) Cast() any {
	if c.Type == "email" {
		return &EmailChannel{
			ID:      c.ID,
			Name:    c.Name,
			Type:    c.Type,
			Events:  c.Events,
			UserIDs: c.UserIDs,
			Emails:  c.Emails,
		}
	}

	if c.Type == "slack" {
		return &SlackChannel{
			ID:                c.ID,
			Name:              c.Name,
			Type:              c.Type,
			Events:            c.Events,
			URL:               c.URL,
			Mentions:          c.Mentions,
			EnabledGraphImage: c.EnabledGraphImage,
		}
	}

	if c.Type == "webhook" {
		return &WebhookChannel{
			ID:     c.ID,
			Name:   c.Name,
			Type:   c.Type,
			Events: c.Events,
			URL:    c.URL,
		}
	}

	return c
}

type EmailChannel struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Events  []string `json:"events"`
	Emails  []string `json:"emails"`
	UserIDs []string `json:"userIds"`
}

type SlackChannel struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Type              string            `json:"type"`
	Events            []string          `json:"events"`
	URL               string            `json:"url"`
	Mentions          map[string]string `json:"mentions"`
	EnabledGraphImage bool              `json:"enabledGraphImage"`
}

type WebhookChannel struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Events []string `json:"events"`
	URL    string   `json:"url"`
}
