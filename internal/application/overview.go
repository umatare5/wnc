package application

import (
	"fmt"

	"github.com/umatare5/cisco-ios-xe-wireless-go/ap"
	"github.com/umatare5/cisco-ios-xe-wireless-go/rf"
	"github.com/umatare5/cisco-ios-xe-wireless-go/rrm"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// OverviewUsecase handles overview operations
type OverviewUsecase struct {
	Config     *config.Config
	Repository *infrastructure.Repository
}

type ShowOverviewData struct {
	ApMac          string             `json:"ap-mac"`
	SlotID         int                `json:"slot-id"`
	Controller     string             `json:"controller"`
	RrmMeasurement rrm.RrmMeasurement `json:"rrm-measurement"`
	RadioOperData  ap.RadioOperData   `json:"radio-oper-data"`
	CapwapData     ap.CapwapData      `json:"capwap-ap"`
	RfTag          rf.RfTag           `json:"rf-tag"`
}

func (ou *OverviewUsecase) ShowOverview(controllers *[]config.Controller, isSecure *bool) []*ShowOverviewData {
	// Initialize with empty slice, not nil slice
	data := []*ShowOverviewData{}

	// Handle nil inputs gracefully
	if controllers == nil || ou.Repository == nil {
		return data
	}

	for _, controller := range *controllers {
		radioOperData := ou.Repository.InvokeApRepository().GetApRadioOperData(controller.Hostname, controller.AccessToken, isSecure)
		if radioOperData == nil {
			// Skip this controller if authentication failed or other error occurred
			continue
		}

		apCapwapData := ou.Repository.InvokeApRepository().GetApCapwapData(controller.Hostname, controller.AccessToken, isSecure)
		if apCapwapData == nil {
			// Skip this controller if authentication failed or other error occurred
			continue
		}

		rfTagsData := ou.Repository.InvokeRfRepository().GetRfTags(controller.Hostname, controller.AccessToken, isSecure)
		if rfTagsData == nil {
			// Skip this controller if authentication failed or other error occurred
			continue
		}

		rrmOperData := ou.Repository.InvokeRrmRepository().GetRrmMeasurement(controller.Hostname, controller.AccessToken, isSecure)
		if rrmOperData == nil {
			// Skip this controller if authentication failed or other error occurred
			continue
		}

		for _, radio := range radioOperData.RadioOperData {
			var merged ShowOverviewData
			merged.ApMac = radio.WtpMac
			merged.SlotID = radio.SlotID
			merged.Controller = controller.Hostname
			merged.RadioOperData = radio

			for _, d := range apCapwapData.CapwapData {
				if d.WtpMac == radio.WtpMac {
					merged.CapwapData = d
				}
			}

			for _, d := range rfTagsData.RfTags.RfTag {
				if d.TagName == merged.CapwapData.TagInfo.ResolvedTagInfo.ResolvedRfTag {
					merged.RfTag = d
				}
			}

			for _, d := range rrmOperData.RrmMeasurement {
				if d.WtpMac == radio.WtpMac {
					if d.RadioSlotID == radio.SlotID {
						merged.RrmMeasurement = d
					}
				}
			}

			data = append(data, &merged)
		}
	}
	return ou.filterByRadio(data)
}

func (ou *OverviewUsecase) filterByRadio(data []*ShowOverviewData) []*ShowOverviewData {
	// Handle nil config gracefully
	if ou.Config == nil {
		return data
	}

	filter := ou.Config.ShowCmdConfig.Radio
	if filter == "" {
		return data
	}
	// Initialize with empty slice, not nil slice
	filteredData := []*ShowOverviewData{}
	for _, d := range data {
		if fmt.Sprintf("%d", d.SlotID) == filter {
			filteredData = append(filteredData, d)
		}
	}
	return filteredData
}
