package subcommand

import (
	"fmt"

	"github.com/umatare5/wnc/internal/config"
	"github.com/urfave/cli/v3"
)

// registerControllersFlag defines the flag for specifying controllers and access tokens.
func registerControllersFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     config.ControllersFlagName,
			Usage:    "Comma-separated list of controllers and their access tokens. Examples: 'wnc1.example.com:token1,wnc2.example.com:token2'",
			Required: true,
			Aliases:  []string{"c"},
			Sources:  cli.EnvVars("WNC_CONTROLLERS"),
		},
	}
}

// registerPrintFormatFlag defines the flag for specifying output format.
func registerPrintFormatFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name: config.PrintFormatFlagName,
			Usage: fmt.Sprintf(
				"Print format for the response. One of: [%s|%s]",
				config.PrintFormatJSON,
				config.PrintFormatTable,
			),
			Value:   config.PrintFormatTable,
			Aliases: []string{"f"},
		},
	}
}

// registerTimeoutFlag defines the flag for HTTP client timeout
func registerTimeoutFlag() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    config.TimeoutFlagName,
			Usage:   "HTTP client timeout in seconds",
			Value:   60,
			Aliases: []string{"t"},
		},
	}
}

func registerRadioFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name: config.RadioFlagName,
			Usage: fmt.Sprintf(
				"Radio interface number to filter the results. One of: [%d: 2.4GHz, %d: 5GHz, %d: 5GHz/6GHz]",
				config.RadioSlotNumSlot0ID,
				config.RadioSlotNumSlot1ID,
				config.RadioSlotNumSlot2ID,
			),
			Aliases: []string{"r"},
		},
	}
}

func registerSSIDFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    config.SSIDFlagName,
			Usage:   "ESSID name to filter the results.",
			Aliases: []string{"s"},
		},
	}
}

func registerClientSortByFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name: config.SortByFlagName,
			Usage: fmt.Sprintf(
				"Sort the results by a specific field. One of: [%s|%s|%s|%s|%s|%s|%s]",
				config.ShowClientHeaderHostname,
				config.ShowClientHeaderIP,
				config.ShowClientHeaderRSSI,
				config.ShowClientHeaderSNR,
				config.ShowClientHeaderThroughput,
				config.ShowClientHeaderRxTraffic,
				config.ShowClientHeaderTxTraffic,
			),
			Aliases: []string{"b"},
			Value:   config.ShowClientHeaderIP,
		},
	}
}

func registerOverviewSortByFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name: config.SortByFlagName,
			Usage: fmt.Sprintf(
				"Sort the results by a specific field. One of: [%s|%s|%s|%s|%s]",
				config.ShowCommonHeaderApName,
				config.OverviewHeaderApMac,
				config.OverviewHeaderChannelNumber,
				config.OverviewHeaderClientCount,
				config.OverviewHeaderTxPower,
			),
			Aliases: []string{"b"},
			Value:   config.ShowCommonHeaderApName,
		},
	}
}

func registerSortOrderFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name: config.SortOrderFlagName,
			Usage: fmt.Sprintf(
				"Sort the results by a specific pattern. One of: [%s|%s]",
				config.OrderByAscending,
				config.OrderByDescending,
			),
			Aliases: []string{"o"},
			Value:   config.OrderByDescending,
		},
	}
}

// registerInsecureFlag defines the flag for skipping TLS certificate verification.
func registerInsecureFlag() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    config.AllowInsecureAccessFlagName,
			Usage:   "Skip TLS certificate verification",
			Value:   false,
			Aliases: []string{"k"},
		},
	}
}
