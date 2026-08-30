package show

import (
	"context"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/render"
	"github.com/umatare5/wnc/internal/wnc"
)

// The three views below list what set and delete write. show ap-tag reports only the tags
// resolved onto an access point, so it cannot show a tag that exists and is bound to nothing.

// PolicyTagRow is one row of show policy-tag: one policy tag paired with one WLAN binding it
// carries, so a tag binding three WLANs is three rows.
type PolicyTagRow struct {
	PolicyTag     *string `json:"policy_tag,omitzero"`
	Description   *string `json:"description,omitzero"`
	WLAN          *string `json:"wlan,omitzero"`
	PolicyProfile *string `json:"policy_profile,omitzero"`
	Controller    string  `json:"controller"`
}

// PolicyTagColumns describes the policy tag view. WLAN is the WLAN profile name the binding keys
// on rather than the SSID, and the two differ on any WLAN not named after its SSID.
func PolicyTagColumns() []render.Column[PolicyTagRow] {
	return []render.Column[PolicyTagRow]{
		{
			Key: DefaultSortPolicyTag, Header: headPolicyTag,
			Cell: func(r PolicyTagRow) string { return render.StrPtr(r.PolicyTag) },
		},
		{
			Key: keyDescription, Header: headDescription,
			Cell: func(r PolicyTagRow) string { return render.StrPtr(r.Description) },
		},
		{Key: "wlan", Header: "WLAN", Cell: func(r PolicyTagRow) string { return render.StrPtr(r.WLAN) }},
		{
			Key: keyPolicyProfile, Header: headPolicyProfile,
			Cell: func(r PolicyTagRow) string { return render.StrPtr(r.PolicyProfile) },
		},
		{
			Key: keyController, Header: headController,
			Cell: func(r PolicyTagRow) string { return render.Str(r.Controller) },
		},
	}
}

// SiteTagRow is one row of show site-tag.
type SiteTagRow struct {
	SiteTag       *string `json:"site_tag,omitzero"`
	Description   *string `json:"description,omitzero"`
	APJoinProfile *string `json:"ap_join_profile,omitzero"`
	FlexProfile   *string `json:"flex_profile,omitzero"`
	LocalSite     *bool   `json:"local_site,omitzero"`
	Controller    string  `json:"controller"`
}

// SiteTagColumns describes the site tag view. Local Site takes no glyph: neither reading is a
// fault or a feature switched off.
func SiteTagColumns() []render.Column[SiteTagRow] {
	return []render.Column[SiteTagRow]{
		{
			Key: DefaultSortSiteTag, Header: headSiteTag,
			Cell: func(r SiteTagRow) string { return render.StrPtr(r.SiteTag) },
		},
		{
			Key: keyDescription, Header: headDescription,
			Cell: func(r SiteTagRow) string { return render.StrPtr(r.Description) },
		},
		{
			Key: "ap_join_profile", Header: "AP Join Profile",
			Cell: func(r SiteTagRow) string { return render.StrPtr(r.APJoinProfile) },
		},
		{
			Key: keyFlexProfile, Header: headFlexProfile,
			Cell: func(r SiteTagRow) string { return render.StrPtr(r.FlexProfile) },
		},
		{
			Key: "local_site", Header: "Local Site",
			Cell: func(r SiteTagRow) string { return render.Bool(r.LocalSite) },
			Sort: func(r SiteTagRow) any {
				if r.LocalSite == nil {
					return nil
				}

				return *r.LocalSite
			},
		},
		{
			Key: keyController, Header: headController,
			Cell: func(r SiteTagRow) string { return render.Str(r.Controller) },
		},
	}
}

// RFTagRow is one row of show rf-tag.
type RFTagRow struct {
	RFTag        *string `json:"rf_tag,omitzero"`
	Description  *string `json:"description,omitzero"`
	Profile24GHz *string `json:"profile_24ghz,omitzero"`
	Profile5GHz  *string `json:"profile_5ghz,omitzero"`
	Profile6GHz  *string `json:"profile_6ghz,omitzero"`
	Controller   string  `json:"controller"`
}

// RFTagColumns describes the RF tag view.
//
// The 2.4 GHz column is the 802.11b leaf and the 5 GHz one the 802.11a leaf, which is the
// pairing the write path's own setters name and the one place these two can be swapped.
func RFTagColumns() []render.Column[RFTagRow] {
	return []render.Column[RFTagRow]{
		{
			Key: DefaultSortRFTag, Header: headRFTag,
			Cell: func(r RFTagRow) string { return render.StrPtr(r.RFTag) },
		},
		{
			Key: keyDescription, Header: headDescription,
			Cell: func(r RFTagRow) string { return render.StrPtr(r.Description) },
		},
		{
			Key: "profile_24ghz", Header: "2.4 GHz Profile",
			Cell: func(r RFTagRow) string { return render.StrPtr(r.Profile24GHz) },
		},
		{
			Key: "profile_5ghz", Header: "5 GHz Profile",
			Cell: func(r RFTagRow) string { return render.StrPtr(r.Profile5GHz) },
		},
		{
			Key: "profile_6ghz", Header: "6 GHz Profile",
			Cell: func(r RFTagRow) string { return render.StrPtr(r.Profile6GHz) },
		},
		{
			Key: keyController, Header: headController,
			Cell: func(r RFTagRow) string { return render.Str(r.Controller) },
		},
	}
}

// FetchPolicyTags reads one controller's policy tags. One collection carries every
// column, so there is no secondary read to degrade and no join to get wrong.
func FetchPolicyTags(ctx context.Context, c *wnc.Client, t config.Target, _ *Reporter) ([]PolicyTagRow, error) {
	tags, err := c.PolicyTags(ctx)
	if err != nil {
		return nil, err
	}

	return policyTagRows(tags, t), nil
}

func FetchSiteTags(ctx context.Context, c *wnc.Client, t config.Target, _ *Reporter) ([]SiteTagRow, error) {
	tags, err := c.SiteTags(ctx)
	if err != nil {
		return nil, err
	}

	return siteTagRows(tags, t), nil
}

func FetchRFTags(ctx context.Context, c *wnc.Client, t config.Target, _ *Reporter) ([]RFTagRow, error) {
	tags, err := c.RFTags(ctx)
	if err != nil {
		return nil, err
	}

	return rfTagRows(tags, t), nil
}

// policyTagRows expands each tag over its bindings. A tag binding nothing keeps a row of
// its own: it exists, it is a candidate for a delete, and dropping it is what an inner
// join would do.
func policyTagRows(tags []wnc.PolicyTag, t config.Target) []PolicyTagRow {
	rows := make([]PolicyTagRow, 0, len(tags))

	for _, tag := range tags {
		row := PolicyTagRow{
			PolicyTag:   optional(tag.Name),
			Description: optionalStr(tag.Description),
			Controller:  t.Name,
		}

		if len(tag.Bindings) == 0 {
			rows = append(rows, row)

			continue
		}

		for _, b := range tag.Bindings {
			bound := row
			bound.WLAN = optional(b.WLANProfile)
			bound.PolicyProfile = optional(b.PolicyProfile)

			rows = append(rows, bound)
		}
	}

	return rows
}

func siteTagRows(tags []wnc.SiteTag, t config.Target) []SiteTagRow {
	rows := make([]SiteTagRow, 0, len(tags))

	for _, tag := range tags {
		rows = append(rows, SiteTagRow{
			SiteTag:       optional(tag.Name),
			Description:   optionalStr(tag.Description),
			APJoinProfile: optionalStr(tag.APJoinProfile),
			FlexProfile:   optionalStr(tag.FlexProfile),
			LocalSite:     tag.LocalSite,
			Controller:    t.Name,
		})
	}

	return rows
}

func rfTagRows(tags []wnc.RFTag, t config.Target) []RFTagRow {
	rows := make([]RFTagRow, 0, len(tags))

	for _, tag := range tags {
		rows = append(rows, RFTagRow{
			RFTag:        optional(tag.Name),
			Description:  optionalStr(tag.Description),
			Profile24GHz: optionalStr(tag.Profile24GHz),
			Profile5GHz:  optionalStr(tag.Profile5GHz),
			Profile6GHz:  optionalStr(tag.Profile6GHz),
			Controller:   t.Name,
		})
	}

	return rows
}
