package domain

type Invitations struct {
	Invitations []Invitation `json:"invitations"`
}

type Invitation struct {
	EMail     string `json:"email"`
	Authority string `json:"authority"` // manager,collaborator,viewer
	ExpiresAt int    `json:"expiresAt,omitempty"`
}

type Revoke struct {
	EMail string `json:"email"`
}
