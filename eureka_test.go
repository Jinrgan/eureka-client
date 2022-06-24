package eureka_client

import (
	"context"
	"testing"
	"time"
)

func TestClient_GetApplications(t *testing.T) {
	ins, err := NewInstance("go-module", 8081)
	if err != nil {
		t.Fatalf("cannot create instance: %v", err)
	}
	clt := Dial(ins, WithURL("http://admin:admin@localhost:8761/eureka"))

	apps, err := clt.GetApplications(context.Background())
	if err != nil {
		t.Errorf("cannot get apps: %+v", err)
	}

	t.Logf("apps: %+v", apps.Applications)
}

func TestClient_Register(t *testing.T) {
	ins, err := NewInstance("go-module", 8081, WithIP("192.168.2.84"), WithVersion("1.0.0"))
	if err != nil {
		t.Fatalf("cannot create instance: %v", err)
	}
	clt := Dial(ins, WithURL("http://admin:admin@localhost:8761/eureka/"))

	err = clt.Register(context.Background())
	if err != nil {
		t.Fatalf("cannot register: %v", err)
	}
}

func TestClient_Heartbeat(t *testing.T) {
	ins, err := NewInstance("go-module", 8081)
	if err != nil {
		t.Fatalf("cannot create instance: %v", err)
	}
	clt := Dial(ins)

	err = clt.Heartbeat(context.Background())
	if err != nil {
		t.Errorf("failed to heartbeat: %v", err)
	}
}

func TestClient_Run(t *testing.T) {
	ins, err := NewInstance("go-module", 8081)
	if err != nil {
		t.Fatalf("cannot create instance: %v", err)
	}
	clt := Dial(ins)

	go func() {
		err := clt.Run(context.Background())
		if err != nil {
			t.Errorf("failed to run client: %v", err)
		}
	}()

	time.Sleep(15 * time.Second)

	clt.Shutdown(context.Background())
}
