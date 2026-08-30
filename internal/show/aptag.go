package show

import (
	"context"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/render"
	"github.com/umatare5/wnc/internal/wnc"
)

// APTagRow is one row of show ap-tag.
type APTagRow struct {
	APName          *string `json:"ap_name,omitzero"`
	APMAC           *string `json:"ap_mac,omitzero"`
	Misconfigured   *bool   `json:"misconfigured,omitzero"`
	MisconfigReason *string `json:"misconfig_reason,omitzero"`
	TagSource       *string `json:"tag_source,omitzero"`
	FilterName      *string `json:"filter_name,omitzero"`
	PolicyTag       *string `json:"policy_tag,omitzero"`
	SiteTag         *string `json:"site_tag,omitzero"`
	RFTag           *string `json:"rf_tag,omitzero"`
	APProfile       *string `json:"ap_profile,omitzero"`
	FlexProfile     *string `json:"flex_profile,omitzero"`
	Controller      string  `json:"controller"`
}

// APTagColumns describes the tag view. Policy, Site and RF Tag are the resolved values in force;
// AP Profile and Flex Profile have no resolved counterpart in the schema, so they are the
// configured site tag's own and agree with Site Tag only while the two site tags do.
func APTagColumns() []render.Column[APTagRow] {
	return []render.Column[APTagRow]{
		{Key: keyAPName, Header: headAPName, Cell: func(r APTagRow) string { return render.StrPtr(r.APName) }},
		{Key: keyAPMAC, Header: "AP MAC", Cell: func(r APTagRow) string { return render.StrPtr(r.APMAC) }},
		{
			Key: "misconfigured", Header: "Misconfigured",
			Cell: func(r APTagRow) string { return render.Bool(r.Misconfigured) },
			// The polarity is inverted here: true is the fault, so it takes the cross.
			Pretty: func(r APTagRow) string { return prettyBool(r.Misconfigured, glyphBad, glyphOK) },
			Sort: func(r APTagRow) any {
				if r.Misconfigured == nil {
					return nil
				}

				return *r.Misconfigured
			},
		},
		{
			// The reason is the enum's own account of the flag beside it, and its
			// "no misconfiguration" member is a value rather than an absence. A dash
			// therefore means the release does not report the leaf at all.
			Key: "misconfig_reason", Header: "Misconfig Reason",
			Cell: func(r APTagRow) string { return render.StrPtr(r.MisconfigReason) },
		},
		{Key: "tag_source", Header: "Tag Source", Cell: func(r APTagRow) string { return render.StrPtr(r.TagSource) }},
		{
			// Named beside Tag Source because it answers it: a filter-sourced tag comes from
			// this rule. The controller sends an empty string where no filter exists and
			// optional collapses it, so a dash is "no filter" rather than a failed read.
			Key: "filter_name", Header: "Filter Name",
			Cell: func(r APTagRow) string { return render.StrPtr(r.FilterName) },
		},
		{Key: keyPolicyTag, Header: headPolicyTag, Cell: func(r APTagRow) string { return render.StrPtr(r.PolicyTag) }},
		{Key: keySiteTag, Header: headSiteTag, Cell: func(r APTagRow) string { return render.StrPtr(r.SiteTag) }},
		{Key: keyRFTag, Header: headRFTag, Cell: func(r APTagRow) string { return render.StrPtr(r.RFTag) }},
		{Key: "ap_profile", Header: "AP Profile", Cell: func(r APTagRow) string { return render.StrPtr(r.APProfile) }},
		{
			Key:    keyFlexProfile,
			Header: headFlexProfile,
			Cell:   func(r APTagRow) string { return render.StrPtr(r.FlexProfile) },
		},
		{Key: keyController, Header: headController, Cell: func(r APTagRow) string { return render.Str(r.Controller) }},
	}
}

// FetchAPTags reads one controller's tag view. One collection carries every column,
// so there is no secondary read to degrade and no join to get wrong.
func FetchAPTags(ctx context.Context, c *wnc.Client, t config.Target, _ *Reporter) ([]APTagRow, error) {
	tags, err := c.APTags(ctx)
	if err != nil {
		return nil, err
	}

	return apTagRows(tags, t), nil
}

// misconfigReason maps the reason through its domain while keeping an absent leaf apart from the
// domain's own "no misconfiguration" member.
func misconfigReason(v *string) *string {
	if v == nil {
		return nil
	}

	return optional(showAPMisconfig(*v))
}

func apTagRows(tags []wnc.APTag, t config.Target) []APTagRow {
	rows := make([]APTagRow, 0, len(tags))

	for _, tag := range tags {
		rows = append(rows, APTagRow{
			APName:          optional(tag.Name),
			APMAC:           optional(tag.WtpMAC),
			Misconfigured:   tag.Misconfigured,
			MisconfigReason: misconfigReason(tag.MisconfigReason),
			TagSource:       optional(showTagSource(tag.TagSource)),
			FilterName:      optional(tag.FilterName),
			PolicyTag:       optional(tag.PolicyTag),
			SiteTag:         optional(tag.SiteTag),
			RFTag:           optional(tag.RFTag),
			APProfile:       optional(tag.APProfile),
			FlexProfile:     optional(tag.FlexProfile),
			Controller:      t.Name,
		})
	}

	return rows
}
