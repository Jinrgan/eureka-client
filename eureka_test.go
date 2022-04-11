package eureka_client

import (
	"context"
	"testing"
)

func TestService_GetApplications(t *testing.T) {
	svc := Dial(WithURL("http://admin:admin@localhost:8761/eureka"))

	apps, err := svc.GetApplications(context.Background())
	if err != nil {
		t.Errorf("cannot get apps: %+v", err)
	}

	t.Logf("apps: %+v", apps.Applications)
}
