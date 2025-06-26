package cli

import (
	"os"
	"sort"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
	"github.com/umatare5/wnc/pkg/tablewriter"
)

// ApTagCli struct
type ApTagCli struct {
	Config     *config.Config
	Repository *infrastructure.Repository
	Usecase    *application.Usecase
}

// ShowApTag retrieves the list of atcess points from the controllers
func (tc *ApTagCli) ShowApTag() {
	isSecure := !tc.Config.ShowCmdConfig.AllowInsecureAccess
	apTags := tc.Usecase.InvokeApUsecase().ShowApTag(
		&tc.Config.ShowCmdConfig.Controllers,
		&isSecure,
	)

	if isJSONFormat(tc.Config.ShowCmdConfig.PrintFormat) {
		printJson(apTags)
		return
	}

	tc.renderShowApTagTable(apTags)
}

// renderShowApTagTable renders the atcess point data in a table format
func (tc *ApTagCli) renderShowApTagTable(apTags []*application.ShowApTagData) {
	table := tablewriter.NewTable(os.Stdout)

	// Set table headers
	headers := tc.getShowApTagTableHeaders()
	table.Header(headers)

	// Set table rows
	tc.sortShowApTagRow(apTags)
	for _, apTag := range apTags {
		row, _ := tc.formatShowApTagRow(apTag)
		table.Append(row)
	}
	// Render the table
	_ = table.Render()
}

func (tc *ApTagCli) getShowApTagTableHeaders() []string {
	return []string{
		"AP Name", "Config", "Policy Tag Name", "RF Tag Name", "Site Tag Name",
		"AP Profile", "Flex Profile", "Tag Source",
	}
}

func (tc *ApTagCli) formatShowApTagRow(ap *application.ShowApTagData) ([]string, error) {
	row := []string{
		ap.CapwapData.Name,
		tc.convertCapwapTagInfoIsApMisconfigurationToConfigCheck(ap.CapwapData.TagInfo.IsApMisconfigured),
		ap.CapwapData.TagInfo.ResolvedTagInfo.ResolvedPolicyTag,
		ap.CapwapData.TagInfo.ResolvedTagInfo.ResolvedRfTag,
		ap.CapwapData.TagInfo.ResolvedTagInfo.ResolvedSiteTag,
		ap.CapwapData.TagInfo.SiteTag.ApProfile,
		ap.CapwapData.TagInfo.SiteTag.FlexProfile,
		ap.CapwapData.TagInfo.TagSource,
	}

	return row, nil
}

func (tc *ApTagCli) sortShowApTagRow(apTags []*application.ShowApTagData) {
	sort.Slice(apTags, func(i, j int) bool {
		return apTags[i].CapwapData.Name < apTags[j].CapwapData.Name
	})
}

func (tc *ApTagCli) convertCapwapTagInfoIsApMisconfigurationToConfigCheck(isMisconfigured bool) string {
	if isAPMisconfigured(isMisconfigured) {
		return "  ❌️"
	}
	return "  ✅️"
}
