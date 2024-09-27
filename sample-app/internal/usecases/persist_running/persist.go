package persistrunning

import (
	"fmt"

	"github.com/yank0vy3rdna/netconf2-meetup-samples/sample-app/internal/repository/sysrepo"
)

type configRepo interface {
	RegisterModuleCallback(callback sysrepo.ModuleChangeCallback) error
	CopyToStartup() error
}

type usecase struct {
	configRepo configRepo
}

func NewUseCase(configRepo configRepo) {
	u := &usecase{configRepo: configRepo}

	u.configRepo.RegisterModuleCallback(u.ConfigChangedCallback)
}

func (u *usecase) ConfigChangedCallback(sysrepo.SessionInterface, string, sysrepo.NotifyEvent) int {
	go func() {
		fmt.Println("copying running to startup ...")
		if err := u.configRepo.CopyToStartup(); err != nil {
			panic(err)
		}
	}()

	return int(sysrepo.ERR_OK)
}
