package domain

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	ID                      string   `json:"id"`
	ScreenName              string   `json:"screenName"`
	Email                   string   `json:"email"`
	Authority               string   `json:"authority"` // owner, manager, collaborator, viewer
	IsInRegistrationProcess bool     `json:"isInRegistrationProcess"`
	IsMFAEnabled            bool     `json:"isMFAEnabled"`
	AuthenticationMethods   []string `json:"authenticationMethods"`
	JoinedAt                int64    `json:"joinedAt"`
}
