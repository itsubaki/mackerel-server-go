package domain

type Invitations struct {
	Invitations []Invitation `json:"invitations"`
}

type Invitation struct {
	OrgID     string `json:"-"`
	EMail     string `json:"email"`
	Authority string `json:"authority"` // manager,collaborator,viewer
	ExpiresAt int64  `json:"expiresAt,omitempty"`
}

type Revoke struct {
	EMail string `json:"email"`
}
