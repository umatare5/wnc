package wnc

import (
	"context"
	"fmt"

	sdk "github.com/umatare5/cisco-ios-xe-wireless-go"
)

// All three reads below ask for the defaults in force. There is no fallback to a plain read: a
// controller that refuses with-defaults answers 400 and the read fails, because a plain answer
// would report a default as an absence.

type PolicyTag struct {
	Name        string
	Description *string
	Bindings    []PolicyBinding
}

// PolicyBinding is one WLAN profile the tag binds and the policy profile it is bound to; both
// leaves are keys of the wlan-policy list, so neither half occurs alone. This is not Binding, which
// names the tag as well because the WLAN view reaches these pairs from the WLAN side.
type PolicyBinding struct {
	WLANProfile   string
	PolicyProfile string
}

// SiteTag is one site tag's configuration. The five fields are exactly the leaves the SDK marks as
// seen on a live controller; the six it declares from the model alone are left out rather than
// rendered on a model's word.
type SiteTag struct {
	Name          string
	Description   *string
	APJoinProfile *string
	FlexProfile   *string

	// LocalSite is a pointer because its absence is not "not local": the leaf carries a
	// schema default, which is why this read asks for the defaults in force.
	LocalSite *bool
}

// RFTag is one RF tag's configuration. The per-slot rf-tag-radio-profiles list is neither read nor
// written, so the three band leaves are the whole of what this CLI carries for an RF tag.
type RFTag struct {
	Name         string
	Description  *string
	Profile24GHz *string
	Profile5GHz  *string
	Profile6GHz  *string
}

// PolicyTags reads the policy tags and the WLAN bindings each one carries.
func (c *Client) PolicyTags(ctx context.Context) ([]PolicyTag, error) {
	entries, err := c.sdk.PolicyTag().ListPolicyTags(ctx, sdk.WithDefaults(sdk.ReportAll))
	if err != nil {
		return nil, fmt.Errorf("reading policy-list-entries: %w", err)
	}

	out := make([]PolicyTag, 0, len(entries))

	for _, e := range entries {
		tag := PolicyTag{Name: e.TagName, Description: e.Description}

		// A tag binding nothing sends no container at all, which is the ordinary case for
		// the built-in default tag.
		if e.WLANPolicies != nil {
			tag.Bindings = make([]PolicyBinding, 0, len(e.WLANPolicies.WLANPolicy))

			for _, b := range e.WLANPolicies.WLANPolicy {
				tag.Bindings = append(tag.Bindings, PolicyBinding{
					WLANProfile:   b.WLANProfileName,
					PolicyProfile: b.PolicyProfileName,
				})
			}
		}

		out = append(out, tag)
	}

	return out, nil
}

func (c *Client) SiteTags(ctx context.Context) ([]SiteTag, error) {
	entries, err := c.sdk.SiteTag().ListSiteTags(ctx, sdk.WithDefaults(sdk.ReportAll))
	if err != nil {
		return nil, fmt.Errorf("reading site-tag-configs: %w", err)
	}

	out := make([]SiteTag, 0, len(entries))

	for _, e := range entries {
		out = append(out, SiteTag{
			Name:          e.SiteTagName,
			Description:   e.Description,
			APJoinProfile: e.ApJoinProfile,
			FlexProfile:   e.FlexProfile,
			LocalSite:     e.IsLocalSite,
		})
	}

	return out, nil
}

// RFTags reads the RF tags through the tag service rather than the accessor show overview uses.
// The two read the same collection and build different things: that one indexes the per-band names
// by (tag, band) to join onto a radio, and this one is the tag itself.
func (c *Client) RFTags(ctx context.Context) ([]RFTag, error) {
	entries, err := c.sdk.RFTag().ListRFTags(ctx, sdk.WithDefaults(sdk.ReportAll))
	if err != nil {
		return nil, fmt.Errorf("reading rf-tags: %w", err)
	}

	out := make([]RFTag, 0, len(entries))

	for _, e := range entries {
		out = append(out, RFTag{
			Name:         e.TagName,
			Description:  e.Description,
			Profile24GHz: e.Dot11BRfProfileName,
			Profile5GHz:  e.Dot11ARfProfileName,
			Profile6GHz:  e.Dot116GhzRFProfName,
		})
	}

	return out, nil
}
