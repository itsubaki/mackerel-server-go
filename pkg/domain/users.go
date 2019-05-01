package domain

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	ID                       string   `json:"id"`
	ScreenName               string   `json:"screenName"`
	Email                    string   `json:"email"`
	Authority                string   `json:"authority"` // owner, manager, collaborator, viewer
	IsInRegisterationProcess bool     `json:"isInRegisterationProcess"`
	IsMFAEnabled             bool     `json:"isMFAEnabled"`
	AuthenticationMethods    []string `json:"authenticationMethods"`
	JoinedAt                 string   `json:"joinedAt"`
}
