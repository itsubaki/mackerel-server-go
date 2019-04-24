package org

import "encoding/json"

type GetOrgOutput struct {
	Name string `json:"name"`
}

func (o *GetOrgOutput) String() string {
	bytea, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}
