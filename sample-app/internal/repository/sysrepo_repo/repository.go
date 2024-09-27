package sysrepo_repo

import (
	"encoding/xml"
	"fmt"

	"github.com/yank0vy3rdna/netconf2-meetup-samples/sample-app/internal/domain"
	"github.com/yank0vy3rdna/netconf2-meetup-samples/sample-app/internal/repository/sysrepo"
)

type repo struct {
	conn          *sysrepo.Connection
	session       sysrepo.SessionInterface
	subscriptions []sysrepo.SubscriptionInterface
}

func NewRepo() (*repo, error) {
	conn, err := sysrepo.NewConnection(&sysrepo.Wrapper{})
	if err != nil {
		return nil, err
	}
	session, err := sysrepo.NewSession(conn, sysrepo.Datastore(sysrepo.DS_RUNNING))
	if err != nil {
		return nil, err
	}
	return &repo{
		conn:          conn,
		session:       session,
		subscriptions: make([]sysrepo.SubscriptionInterface, 0),
	}, nil
}

func (r *repo) Close() {
	for _, subscription := range r.subscriptions {
		subscription.Unsubscribe()
	}
	r.session.Stop()
	r.conn.Disconnect()
}

const module = "loadbalancer"

func (r *repo) GetCurrentConfig() (domain.LoadBalancerConfig, error) {
	data, err := r.session.GetDataByModuleName(module)
	if err != nil {
		return domain.LoadBalancerConfig{}, err
	}
	defer data.Free()
	dataXml := data.DataTreeToString(sysrepo.LydFormat(sysrepo.LYD_XML))
	var config domain.LoadBalancerConfig
	err = xml.Unmarshal([]byte(dataXml), &config)
	if err != nil {
		return domain.LoadBalancerConfig{}, err
	}
	return config, nil
}

func (r *repo) RegisterModuleCallback(callback sysrepo.ModuleChangeCallback) error {
	xpath := fmt.Sprintf("/%s:*//.", module)

	subscription, err := r.session.Subscribe(module, xpath, callback)
	if err != nil {
		return fmt.Errorf("failed to subscribe to module changes: %w", err)
	}

	r.subscriptions = append(r.subscriptions, subscription)

	return nil
}

func (r *repo) CopyToStartup() error {
	session, err := sysrepo.NewSession(r.conn, sysrepo.Datastore(sysrepo.DS_STARTUP))
	if err != nil {
		return err
	}
	defer session.Stop()

	return session.CopyConfig(sysrepo.Datastore(sysrepo.DS_RUNNING))
}
