package wnc

import (
	"context"

	"github.com/umatare5/cisco-ios-xe-wireless-go/service/rf"
	"github.com/umatare5/cisco-ios-xe-wireless-go/service/site"
	"github.com/umatare5/cisco-ios-xe-wireless-go/service/wlan"
)

// TagFields are the bindings one write carries. A nil field was not named on the command line and
// is left alone, because an empty string sent to mean "unset" would tell the operator a binding
// was cleared that was not.
type TagFields struct {
	Description   *string
	Profile24GHz  *string
	Profile5GHz   *string
	Profile6GHz   *string
	APJoinProfile *string
	FlexProfile   *string
	LocalSite     *bool
	WLAN          *string
	PolicyProfile *string
}

// Empty reports that no binding was named, which distinguishes a create carrying only the
// key from an update that would change nothing.
func (f TagFields) Empty() bool {
	return f.Description == nil && f.Profile24GHz == nil && f.Profile5GHz == nil &&
		f.Profile6GHz == nil && f.APJoinProfile == nil && f.FlexProfile == nil &&
		f.LocalSite == nil && f.WLAN == nil && f.PolicyProfile == nil
}

// tagExists turns a keyed tag read into a presence answer. The URL keys the record itself rather
// than a node inside it, so a name the controller does not hold answers an unambiguous 404 and
// every other failure stays a failure.
func tagExists[T any](tag *T, err error) (bool, error) {
	if err != nil {
		if cause, _ := Classify(err); cause == CauseNotFound {
			return false, nil
		}

		return false, err
	}

	return tag != nil, nil
}

func (c *Client) PolicyTagExists(ctx context.Context, name string) (bool, error) {
	return tagExists(c.sdk.PolicyTag().GetPolicyTag(ctx, name))
}

func (c *Client) SiteTagExists(ctx context.Context, name string) (bool, error) {
	return tagExists(c.sdk.SiteTag().GetSiteTag(ctx, name))
}

func (c *Client) RFTagExists(ctx context.Context, name string) (bool, error) {
	return tagExists(c.sdk.RFTag().GetRFTag(ctx, name))
}

// CreatePolicyTag creates a policy tag carrying only the fields that were named. The three
// configuration modules declare no leafref, so a named profile that does not exist is accepted
// and persists.
func (c *Client) CreatePolicyTag(ctx context.Context, name string, f TagFields) error {
	entry := wlan.PolicyListEntry{TagName: name, Description: f.Description}

	if f.WLAN != nil && f.PolicyProfile != nil {
		entry.WLANPolicies = &wlan.WLANPolicies{
			WLANPolicy: []wlan.WLANPolicyMap{{
				WLANProfileName:   *f.WLAN,
				PolicyProfileName: *f.PolicyProfile,
			}},
		}
	}

	return c.sdk.PolicyTag().CreatePolicyTag(ctx, &entry)
}

// UpdatePolicyTag applies the named fields one at a time. Each SDK setter reads the tag
// and writes it back, so a field the command did not name survives the exchange.
func (c *Client) UpdatePolicyTag(ctx context.Context, name string, f TagFields) error {
	if f.Description != nil {
		if err := c.sdk.PolicyTag().SetDescription(ctx, name, *f.Description); err != nil {
			return err
		}
	}

	if f.WLAN != nil && f.PolicyProfile != nil {
		return c.sdk.PolicyTag().SetPolicyProfile(ctx, name, *f.WLAN, *f.PolicyProfile)
	}

	return nil
}

func (c *Client) DeletePolicyTag(ctx context.Context, name string) error {
	return c.sdk.PolicyTag().DeletePolicyTag(ctx, name)
}

// CreateSiteTag creates a site tag. A flex profile forces is-local-site false rather than leaving
// it to the operator: the leaf carries when "../is-local-site = 'false'" and is-local-site defaults
// to TRUE, so a body naming a flex profile without it is a when-violation, and the SDK's own setter
// does the same on the update path.
func (c *Client) CreateSiteTag(ctx context.Context, name string, f TagFields) error {
	local := f.LocalSite
	if f.FlexProfile != nil {
		no := false
		local = &no
	}

	return c.sdk.SiteTag().CreateSiteTag(ctx, &site.SiteListEntry{
		SiteTagName:   name,
		Description:   f.Description,
		ApJoinProfile: f.APJoinProfile,
		FlexProfile:   f.FlexProfile,
		IsLocalSite:   local,
	})
}

// UpdateSiteTag applies the named fields. The flex profile is applied last because the
// SDK's setter for it also forces is-local-site to false — a flex profile is only in
// force on a non-local site — so applying --local-site afterwards would undo the binding.
func (c *Client) UpdateSiteTag(ctx context.Context, name string, f TagFields) error {
	if f.Description != nil {
		if err := c.sdk.SiteTag().SetDescription(ctx, name, *f.Description); err != nil {
			return err
		}
	}

	if f.APJoinProfile != nil {
		if err := c.sdk.SiteTag().SetAPJoinProfile(ctx, name, *f.APJoinProfile); err != nil {
			return err
		}
	}

	if f.LocalSite != nil {
		if err := c.sdk.SiteTag().SetLocalSite(ctx, name, *f.LocalSite); err != nil {
			return err
		}
	}

	if f.FlexProfile != nil {
		return c.sdk.SiteTag().SetFlexProfile(ctx, name, *f.FlexProfile)
	}

	return nil
}

func (c *Client) DeleteSiteTag(ctx context.Context, name string) error {
	return c.sdk.SiteTag().DeleteSiteTag(ctx, name)
}

// CreateRFTag creates an RF tag carrying only the fields that were named.
func (c *Client) CreateRFTag(ctx context.Context, name string, f TagFields) error {
	entry := rf.RFTag{
		TagName:             name,
		Description:         f.Description,
		Dot11BRfProfileName: f.Profile24GHz,
		Dot11ARfProfileName: f.Profile5GHz,
		Dot116GhzRFProfName: f.Profile6GHz,
	}

	return c.sdk.RFTag().CreateRFTag(ctx, &entry)
}

// UpdateRFTag applies the named fields one at a time. Each SDK setter reads the tag and PATCHes it
// back, as the site and policy setters do, so a field the command did not name survives.
func (c *Client) UpdateRFTag(ctx context.Context, name string, f TagFields) error {
	if f.Description != nil {
		if err := c.sdk.RFTag().SetDescription(ctx, name, *f.Description); err != nil {
			return err
		}
	}

	if f.Profile24GHz != nil {
		if err := c.sdk.RFTag().SetDot11BRfProfile(ctx, name, *f.Profile24GHz); err != nil {
			return err
		}
	}

	if f.Profile5GHz != nil {
		if err := c.sdk.RFTag().SetDot11ARfProfile(ctx, name, *f.Profile5GHz); err != nil {
			return err
		}
	}

	if f.Profile6GHz != nil {
		return c.sdk.RFTag().SetDot116GhzRFProfile(ctx, name, *f.Profile6GHz)
	}

	return nil
}

func (c *Client) DeleteRFTag(ctx context.Context, name string) error {
	return c.sdk.RFTag().DeleteRFTag(ctx, name)
}
