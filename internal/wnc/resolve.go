package wnc

import (
	"context"

	sdk "github.com/umatare5/cisco-ios-xe-wireless-go"
)

// APRadioMACByName settles that a controller holds an access point under one name and returns the
// base radio address the row carries. ap-name-mac-map is keyed on wtp-name, so this reads one row
// and a name no access point holds answers 404; the bool is what that 404 cannot say, because a
// 200 carrying no row is reported as an absence rather than keying the radio read on an empty
// string.
func (c *Client) APRadioMACByName(
	ctx context.Context, apName string,
) (radioMAC string, found bool, err error) {
	resp, err := c.sdk.AP().GetNameMACMapByWTPName(ctx, apName)
	if err != nil {
		return "", false, err
	}

	if resp == nil || len(resp.ApNameMACMap) == 0 {
		return "", false, nil
	}

	return resp.ApNameMACMap[0].WtpMAC, true, nil
}

// ClientByMAC settles that a controller holds a client at one address and returns the address the
// row carries, which is the controller's own spelling and what goes on the wire.
func (c *Client) ClientByMAC(
	ctx context.Context, mac string,
) (clientMAC string, found bool, err error) {
	resp, err := c.sdk.Client().GetCommonInfoByMAC(ctx, mac)
	if err != nil {
		return "", false, err
	}

	if resp == nil || len(resp.CommonOperData) == 0 {
		return "", false, nil
	}

	return resp.CommonOperData[0].ClientMAC, true, nil
}

// The two leaves the username resolve joins on. common-oper-data is keyed on client-mac, so a
// username cannot be read by key and the collection is filtered instead, pruned to what the filter
// and the count need; none of the dropped leaves is credential material.
const clientUsernameFields = "client-mac;username"

// ClientsByUsername counts the sessions a controller holds under one username. The count is the
// whole answer because the RPC's user-name leaf states no cardinality, and the match is on equality
// because most records carry the leaf empty, so an empty argument would select nearly the whole
// fleet; the caller refuses one before this is reached.
func (c *Client) ClientsByUsername(ctx context.Context, username string) (int, error) {
	resp, err := c.sdk.Client().ListCommonInfo(ctx, sdk.WithFields(clientUsernameFields))
	if err != nil {
		return 0, err
	}

	if resp == nil {
		return 0, nil
	}

	matched := 0

	for _, cl := range resp.CommonOperData {
		if cl.Username == username {
			matched++
		}
	}

	return matched, nil
}
