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
	LinkId   int    `json:"link_id"`
	Day      string `json:"day"`
	IsUnique bool   `json:"is_unique"`
}

func (p *Plugin) Init(cfg config.Configurer, log *zap.Logger) error {
	p.clicks = make(chan string)
	p.log = log

	const op = errors.Op("roadrunner_http_test_init")
	if !cfg.Has(name) {
		return errors.E(errors.Disabled)
	}

	p.cfg = &Config{}
	err := cfg.UnmarshalKey(name, p.cfg)

	if err != nil {
		return errors.E(op, err)
	}

	p.cfg.InitDefaults()

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
				go func() {

					p.log.Info(c)
					click := Click{}
					err := json.Unmarshal([]byte(c), &click)

					if err != nil {
						panic(err.Error()) // proper error handling instead of panic in your app
					}

					if click.IsUnique {
						_, err = p.db.Exec("INSERT INTO daily (day, link_id, clicks, unique_clicks) VALUES (?, ?, 1, 1) ON DUPLICATE KEY UPDATE clicks=clicks+1, unique_clicks=unique_clicks+1;", click.Day, click.LinkId)
					} else {
						_, err = p.db.Exec("INSERT INTO daily (day, link_id, clicks) VALUES (?, ?, 1) ON DUPLICATE KEY UPDATE clicks=clicks+1;", click.Day, click.LinkId)
					}

					if err != nil {
						panic(err.Error()) // proper error handling instead of panic in your app
					}
				}()
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
