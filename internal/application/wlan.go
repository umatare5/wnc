package application

import (
	"github.com/umatare5/cisco-ios-xe-wireless-go/wlan"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// WlanUsecase handles WLAN-related operations
type WlanUsecase struct {
	Config     *config.Config
	Repository *infrastructure.Repository
}

// ShowWlanData holds merged data from multiple controllers
type ShowWlanData struct {
	TagName      string            `json:"tag-name"`
	PolicyName   string            `json:"policy-name"`
	WlanName     string            `json:"wlan-name"`
	WlanCfgEntry wlan.WlanCfgEntry `json:"wlan-cfg-entry"`
	WlanPolicy   wlan.WlanPolicy   `json:"wlan-policy"`
	Controller   string            `json:"controller"`
}

// ShowWlan retrieves and merges WLAN data from multiple controllers
func (u *WlanUsecase) ShowWlan(controllers *[]config.Controller, isSecure *bool) []*ShowWlanData {
	var data []*ShowWlanData

	// Return empty slice if repository is nil
	if u.Repository == nil {
		return data
	}

	// Return empty slice if controllers is nil
	if controllers == nil {
		return data
	}

	for _, controller := range *controllers {
		wlanCfg := u.Repository.InvokeWlanRepository().GetWlanCfg(controller.Hostname, controller.AccessToken, isSecure)
		if wlanCfg == nil {
			// Skip this controller if authentication failed or other error occurred
			continue
		}

		for _, pol := range wlanCfg.CiscoIOSXEWirelessWlanCfgWlanCfgData.PolicyListEntries.PolicyListEntry {
			for _, p := range pol.WlanPolicies.WlanPolicy {
				var merged ShowWlanData
				merged.TagName = pol.TagName
				merged.Controller = controller.Hostname
				merged.PolicyName = p.PolicyProfileName
				merged.WlanName = p.WlanProfileName

				for _, d := range wlanCfg.CiscoIOSXEWirelessWlanCfgWlanCfgData.WlanCfgEntries.WlanCfgEntry {
					if d.ProfileName == merged.WlanName {
						merged.WlanCfgEntry = d
					}
				}
				for _, d := range wlanCfg.CiscoIOSXEWirelessWlanCfgWlanCfgData.WlanPolicies.WlanPolicy {
					if d.PolicyProfileName == merged.PolicyName {
						merged.WlanPolicy = d
					}
				}

				data = append(data, &merged)
			}
		}
	}

	return data
}
