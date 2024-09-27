package httpproxy

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	proxy_usecase "github.com/yank0vy3rdna/netconf2-meetup-samples/sample-app/internal/usecases/proxy"
)

type proxyUseCase interface {
	ProxyTo(net.IP) (uint16, error)
}
type proxy struct {
	server *http.Server

	proxyUseCase proxyUseCase
}

func NewProxy(proxyUseCase proxyUseCase) *proxy {
	p := &proxy{proxyUseCase: proxyUseCase}
	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        p,
	}
	go func() {
		err := s.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		panic(err)
	}()
	p.server = s
	return p
}

func (p *proxy) Close(ctx context.Context) {
	if err := p.server.Shutdown(ctx); err != nil {
		panic(err)
	}
}

type req struct {
	From string `json:"from"`
}
type resp struct {
	To uint16 `json:"to"`
}

func (p *proxy) ServeHTTP(respWriter http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	var request req
	err = json.Unmarshal(body, &request)
	if err != nil {
		respWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	ipFrom := net.ParseIP(request.From)
	if ipFrom == nil {
		respWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	ipFrom = ipFrom.To4()
	if ipFrom == nil {
		respWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	to, err := p.proxyUseCase.ProxyTo(ipFrom)
	if errors.Is(err, proxy_usecase.ErrNoRulesMatched) {
		respWriter.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		respWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(resp{To: to})
	if err != nil {
		respWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	respWriter.Write(bytes)
}
