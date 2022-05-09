package register

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bpalermo/memento/pkg/endpoint"
	"github.com/bpalermo/memento/pkg/util"
	"github.com/sirupsen/logrus"
	clientV3 "go.etcd.io/etcd/client/v3"
	"path"
	"strings"
	"time"
)

type EtcdRegister struct {
	cli            *clientV3.Client
	logger         *logrus.Logger
	leaseID        clientV3.LeaseID
	keepAliveChan  <-chan *clientV3.LeaseKeepAliveResponse
	serviceName    string
	basePath       string
	endpoint       *endpoint.Endpoint
	ttl            time.Duration
	defaultTimeout time.Duration
}

// NewEtcdRegister EtcdRegister factory method
func NewEtcdRegister(cfg clientV3.Config, logger *logrus.Logger, basePath string, serviceName string, servicePort uint16, ttl time.Duration, defaultTimeout time.Duration) (register *EtcdRegister, err error) {
	localIp, err := util.LocalIP()
	if err != nil {
		logger.WithError(err).Error("could not determine local IP")
		return nil, err
	}

	serviceEndpoint := endpoint.NewEndpoint(localIp, servicePort)

	logger.Debug("Initializing Etcd client")
	cli, err := clientV3.New(cfg)
	if err != nil {
		logger.WithError(err).Error("could not initialize Etcd client")
		return nil, err
	}
	logger.Debugf("Registry initialized with basePath=%s, endpoint=%v, defaultTimeout=%s", basePath, serviceEndpoint, defaultTimeout)

	return &EtcdRegister{
		cli:            cli,
		logger:         logger,
		basePath:       basePath,
		serviceName:    serviceName,
		endpoint:       serviceEndpoint,
		ttl:            ttl,
		defaultTimeout: defaultTimeout,
	}, nil
}

func (r *EtcdRegister) Register() (err error) {
	if r.basePath == "" {
		return fmt.Errorf("basePath must be non empty")
	}
	if r.endpoint == nil {
		return fmt.Errorf("endpoint must be non nil")
	}

	if err = r.putKeyWithLease(); err != nil {
		return err
	}

	return nil
}

func (r *EtcdRegister) putKeyWithLease() (err error) {
	r.logger.Debugf("Creating lease grant for %s/%s:%d", r.serviceName, r.endpoint.Ip, r.endpoint.Port)
	ctx, cancel := context.WithTimeout(context.Background(), r.defaultTimeout)
	defer cancel()
	grantResponse, err := r.cli.Grant(ctx, int64(r.ttl.Seconds()))
	if err != nil {
		return err
	}

	r.logger.Debugf("Creating endpoint for %s/%s:%d", r.serviceName, r.endpoint.Ip, r.endpoint.Port)
	ctx, cancel = context.WithTimeout(context.Background(), r.defaultTimeout)
	defer cancel()
	_, err = r.cli.Put(ctx, r.endpointPath(r.serviceName, r.endpoint), encode(r.endpoint), clientV3.WithLease(grantResponse.ID))
	if err != nil {
		return err
	}

	r.logger.Debugf("Starting automatic keep alive for %s/%s:%d", r.serviceName, r.endpoint.Ip, r.endpoint.Port)
	keepAliveChan, err := r.cli.KeepAlive(ctx, grantResponse.ID)
	if err != nil {
		return err
	}

	r.leaseID = grantResponse.ID
	r.keepAliveChan = keepAliveChan

	r.logger.Infof("Registration succeeded for %s/%s:%d", r.serviceName, r.endpoint.Ip, r.endpoint.Port)

	return nil
}

// Listen listen and watch
func (r *EtcdRegister) Listen() {
	for leaseKeepResp := range r.keepAliveChan {
		r.logger.Debugf("[%s %s:%d] %s", r.serviceName, r.endpoint.Ip, r.endpoint.Port, leaseKeepResp)
	}
}

// Close revoke lease and close client
func (r *EtcdRegister) Close() (err error) {
	// Revoke lease
	if _, err = r.cli.Revoke(context.Background(), r.leaseID); err != nil {
		return err
	}
	r.logger.Infof("Lease %d revoked for %s/%s:%d", r.leaseID, r.serviceName, r.endpoint.Ip, r.endpoint.Port)
	return r.cli.Close()
}

func (r *EtcdRegister) endpointPath(serviceName string, endpoint *endpoint.Endpoint) string {
	service := strings.Replace(serviceName, "/", "-", -1)
	return path.Join(r.basePath, service, fmt.Sprintf("%s:%d", endpoint.Ip, endpoint.Port))
}

func encode(t interface{}) string {
	if t != nil {
		b, _ := marshal(t)
		return string(b)
	}
	return ""
}

func marshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
