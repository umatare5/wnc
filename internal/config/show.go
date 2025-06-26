package config

import (
	"errors"
	"strings"

	"github.com/jinzhu/configor"
	"github.com/umatare5/wnc/pkg/log"
	"github.com/urfave/cli/v3"
)

const (
	OverviewHeaderApMac              = "APMac"
	OverviewHeaderApRadioID          = "Radio"
	OverviewHeaderApOperStatus       = "Status"
	OverviewHeaderChannelNumber      = "Channel"
	OverviewHeaderChannelUtilization = "ChannelUtilization"
	OverviewHeaderClientCount        = "ClientCount"
	OverviewHeaderRFTagName          = "RFTagName"
	OverviewHeaderTxPower            = "TxPower"
	ShowClientHeaderBand             = "Band"
	ShowClientHeaderHostname         = "Hostname"
	ShowClientHeaderIP               = "IPAddress"
	ShowClientHeaderMacAddress       = "MACAddress"
	ShowClientHeaderProtocol         = "Protocol"
	ShowClientHeaderRSSI             = "RSSI"
	ShowClientHeaderRxTraffic        = "RxTraffic"
	ShowClientHeaderSNR              = "SNR"
	ShowClientHeaderSSID             = "SSID"
	ShowClientHeaderState            = "State"
	ShowClientHeaderStream           = "Stream"
	ShowClientHeaderThroughput       = "Throughput"
	ShowClientHeaderTxTraffic        = "TxTraffic"
	ShowClientHeaderUsername         = "Username"
	ShowCommonHeaderApName           = "APName"
	ShowCommonHeaderController       = "Controller"
)

// ShowCmdConfig holds show command configuration
type ShowCmdConfig struct {
	Controllers         []Controller
	AllowInsecureAccess bool
	PrintFormat         string
	Timeout             int
	APName              string
	Radio               string
	SSID                string
	SortBy              string
	SortOrder           string
}

type Controller struct {
	Hostname    string
	AccessToken string
}

// SetShowCmdConfig initializes the configuration
func (c *Config) SetShowCmdConfig(cli *cli.Command) {
	err := c.validateShowCmdFlags(cli)
	if err != nil {
		log.Fatal(err)
	}

	cfg := ShowCmdConfig{
		Controllers:         c.parseControllers(cli.String(ControllersFlagName)),
		AllowInsecureAccess: cli.Bool(AllowInsecureAccessFlagName),
		PrintFormat:         cli.String(PrintFormatFlagName),
		Timeout:             cli.Int(TimeoutFlagName),
		APName:              cli.String(APNameFlagName),
		Radio:               cli.String(RadioFlagName),
		SSID:                cli.String(SSIDFlagName),
		SortBy:              cli.String(SortByFlagName),
		SortOrder:           cli.String(SortOrderFlagName),
	}

	err = configor.New(&configor.Config{}).Load(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	c.ShowCmdConfig = cfg
}

// validateShowCmdFlags checks if the flags are valid
func (c *Config) validateShowCmdFlags(cli *cli.Command) error {
	if err := c.validateControllersFormat(cli.String(ControllersFlagName)); err != nil {
		log.Fatal(err)
	}
	if err := c.validatePrintFormat(cli.String(PrintFormatFlagName)); err != nil {
		log.Fatal(err)
	}

	return nil
}

// validateControllersFormat checks if the controllers format is valid
func (c *Config) validateControllersFormat(input string) error {
	pairs := strings.SplitSeq(input, ",")
	for pair := range pairs {
		hostname, accessToken, err := c.parseControllerPair(pair)
		if err != nil {
			return err
		}

		if hostname == "" {
			return errors.New("invalid controllers format: hostname is empty.")
		}
		if accessToken == "" {
			return errors.New("invalid controllers format: access token is empty.")
		}
	}
	return nil
}

// parseControllerPair parses a single controller:token pair, handling URLs with schemas
func (c *Config) parseControllerPair(pair string) (hostname, accessToken string, err error) {
	// Find the last colon to split hostname and token
	// This handles URLs like https://hostname:port:token correctly
	lastColonIndex := strings.LastIndex(pair, ":")
	if lastColonIndex == -1 {
		return "", "", errors.New("invalid controllers format: controllers does not contain ':'.")
	}

	hostname = strings.TrimSpace(pair[:lastColonIndex])
	accessToken = strings.TrimSpace(pair[lastColonIndex+1:])

	// Remove schema prefix if present (https:// or http://)
	if strings.HasPrefix(hostname, "https://") {
		hostname = strings.TrimPrefix(hostname, "https://")
	} else if strings.HasPrefix(hostname, "http://") {
		hostname = strings.TrimPrefix(hostname, "http://")
	}

	return hostname, accessToken, nil
}

// validatePrintFormat checks if the output format is valid
func (c *Config) validatePrintFormat(format string) error {
	switch format {
	case PrintFormatJSON, PrintFormatTable:
		return nil
	default:
		return errors.New(`invalid format: must be "json" or "table"`)
	}
}

// parseControllers parses the controllers flag into a slice of Controller structs
func (c *Config) parseControllers(input string) []Controller {
	pairs := strings.Split(input, ",")
	controllers := []Controller{}

	for _, pair := range pairs {
		hostname, accessToken, err := c.parseControllerPair(pair)
		if err != nil {
			// This should not happen as validation already passed
			continue
		}
		controllers = append(controllers, Controller{
			Hostname:    hostname,
			AccessToken: accessToken,
		})
	}

	return controllers
}
