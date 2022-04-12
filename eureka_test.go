package eureka_client

import (
	"context"
	"testing"
)

func TestClient_GetApplications(t *testing.T) {
	clt := Dial(WithURL("http://admin:admin@localhost:8761/eureka"))

	apps, err := clt.GetApplications(context.Background())
	if err != nil {
		t.Errorf("cannot get apps: %+v", err)
	}

	t.Logf("apps: %+v", apps.Applications)
}

func TestClient_Register(t *testing.T) {
	clt := Dial(WithURL("http://admin:admin@localhost:8761/eureka"))

	err := clt.Register(context.Background(), NewInstance("go-module", "192.168.31.236", 8081))
	if err != nil {
		t.Fatalf("cannot register: %v", err)
	}
}
