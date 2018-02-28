package main

import (
	"os"

	"github.com/OopsMouse/arbitgo/models"

	"github.com/OopsMouse/arbitgo/infrastructure"
	"github.com/OopsMouse/arbitgo/usecase"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "arbitgo"
	app.Usage = "A Bot for arbit rage with one exchange, multi currency"
	app.Version = "0.0.1"

	var dryrun bool
	var apiKey string
	var secret string
	var assetString string
	var maxqty float64
	var threshold float64

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "dryrun, d",
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
			Name:        "asset, as",
			Usage:       "main asset",
			Destination: &assetString,
			Value:       "BTC",
		},
		cli.Float64Flag{
			Name:        "maxqty, m",
			Usage:       "max qty of main asset",
			Destination: &threshold,
		},
		cli.Float64Flag{
			Name:        "threshold, t",
			Usage:       "profit threshold",
			Destination: &threshold,
		},
	}

	app.Action = func(c *cli.Context) error {
		if apiKey == "" || secret == "" {
			return cli.NewExitError("api key and secret is required", 0)
		}
		return nil
	}

	app.Run(os.Args)

	mainAsset := models.Asset(assetString)
	exchange := CreateExchange(apiKey, secret, mainAsset, dryrun)
	analyzer := CreateAnalyzer(mainAsset, exchange.GetCharge(), maxqty, threshold)
	arbitrader := CreateTrader(exchange, analyzer, mainAsset)
	arbitrader.Run()
}

func CreateExchange(apikey string, secret string, mainAsset models.Asset, dryRun bool) usecase.Exchange {
	binance := infrastructure.NewBinance(
		apikey,
		secret,
	)

	if dryRun {
		balances := map[models.Asset]*models.Balance{}
		balances[mainAsset] = &models.Balance{
			Asset: mainAsset,
			Free:  0.12,
			Total: 0.12,
		}
		return infrastructure.NewExchangeStub(
			binance,
			balances,
		)
	}

	return binance
}

func CreateAnalyzer(mainAsset models.Asset, charge float64, maxqty float64, threshold float64) usecase.MarketAnalyzer {
	return usecase.NewMarketAnalyzer(
		mainAsset,
		charge,
		maxqty,
		threshold,
	)
}

func CreateTrader(exchange usecase.Exchange, analyzer usecase.MarketAnalyzer, mainAsset models.Asset) *usecase.Arbitrader {
	return usecase.NewArbitrader(
		exchange,
		analyzer,
		mainAsset,
	)
}