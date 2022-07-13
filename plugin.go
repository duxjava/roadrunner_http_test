package roadrunner_http_test

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/roadrunner-server/api/v2/plugins/config"
	"github.com/roadrunner-server/errors"
	"go.uber.org/zap"
	"time"
)

// plugin name
const name = "roadrunner_http_test"

// Plugin structure should have exactly the `Plugin` name to be found by RR
type Plugin struct {
	clicks chan string
	log    *zap.Logger
	cfg    *Config
	db     *sql.DB
}

type Click struct {
	Id       int    `json:"id"`
	Day      string `json:"day"`
	IsUnique bool   `json:"isUnique"`
}

type Link struct {
	Id     int    `json:"id"`
	UserId string `json:"user_id"`
	Url    bool   `json:"url"`
	GenUrl bool   `json:"gen_url"`
}

func (p *Plugin) Init(cfg config.Configurer, log *zap.Logger) error {
	p.clicks = make(chan string)
	p.log = log

	const op = errors.Op("my_plugin_init")
	if !cfg.Has(name) {
		return errors.E(errors.Disabled)
	}

	p.cfg = &Config{}
	err := cfg.UnmarshalKey(name, p.cfg)

	if err != nil {
		return errors.E(op, err)
	}

	db, err := sql.Open("mysql", p.cfg.Mysql.Connection)
	db.SetMaxIdleConns(p.cfg.Mysql.Maxidle)
	db.SetMaxOpenConns(p.cfg.Mysql.Maxopen)
	db.SetConnMaxLifetime(time.Second * p.cfg.Mysql.Lifetime)

	if err != nil {
		return errors.E(op, err)
	}

	err = db.Ping()

	if err != nil {
		return errors.E(op, err)
	}

	p.db = db

	return nil
}

func (p *Plugin) Serve() chan error {

	errCh := make(chan error, 1)

	go func() {
		for {
			select {
			case c := <-p.clicks:
				p.log.Info(c)
				click := Click{}
				_ = json.Unmarshal([]byte(c), &click)

				link := Link{}

				err := p.db.QueryRow("SELECT id, user_id, url, gen_url FROM links WHERE id = ?", click.Id).Scan(&link.Id, &link.UserId, &link.Url, &link.GenUrl)

				if err != nil {
					panic(err.Error()) // proper error handling instead of panic in your app
				}

				s, _ := json.Marshal(link)

				p.log.Info(string(s))

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
	clicks chan string
}

// RPC interface implementation, RR will find this interface and automatically expose the RPC endpoint with methods (rpc structure)
func (p *Plugin) RPC() interface{} {
	rpc := &rpc{}
	rpc.clicks = p.clicks
	return rpc
}

// AddClick Generate this is the function exposed to PHP $rpc->call(), can be any name
func (r *rpc) AddClick(input string, output *string) error {
	r.clicks <- input
	*output = input
	return nil
}
