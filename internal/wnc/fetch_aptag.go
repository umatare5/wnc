package wnc

import (
	"context"
	"fmt"

	sdk "github.com/umatare5/cisco-ios-xe-wireless-go"
	"github.com/umatare5/cisco-ios-xe-wireless-go/service/ap"
)

// The fields expression names the whole tag-info container rather than the leaves inside it: an
// expression naming a node the release does not declare answers 200 with a chunked body that stops
// mid-object and carries no Content-Length to detect it by. Pruning also keeps the serial number
// and the four certificate leaves on the controller.
const apTagFields = "wtp-mac;name;tag-info"

// APTag is one access point's tag assignment. The three tag names are the RESOLVED ones, what is in
// force; the two profile names have no resolved counterpart in the schema, so they come from the
// configured site-tag container and agree with SiteTag only while the configured and resolved site
// tags match.
type APTag struct {
	Name            string
	WtpMAC          string
	Misconfigured   *bool
	MisconfigReason *string
	TagSource       string
	FilterName      string
	PolicyTag       string
	SiteTag         string
	RFTag           string
	APProfile       string
	FlexProfile     string
}

// APTags reads every column of the tag view from one pruned collection; the controller's own
// identity on the row is the CLI's, not a leaf. An access point holding no tag information yields
// empty strings rather than an absence, because every container on this path is a non-pointer
// struct.
func (c *Client) APTags(ctx context.Context) ([]APTag, error) {
	resp, err := c.sdk.AP().ListCAPWAPData(ctx, sdk.WithFields(apTagFields))
	if err != nil {
		return nil, fmt.Errorf("reading capwap-data: %w", err)
	}

	if resp == nil {
		return nil, nil
	}

	out := make([]APTag, 0, len(resp.CAPWAPData))

	for _, record := range resp.CAPWAPData {
		tags := record.TagInfo

		out = append(out, APTag{
			Name:            record.Name,
			WtpMAC:          record.WtpMAC,
			Misconfigured:   tags.IsApMisconfigured,
			MisconfigReason: misconfigSpelling(tags.ApMisconfig),
			TagSource:       tags.TagSource,
			FilterName:      tags.FilterInfo.FilterName,
			PolicyTag:       tags.ResolvedTagInfo.ResolvedPolicyTag,
			SiteTag:         tags.ResolvedTagInfo.ResolvedSiteTag,
			RFTag:           tags.ResolvedTagInfo.ResolvedRFTag,
			APProfile:       tags.SiteTag.ApProfile,
			FlexProfile:     tags.SiteTag.FlexProfile,
		})
	}

	return out, nil
}

// misconfigSpelling carries the wire spelling out of the SDK's named type, because
// internal/show owns the display mapping and cannot import the SDK.
func misconfigSpelling(reason *ap.ApMisconfig) *string {
	if reason == nil {
		return nil
	}

	s := string(*reason)

	return &s
}
