package roadrunner_http_test

import "go.uber.org/zap"

// plugin name
const name = "roadrunner_http_test"

// Plugin structure should have exactly the `Plugin` name to be found by RR
type Plugin struct {
	clicks *chan string
	log    *zap.Logger
}

func (p *Plugin) Init(log *zap.Logger) error {
	*p.clicks = make(chan string)
	p.log = log
	return nil
}

func (p *Plugin) Serve() chan error {

	errCh := make(chan error, 1)

	go func() {
		for {
			select {
			case click := <-*p.clicks:
				p.log.Info(click)
			default:
			}
		}
	}()
	return errCh
}

func (p *Plugin) Stop() error {
	return nil
}

// Name this is not mandatory, but if you implement this interface and provide a plugin name, RR will expose the RPC method of this plugin using this name
func (p *Plugin) Name() string {
	return name
}

// ----------------------------------------------------------------------------
// RPC
// ----------------------------------------------------------------------------

type rpc struct {
	srv *Plugin
}

// RPC interface implementation, RR will find this interface and automatically expose the RPC endpoint with methods (rpc structure)
func (p *Plugin) RPC() interface{} {
	return &rpc{}
}

// AddClick Generate this is the function exposed to PHP $rpc->call(), can be any name
func (r *rpc) AddClick(input string, output *string) error {
	*r.srv.clicks <- input
	*output = input
	return nil
}
