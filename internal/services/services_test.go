package services

import "testing"

func TestGetServicesOutput(t *testing.T) {
	out := &GetServicesOutput{
		Services: []Service{
			{
				Name: "test-service-01",
				Memo: "memo",
				Roles: []string{
					"app",
					"db",
					"cache",
				},
			},
		},
	}

	if out.String() != "{\"services\":[{\"name\":\"test-service-01\",\"memo\":\"memo\",\"roles\":[\"app\",\"db\",\"cache\"]}]}" {
		t.Error(out)
	}
}
