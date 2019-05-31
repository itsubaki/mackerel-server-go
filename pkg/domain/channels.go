package domain

type Channels struct {
	Channels []Channel `json:"channels"`
}

type Channel struct {
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
