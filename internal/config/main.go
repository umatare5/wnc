package config

const (
	ControllersFlagName         = "controllers"
	AllowInsecureAccessFlagName = "insecure"
	PrintFormatFlagName         = "format"
	TimeoutFlagName             = "timeout"
	RadioFlagName               = "radio"
	SSIDFlagName                = "ssid"
	SortByFlagName              = "sort-by"
	SortOrderFlagName           = "sort-order"
	APNameFlagName              = "ap-name"
	PrintFormatJSON             = "json"
	PrintFormatTable            = "table"
	OrderByAscending            = "asc"
	OrderByDescending           = "desc"
	RadioSlotNumSlot0ID         = 0
	RadioSlotNumSlot1ID         = 1
	RadioSlotNumSlot2ID         = 2
)

type Config struct {
	GenerateCmdConfig GenerateCmdConfig
	ShowCmdConfig     ShowCmdConfig
}

func New() Config {
	return Config{
		GenerateCmdConfig: GenerateCmdConfig{},
		ShowCmdConfig:     ShowCmdConfig{},
	}
}
