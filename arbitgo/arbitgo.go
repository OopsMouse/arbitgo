package main

import (
	"os"

	"github.com/OopsMouse/arbitgo/models"

	"github.com/OopsMouse/arbitgo/infrastructure"
	"github.com/OopsMouse/arbitgo/usecase"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "arbitgo"
	app.Usage = "A Bot for arbit rage with one exchange, multi currency"
	app.Version = "0.0.1"

	var debug bool
	var dryrun bool
	var apiKey string
	var secret string
	var server string

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "debug mode",
			Destination: &debug,
		},
		cli.BoolFlag{
			Name:        "dryrun, dry, d",
			Usage:       "dry run mode",
			Destination: &dryrun,
		},
		cli.StringFlag{
			Name:        "apikey, a",
			Usage:       "api key of exchange",
			Destination: &apiKey,
			EnvVar:      "EXCHANGE_APIKEY",
		},
		cli.StringFlag{
			Name:        "secret, s",
			Usage:       "secret of exchange",
			Destination: &secret,
			EnvVar:      "EXCHANGE_SECRET",
		},
		cli.StringFlag{
			Name:        "server",
			Usage:       "server host",
			Destination: &server,
		},
	}

	app.Action = func(c *cli.Context) error {
		if apiKey == "" || secret == "" {
			return cli.NewExitError("api key and secret is required", 0)
		}
		return nil
	}

	app.Run(os.Args)

	logInit(debug)
	exchange := newExchange(apiKey, secret, dryrun)
	arbitrader := newTrader(exchange, &server)
	arbitrader.Run()
}

func newExchange(apikey string, secret string, dryRun bool) usecase.Exchange {
	binance := infrastructure.NewBinance(
		apikey,
		secret,
	)

	if dryRun {
		balances := map[string]*models.Balance{}
		balances["BTC"] = &models.Balance{
			Asset: "BTC",
			Free:  0.01,
			Total: 0.01,
		}
		return infrastructure.NewExchangeStub(
			binance,
			balances,
		)
	}

	return binance
}

func newTrader(exchange usecase.Exchange, server *string) *usecase.Trader {
	return usecase.NewTrader(
		exchange,
		server,
	)
}

func logInit(debug bool) {
	format := &log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	log.SetFormatter(format)
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	log.SetOutput(os.Stdout)
}
