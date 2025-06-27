package show

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
	"github.com/umatare5/wnc/pkg/tablewriter"
)

// OverviewCli struct
type OverviewCli struct {
	Config     *config.Config
	Repository *infrastructure.Repository
	Usecase    *application.Usecase
}

// ShowOverview retrieves the list of atcess points from the controllers
func (oc *OverviewCli) ShowOverview() {
	isSecure := !oc.Config.ShowCmdConfig.AllowInsecureAccess
	data := oc.Usecase.InvokeOverviewUsecase().ShowOverview(
		&oc.Config.ShowCmdConfig.Controllers,
		&isSecure,
	)

	if oc.Config.ShowCmdConfig.PrintFormat == config.PrintFormatJSON {
		printJson(data)
		return
	}

	// Skip table rendering if no data is available
	if len(data) == 0 {
		return
	}

	oc.renderShowOverviewTable(data)
}

// renderShowOverviewTable renders the atcess point data in a table format
func (oc *OverviewCli) renderShowOverviewTable(data []*application.ShowOverviewData) {
	table := tablewriter.NewTable(os.Stdout)

	// Set table headers
	headers := oc.getShowOverviewTableHeaders()
	table.Header(headers)

	// Set table rows
	oc.sortShowOverviewRow(data)
	for _, Overview := range data {
		row, _ := oc.formatShowOverviewRow(Overview)
		table.Append(row)
	}
	// Render the table
	_ = table.Render()
}

func (oc *OverviewCli) getShowOverviewTableHeaders() []string {
	return []string{
		config.ShowCommonHeaderApName,
		config.OverviewHeaderApMac,
		config.OverviewHeaderApRadioID,
		config.OverviewHeaderApOperStatus,
		config.OverviewHeaderChannelNumber,
		config.OverviewHeaderTxPower,
		config.OverviewHeaderClientCount,
		config.OverviewHeaderChannelUtilization,
		config.OverviewHeaderRFTagName,
		config.ShowCommonHeaderController,
	}
}

func (oc *OverviewCli) formatShowOverviewRow(data *application.ShowOverviewData) ([]string, error) {
	var row []string

	// Check if essential data is available to avoid index out of range panics
	powerValue := "N/A"
	if len(data.RadioOperData.RadioBandInfo) > 0 {
		powerValue = fmt.Sprintf("%d dBm", data.RadioOperData.RadioBandInfo[0].PhyTxPwrLvlCfg.PhyTxPwrLvlCfgCfgData.CurrTxPowerInDbm)
	}

	row = []string{
		data.CapwapData.Name,
		data.CapwapData.WtpMac,
		fmt.Sprintf("%d", data.SlotID),
		oc.convertRadioOperDataOperState(data.RadioOperData.OperState),
		fmt.Sprintf("%d MHz %s",
			data.RadioOperData.PhyHtCfg.PhyHtCfgCfgData.ChanWidth,
			data.RadioOperData.PhyHtCfg.PhyHtCfgCfgData.FreqString,
		),
		powerValue,
		oc.convertRrmMeasurementLoadStations(data.RrmMeasurement.Load.Stations),
		oc.convertUtilizationsToIndicator(
			data.RrmMeasurement.Load.RxUtilPercentage,
			data.RrmMeasurement.Load.TxUtilPercentage,
			data.RrmMeasurement.Load.RxNoiseChannelUtilization),
	}

	if data.SlotID == config.RadioSlotNumSlot0ID {
		row = append(row, []string{
			data.RfTag.Dot11BRfProfileName,
		}...)
	}
	if data.SlotID == config.RadioSlotNumSlot1ID {
		row = append(row, []string{
			data.RfTag.Dot11ARfProfileName,
		}...)
	}
	if data.SlotID == config.RadioSlotNumSlot2ID {
		row = append(row, []string{
			data.RfTag.Dot116GhzRfProfName,
		}...)
	}

	row = append(row, []string{data.Controller}...)

	return row, nil
}

func (oc *OverviewCli) sortShowOverviewRow(data []*application.ShowOverviewData) {
	sort.Slice(data, func(i, j int) bool {
		sortBy := oc.Config.ShowCmdConfig.SortBy
		sortOrder := oc.Config.ShowCmdConfig.SortOrder

		var d bool
		switch sortBy {
		case config.ShowCommonHeaderApName:
			d = data[i].CapwapData.Name < data[j].CapwapData.Name
		case config.OverviewHeaderApMac:
			d = data[i].CapwapData.WtpMac < data[j].CapwapData.WtpMac
		case config.OverviewHeaderChannelNumber:
			d = data[i].RadioOperData.PhyHtCfg.PhyHtCfgCfgData.ChanWidth < data[j].RadioOperData.PhyHtCfg.PhyHtCfgCfgData.ChanWidth
		case config.OverviewHeaderTxPower:
			d = data[i].RadioOperData.RadioBandInfo[0].PhyTxPwrLvlCfg.PhyTxPwrLvlCfgCfgData.CurrTxPowerInDbm < data[j].RadioOperData.RadioBandInfo[0].PhyTxPwrLvlCfg.PhyTxPwrLvlCfgCfgData.CurrTxPowerInDbm
		case config.OverviewHeaderClientCount:
			d = data[i].RrmMeasurement.Load.Stations < data[j].RrmMeasurement.Load.Stations
		default:
			d = false
		}
		if sortOrder == config.OrderByDescending {
			return !d
		}
		return d
	})
}

func (oc *OverviewCli) convertUtilizationsToIndicator(rx, tx, noise int) string {
	total := rx + tx + noise
	if total > 100 {
		total = 100 // Cap total at 100%
	} else if total < 0 {
		total = 0 // Ensure total is not negative
	}
	filled := total / 10
	empty := 10 - filled
	return fmt.Sprintf("[%s%s] %d%%", strings.Repeat("#", filled), strings.Repeat(" ", empty), total)
}

func (oc *OverviewCli) convertRadioOperDataOperState(v string) string {
	if v == "radio-up" {
		return "  ✅️"
	}
	return "  ❌️"
}

func (oc *OverviewCli) convertRrmMeasurementLoadStations(v int) string {
	return fmt.Sprintf("%d clients", v)
}
