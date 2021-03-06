package frontend

import (
	"fmt"
	"net"
	"strconv"

	"github.com/convox/praxis/helpers"
)

type Proxy struct {
	Host   string
	Port   string
	Target string

	frontend *Frontend
	listener net.Listener
}

func NewProxy(host, port, target string, frontend *Frontend) *Proxy {
	return &Proxy{
		Host:     host,
		Port:     port,
		Target:   target,
		frontend: frontend,
	}
}

func (p *Proxy) Close() error {
	if err := p.listener.Close(); err != nil {
		return err
	}

	return nil
}

func (p *Proxy) Serve() error {
	log := Log.At("proxy.create").Namespace("host=%s port=%s target=%s", p.Host, p.Port, p.Target).Start()

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", p.Host, p.Port))
	if err != nil {
		return err
	}

	p.listener = ln

	log.Success()

	if err := p.handleListener(); err != nil {
		return err
	}

	return nil
}

func (p *Proxy) handleListener() error {
	for {
		cn, err := p.listener.Accept()
		if err != nil {
			return err
		}

		go p.handleConnection(cn)
	}
}

func (p *Proxy) handleConnection(cn net.Conn) {
	log := Log.At("proxy.connect").Start()

	defer cn.Close()

	ip, port, err := net.SplitHostPort(cn.LocalAddr().String())
	if err != nil {
		log.Error(err)
		return
	}

	log = log.Namespace("ip=%s port=%s", ip, port)

	pi, err := strconv.Atoi(port)
	if err != nil {
		log.Error(err)
		return
	}

	ep, ok := p.frontend.endpoints[fmt.Sprintf("%s:%d", ip, pi)]
	if !ok {
		log.Error(fmt.Errorf("no endpoint"))
		return
	}

	log = log.Namespace("target=%s host=%q", ep.Target, ep.Host)

	out, err := net.Dial("tcp", ep.Target)
	if err != nil {
		log.Error(err)
		return
	}

	defer out.Close()

	if err := helpers.Pipe(cn, out); err != nil {
		log.Error(err)
		return
	}

	log.Success()
}
