package application

import (
	"github.com/umatare5/cisco-ios-xe-wireless-go/ap"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// ApUsecase handles AP-related operations
type ApUsecase struct {
	Config     *config.Config
	Repository *infrastructure.Repository
}

type ShowApCommonData struct {
	ApMac      string        `json:"ap-mac"`
	Controller string        `json:"controller"`
	CapwapData ap.CapwapData `json:"capwap-ap"`
}

type ShowApData struct {
	ShowApCommonData
	LLDPnei    ap.LldpNeigh  `json:"lldp-neigh"`
	ApOperData ap.ApOperData `json:"ap-oper-data"`
}

type ShowApTagData struct {
	ShowApCommonData
}

// ShowAp retrieves and merges AP ap from multiple controllers
func (au *ApUsecase) ShowAp(controllers *[]config.Controller, isSecure *bool) []*ShowApData {
	var data []*ShowApData

	// Return empty slice if repository is nil
	if au.Repository == nil {
		return data
	}

	// Return empty slice if controllers is nil
	if controllers == nil {
		return data
	}

	for _, controller := range *controllers {
		apCapwapData := au.Repository.InvokeApRepository().GetApCapwapData(controller.Hostname, controller.AccessToken, isSecure)
		if apCapwapData == nil {
			// Skip this controller if authentication failed or other error occurred
			continue
		}

		apOperData := au.Repository.InvokeApRepository().GetApOperData(controller.Hostname, controller.AccessToken, isSecure)
		if apOperData == nil {
			// Skip this controller if authentication failed or other error occurred
			continue
		}

		lldp := au.Repository.InvokeApRepository().GetApLldpNeigh(controller.Hostname, controller.AccessToken, isSecure)
		if lldp == nil {
			// Skip this controller if authentication failed or other error occurred
			continue
		}

		for _, ap := range apCapwapData.CapwapData {
			var merged ShowApData
			merged.ApMac = ap.WtpMac
			merged.Controller = controller.Hostname
			merged.CapwapData = ap

			// Search ApLldpNeigh
			for _, d := range lldp.LldpNeigh {
				if d.WtpMac == ap.WtpMac {
					merged.LLDPnei = d
				}
			}

			// Search ApLldpNeigh
			for _, d := range apOperData.OperData {
				if d.WtpMac == ap.WtpMac {
					merged.ApOperData = d
				}
			}

			data = append(data, &merged)
		}
	}

	return data
}

// ShowAp retrieves and merges AP ap from multiple controllers
func (au *ApUsecase) ShowApTag(controllers *[]config.Controller, isSecure *bool) []*ShowApTagData {
	var data []*ShowApTagData

	// Return empty slice if repository is nil
	if au.Repository == nil {
		return data
	}

	// Return empty slice if controllers is nil
	if controllers == nil {
		return data
	}

	for _, controller := range *controllers {
		result := au.Repository.InvokeApRepository().GetApCapwapData(controller.Hostname, controller.AccessToken, isSecure)
		if result == nil {
			// Skip this controller if authentication failed or other error occurred
			continue
		}

		for _, ap := range result.CapwapData {
			var merged ShowApTagData
			merged.ApMac = ap.WtpMac
			merged.Controller = controller.Hostname
			merged.CapwapData = ap
			data = append(data, &merged)
		}
	}
	return data
}
