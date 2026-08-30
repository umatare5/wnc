package show

import (
	"context"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/render"
	"github.com/umatare5/wnc/internal/wnc"
)

// OverviewRow is one row of show overview: one access point radio. ap_mac is the access point's
// base radio address and is the same on every radio of one access point, so the row identity is
// the (controller, ap_mac, slot) triple and the column is not called "radio mac".
type OverviewRow struct {
	APName     *string `json:"ap_name,omitzero"`
	APMAC      *string `json:"ap_mac,omitzero"`
	Slot       *int    `json:"slot,omitzero"`
	Mode       *string `json:"mode,omitzero"`
	Band       *string `json:"band,omitzero"`
	Admin      *string `json:"admin,omitzero"`
	Oper       *string `json:"oper,omitzero"`
	Channel    *int    `json:"channel,omitzero"`
	Width      *int    `json:"channel_width_mhz,omitzero"`
	TxPower    *int8   `json:"tx_power_dbm,omitzero"`
	Clients    *int    `json:"clients,omitzero"`
	ChUtil     *int    `json:"ch_util_percent,omitzero"`
	RFProfile  *string `json:"rf_profile,omitzero"`
	Controller string  `json:"controller"`
}

// OverviewColumns describes the per-radio view. Channel is a channel number and not a frequency,
// its YANG declaring no units, so the cell says "ch".
//
// Of the units here only dBm and the percentage are declared by the controller, on
// curr-tx-power-in-dbm and cca-util-percentage; MHz rests on chan-width's own enum descriptions,
// because the leaf is a bare uint8 and the module declares no Hz unit.
func OverviewColumns() []render.Column[OverviewRow] {
	return []render.Column[OverviewRow]{
		{Key: keyAPName, Header: headAPName, Cell: func(r OverviewRow) string { return render.StrPtr(r.APName) }},
		{Key: keyAPMAC, Header: "AP MAC", Cell: func(r OverviewRow) string { return render.StrPtr(r.APMAC) }},
		{
			Key: keySlot, Header: "Slot",
			Cell: func(r OverviewRow) string { return render.IntPtr(r.Slot) },
			Sort: func(r OverviewRow) any { return render.SortValue(r.Slot) },
		},
		{Key: keyMode, Header: "Mode", Cell: func(r OverviewRow) string { return render.StrPtr(r.Mode) }},
		{Key: keyBand, Header: "Band", Cell: func(r OverviewRow) string { return render.StrPtr(r.Band) }},
		{
			Key: keyAdmin, Header: "Admin",
			Cell:   func(r OverviewRow) string { return render.StrPtr(r.Admin) },
			Pretty: func(r OverviewRow) string { return prettyState(r.Admin, dispEnabled, dispDisabled, glyphNo) },
		},
		{
			Key: "oper", Header: "Oper",
			Cell:   func(r OverviewRow) string { return render.StrPtr(r.Oper) },
			Pretty: func(r OverviewRow) string { return prettyState(r.Oper, dispUp, dispDown, glyphBad) },
		},
		{
			Key: keyChannel, Header: "Channel",
			Cell: func(r OverviewRow) string { return render.UnitPtr(r.Channel, "ch") },
			Sort: func(r OverviewRow) any { return render.SortValue(r.Channel) },
		},
		{
			Key: "channel_width_mhz", Header: "Width",
			Cell: func(r OverviewRow) string { return render.UnitPtr(r.Width, "MHz") },
			Sort: func(r OverviewRow) any { return render.SortValue(r.Width) },
		},
		{
			Key: "tx_power_dbm", Header: "TxPower",
			Cell: func(r OverviewRow) string { return render.UnitPtr(r.TxPower, "dBm") },
			Sort: func(r OverviewRow) any { return render.SortValue(r.TxPower) },
		},
		{
			Key: "clients", Header: "Clients",
			Cell: func(r OverviewRow) string { return render.UnitPtr(r.Clients, "clients") },
			Sort: func(r OverviewRow) any { return render.SortValue(r.Clients) },
		},
		{
			Key: "ch_util_percent", Header: "ChUtil",
			Cell: func(r OverviewRow) string { return render.UnitPtr(r.ChUtil, "%") },
			Sort: func(r OverviewRow) any { return render.SortValue(r.ChUtil) },
		},
		{
			Key:    "rf_profile",
			Header: "RF Profile",
			Cell:   func(r OverviewRow) string { return render.StrPtr(r.RFProfile) },
		},
		{
			Key:    keyController,
			Header: headController,
			Cell:   func(r OverviewRow) string { return render.Str(r.Controller) },
		},
	}
}

// FetchOverview builds the fetcher for one band filter.
func FetchOverview(band string) Fetcher[OverviewRow] {
	return func(ctx context.Context, c *wnc.Client, t config.Target, rep *Reporter) ([]OverviewRow, error) {
		radios, reads, err := c.Radios(ctx)
		if err != nil {
			return nil, err
		}

		rep.Degraded("capwap-data", reads.CAPWAP)
		rep.Degraded("common-oper-data", reads.Clients)
		rep.Degraded("rrm-measurement", reads.RRM)
		rep.Degraded("rf-tags", reads.RFTags)

		return overviewRows(radios, band, t, rep), nil
	}
}

// overviewRows filters and converts the joined radios.
func overviewRows(radios []wnc.Radio, band string, t config.Target, rep *Reporter) []OverviewRow {
	rows := make([]OverviewRow, 0, len(radios))
	unreported := 0

	for _, r := range radios {
		display := showRadioBand(r.Band)

		if band != "" {
			if display == "" {
				unreported++

				continue
			}

			if display != band {
				continue
			}
		}

		rows = append(rows, OverviewRow{
			APName: optional(r.APName),
			APMAC:  optional(r.APMAC),
			// Zero is a slot number the controller reports, not an omitted leaf, so the
			// slot is not zero-sentinelled.
			Slot:       ptr(r.Slot),
			Mode:       optional(showRadioMode(r.Mode)),
			Band:       optional(display),
			Admin:      optional(showRadioAdmin(r.AdminState)),
			Oper:       optional(showRadioOper(r.OperState)),
			Channel:    r.Channel,
			Width:      r.Width,
			TxPower:    r.TxPowerDBm,
			Clients:    r.Clients,
			ChUtil:     r.ChUtil,
			RFProfile:  optional(r.RFProfile),
			Controller: t.Name,
		})
	}

	rep.Excluded(unreported, "the band was not reported for them")

	return rows
}
