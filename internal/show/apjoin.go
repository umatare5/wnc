package show

import (
	"context"
	"time"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/render"
	"github.com/umatare5/wnc/internal/wnc"
)

// Status display strings. The device prints these two words in "show wireless stats ap
// join summary", and naming them keeps the cell and the glyph mapping from drifting.
const (
	joinedYes = "Joined"
	joinedNo  = "Not Joined"
)

// APJoinRow is one row of show ap-join: one access point the controller remembers, joined or not.
// Status is a *string rather than the *bool the controller sends, because the device's own words
// are Joined and Not Joined where render.Bool would print Yes and No.
type APJoinRow struct {
	APName            *string `json:"ap_name,omitzero"`
	RadioMAC          *string `json:"radio_mac,omitzero"`
	EthernetMAC       *string `json:"ethernet_mac,omitzero"`
	IPAddress         *string `json:"ip_address,omitzero"`
	Status            *string `json:"status,omitzero"`
	LastFailurePhase  *string `json:"last_failure_phase,omitzero"`
	LastJoinFailure   *string `json:"last_join_failure,omitzero"`
	LastConfigFailure *string `json:"last_config_failure,omitzero"`
	LastDiscFailure   *string `json:"last_disc_failure,omitzero"`
	DisconnectReason  *string `json:"disconnect_reason,omitzero"`
	RebootReason      *string `json:"reboot_reason,omitzero"`
	LastJoin          *int64  `json:"last_join_seconds,omitzero"`
	LastConfig        *int64  `json:"last_config_seconds,omitzero"`
	LastDiscovery     *int64  `json:"last_discovery_seconds,omitzero"`
	LastError         *int64  `json:"last_error_seconds,omitzero"`
	Controller        string  `json:"controller"`
}

// APJoinColumns describes the join view. It is the only one that shows an access point absent from
// capwap-data, whose record here carries the phase that failed and the reason it disconnected.
func APJoinColumns() []render.Column[APJoinRow] {
	return []render.Column[APJoinRow]{
		{Key: keyAPName, Header: headAPName, Cell: func(r APJoinRow) string { return render.StrPtr(r.APName) }},
		{Key: keyRadioMAC, Header: headRadioMAC, Cell: func(r APJoinRow) string { return render.StrPtr(r.RadioMAC) }},
		{
			Key: keyEthernetMAC, Header: headEthernetMAC,
			Cell: func(r APJoinRow) string { return render.StrPtr(r.EthernetMAC) },
		},
		{
			Key:    keyIPAddress,
			Header: headIPAddress,
			Cell:   func(r APJoinRow) string { return render.StrPtr(r.IPAddress) },
		},
		{
			Key: keyStatus, Header: "Status",
			Cell:   func(r APJoinRow) string { return render.StrPtr(r.Status) },
			Pretty: func(r APJoinRow) string { return prettyState(r.Status, joinedYes, joinedNo, glyphBad) },
		},
		{
			Key: "last_failure_phase", Header: "Last Failure Phase",
			Cell: func(r APJoinRow) string { return render.StrPtr(r.LastFailurePhase) },
		},
		{
			Key: "last_join_failure", Header: "Last Join Failure",
			Cell: func(r APJoinRow) string { return render.StrPtr(r.LastJoinFailure) },
		},
		{
			Key: "last_config_failure", Header: "Last Config Failure",
			Cell: func(r APJoinRow) string { return render.StrPtr(r.LastConfigFailure) },
		},
		{
			Key: "last_disc_failure", Header: "Last Discovery Failure",
			Cell: func(r APJoinRow) string { return render.StrPtr(r.LastDiscFailure) },
		},
		{
			Key: "disconnect_reason", Header: "Last Disconnect Reason",
			Cell: func(r APJoinRow) string { return render.StrPtr(r.DisconnectReason) },
		},
		{
			Key: "reboot_reason", Header: "Reboot Reason",
			Cell: func(r APJoinRow) string { return render.StrPtr(r.RebootReason) },
		},
		{
			Key: "last_join_seconds", Header: "Last Join",
			Cell: func(r APJoinRow) string { return render.Duration(r.LastJoin) },
			Sort: func(r APJoinRow) any { return render.SortValue(r.LastJoin) },
		},
		{
			Key: "last_config_seconds", Header: "Last Config",
			Cell: func(r APJoinRow) string { return render.Duration(r.LastConfig) },
			Sort: func(r APJoinRow) any { return render.SortValue(r.LastConfig) },
		},
		{
			Key: "last_discovery_seconds", Header: "Last Discovery",
			Cell: func(r APJoinRow) string { return render.Duration(r.LastDiscovery) },
			Sort: func(r APJoinRow) any { return render.SortValue(r.LastDiscovery) },
		},
		{
			Key: "last_error_seconds", Header: "Last Error",
			Cell: func(r APJoinRow) string { return render.Duration(r.LastError) },
			Sort: func(r APJoinRow) any { return render.SortValue(r.LastError) },
		},
		{
			Key:    keyController,
			Header: headController,
			Cell:   func(r APJoinRow) string { return render.Str(r.Controller) },
		},
	}
}

// FetchAPJoins reads one controller's join view. One collection carries every column,
// so there is no secondary read to degrade and no join to get wrong.
func FetchAPJoins(ctx context.Context, c *wnc.Client, t config.Target, _ *Reporter) ([]APJoinRow, error) {
	joins, err := c.APJoins(ctx)
	if err != nil {
		return nil, err
	}

	return apJoinRows(joins, time.Now(), t), nil
}

// apJoinRows converts the join records. Every instant on ap-join-stats reads back as the Unix
// epoch when the event never happened, which render.SecondsSince rejects.
func apJoinRows(joins []wnc.APJoin, now time.Time, t config.Target) []APJoinRow {
	rows := make([]APJoinRow, 0, len(joins))

	for _, j := range joins {
		rows = append(rows, APJoinRow{
			APName:            optional(j.Name),
			RadioMAC:          optional(j.WtpMAC),
			EthernetMAC:       optional(j.EthernetMAC),
			IPAddress:         optional(j.IPAddr),
			Status:            joinStatus(j.Joined),
			LastFailurePhase:  optional(showJoinPhase(j.LastFailurePhase)),
			LastJoinFailure:   optional(showJoinFault(j.LastJoinFailure)),
			LastConfigFailure: optional(showJoinFault(j.LastConfigFailure)),
			LastDiscFailure:   optional(showJoinFault(j.LastDiscFailure)),
			DisconnectReason:  optional(j.DisconnectReason),
			RebootReason:      optional(showJoinFault(j.RebootReason)),
			LastJoin:          render.SecondsSince(now, j.LastJoin),
			LastConfig:        render.SecondsSince(now, j.LastConfig),
			LastDiscovery:     render.SecondsSince(now, j.LastDiscovery),
			LastError:         render.SecondsSince(now, j.LastError),
			Controller:        t.Name,
		})
	}

	return rows
}

// joinStatus renders the join flag in the device's own two words. A nil stays absent:
// a release that omits the leaf must not be reported as not joined.
func joinStatus(p *bool) *string {
	if p == nil {
		return nil
	}

	if *p {
		return ptr(joinedYes)
	}

	return ptr(joinedNo)
}
