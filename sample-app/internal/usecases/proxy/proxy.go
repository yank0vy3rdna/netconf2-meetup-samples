package proxy

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/yank0vy3rdna/netconf2-meetup-samples/sample-app/internal/domain"
	"github.com/yank0vy3rdna/netconf2-meetup-samples/sample-app/internal/repository/sysrepo"
)

type configRepo interface {
	GetCurrentConfig() (domain.LoadBalancerConfig, error)
	RegisterModuleCallback(callback sysrepo.ModuleChangeCallback) error
}

type usecase struct {
	configRepo configRepo

	currentConfig domain.LoadBalancerConfig
	mu            sync.RWMutex
}

func NewUseCase(configRepo configRepo) (*usecase, error) {
	u := &usecase{configRepo: configRepo}

	if err := u.updateConfig(); err != nil {
		return nil, err
	}
	u.configRepo.RegisterModuleCallback(u.ConfigChangedCallback)

	return u, nil
}

func (u *usecase) updateConfig() error {
	config, err := u.configRepo.GetCurrentConfig()
	if err != nil {
		return err
	}

	u.mu.Lock()
	u.currentConfig = config
	u.mu.Unlock()

	bytes, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("config updated: %s\n", string(bytes))

	return nil
}

func (u *usecase) ConfigChangedCallback(sysrepo.SessionInterface, string, sysrepo.NotifyEvent) int {
	go func() {
		if err := u.updateConfig(); err != nil {
			panic(err)
		}
	}()

	return int(sysrepo.ERR_OK)
}
