package register

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	clientV3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

type etcdContainer struct {
	testcontainers.Container
	Endpoint string
}

func setupEtcd(ctx context.Context) (*etcdContainer, error) {
	req := testcontainers.ContainerRequest{
		Image: "quay.io/coreos/etcd:v3.5.4",
		Entrypoint: []string{
			"/usr/local/bin/etcd",
		},
		Cmd: []string{
			"--name",
			"node1",
			"--initial-advertise-peer-urls",
			"http://127.0.0.1:2380",
			"--listen-peer-urls",
			"http://0.0.0.0:2380",
			"--advertise-client-urls",
			"http://127.0.0.1:2379",
			"--listen-client-urls",
			"http://0.0.0.0:2379",
			"--initial-cluster",
			"node1=http://127.0.0.1:2380",
		},
		ExposedPorts: []string{
			"2379:2379/tcp",
			"2380:2380/tcp",
		},
		WaitingFor: wait.ForLog("serving client traffic insecurely; this is strongly discouraged!"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "2379")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	return &etcdContainer{Container: container, Endpoint: fmt.Sprintf("%s:%s", hostIP, mappedPort.Port())}, nil
}

func TestEtcdRegister_Register(t *testing.T) {
	ctx := context.Background()

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	container, err := setupEtcd(ctx)
	if err != nil {
		t.Fatal(err)
	}

	r, err := NewEtcdRegister(clientV3.Config{
		Endpoints: []string{
			container.Endpoint,
		},
		DialTimeout: 15 * time.Second,
	}, logger, "/discovery", "fake", 18080, 15*time.Second, 45*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, r)

	err = r.Register()
	if err != nil {
		t.Fatal(err)
	}

	defer terminate(t, container, ctx)
}

func terminate(t *testing.T, container *etcdContainer, ctx context.Context) {
	err := container.Terminate(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
