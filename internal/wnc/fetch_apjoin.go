package wnc

import (
	"context"
	"fmt"
	"time"
)

// APJoin is one access point's join outcome, whether or not it is joined. This is the one view
// that reports an access point capwap-data cannot show, and every counter on the list is omitted
// because nothing on it declares a clear time, so a total since boot has no window.
type APJoin struct {
	Name              string
	WtpMAC            string
	EthernetMAC       string
	IPAddr            string
	Joined            *bool
	LastFailurePhase  string
	LastJoinFailure   string
	LastConfigFailure string
	LastDiscFailure   string
	DisconnectReason  string
	RebootReason      string
	LastJoin          time.Time
	LastConfig        time.Time
	LastDiscovery     time.Time
	LastError         time.Time
}

// APJoins reads the join view. One collection carries every column, so there is
// nothing to join and no secondary read to degrade.
//
// A controller with no access point answers 204, which the SDK returns as a non-nil
// pointer to a zero struct rather than as nil, so the empty slice below is what an
// empty fleet produces.
func (c *Client) APJoins(ctx context.Context) ([]APJoin, error) {
	resp, err := c.sdk.AP().ListAPJoinStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading ap-join-stats: %w", err)
	}

	if resp == nil {
		return nil, nil
	}

	out := make([]APJoin, 0, len(resp.ApJoinStats))

	for _, ap := range resp.ApJoinStats {
		join := ap.ApJoinInfo
		disc := ap.ApDiscoveryInfo

		out = append(out, APJoin{
			Name:              join.ApName,
			WtpMAC:            ap.WtpMAC,
			EthernetMAC:       join.ApEthernetMAC,
			IPAddr:            join.ApIPAddr,
			Joined:            join.IsJoined,
			LastFailurePhase:  join.LastErrorType,
			LastJoinFailure:   join.LastJoinFailureType,
			LastConfigFailure: join.LastConfigFailureType,
			LastDiscFailure:   disc.LastDiscFailureType,
			DisconnectReason:  ap.ApDisconnectReason,
			RebootReason:      ap.RebootReason,
			LastJoin:          join.LastSuccJoinAtmptTime,
			LastConfig:        join.LastSuccConfAtmptTime,
			LastDiscovery:     disc.LastSuccessDiscTime,
			LastError:         join.LastErrorTime,
		})
	}

	return out, nil
}
