package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/wnc"
)

// adminVerb is the one thing that differs between the two trees. It carries no past
// participle on purpose: neither RPC declares an output container, so nothing this CLI
// prints may claim the radio is disabled rather than that the instruction was sent.
type adminVerb struct {
	name   string
	title  string
	gerund string
	on     bool
}

func adminVerbs() []adminVerb {
	return []adminVerb{
		{name: "enable", title: "Enable", gerund: "enabling", on: true},
		{name: "disable", title: "Disable", gerund: "disabling", on: false},
	}
}

// enableCommand and disableCommand are generated from one declaration so the two trees
// cannot disagree about which targets exist, as the tag trees are.
func enableCommand() *cli.Command { return adminCommand(adminVerbs()[0]) }

func disableCommand() *cli.Command { return adminCommand(adminVerbs()[1]) }

func adminCommand(v adminVerb) *cli.Command {
	return &cli.Command{
		Name:     v.name,
		Usage:    v.title + " an access point or one of its radios",
		Flags:    execFlags(),
		Action:   parentAction,
		Commands: []*cli.Command{adminAPLeaf(v), adminRadioLeaf(v)},
	}
}

// adminAPLeaf sets the access point's own admin state. It reuses the reset tree's runner
// because the guard sequence is identical: one name, one controller, resolve it against the
// controller, then act.
func adminAPLeaf(v adminVerb) *cli.Command {
	return &cli.Command{
		Name:      leafAP,
		Usage:     v.title + " one access point",
		UsageText: synopsis(config.FlagAPName),
		Description: "--ap-name is the ap_name column of wnc show ap. The controller resolves it\n" +
			"first, so a name it holds no access point under is refused before the RPC.\n" +
			"This sets the access point's admin state, not one radio's.\n" +
			"Pass --dry-run to name the target and change nothing.",
		Flags: []cli.Flag{apNameFlag(), yesFlag()},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runAPAction(ctx, cmd, writeWording{
				question: v.title + " %s on %s? This sets the access point's admin state, " +
					"not one radio's. [y/N]: ",
				dryRun: "%s on %s: would " + v.name + "\n",
				sent:   "%s on %s: " + v.name + " sent\n",
				failed: v.gerund + " %s on %s",
			}, func(ctx context.Context, c *wnc.Client, name string) error {
				return c.SetAPAdminState(ctx, name, v.on)
			})
		},
	}
}

func adminRadioLeaf(v adminVerb) *cli.Command {
	return &cli.Command{
		Name:      "radio",
		Usage:     v.title + " one radio of one access point",
		UsageText: synopsis(config.FlagAPName, config.FlagSlot),
		Description: "--ap-name is the ap_name column of wnc show ap, and --slot is the Slot column\n" +
			"of wnc show overview. The band the RPC needs is read from the controller and\n" +
			"follows the radio type, so a dual-band radio takes one number either way.\n" +
			"Pass --dry-run to name the target and change nothing.",
		Flags: []cli.Flag{apNameFlag(), yesFlag(), slotFlag()},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runRadioAdmin(ctx, cmd, v)
		},
	}
}

// slotFlag names the radio. HideDefault stops the help printing "(default: 0)" for a flag slotArg
// requires, since zero is a slot the controller reports and that text read as an omittable
// default; Required would suppress it too, but would replace slotArg's message with urfave's own.
func slotFlag() cli.Flag {
	return &cli.IntFlag{
		Name:        config.FlagSlot,
		Usage:       "radio slot, as shown in the Slot column of wnc show overview",
		HideDefault: true,
	}
}

// slotArg reads the slot. IsSet and not a zero test: zero is a slot the controller reports, so a
// missing flag and a named zero must not read alike. The bound is the wnc layer's, derived
// from the same must table a pair is checked against once the band is known.
func slotArg(cmd *cli.Command) (int, error) {
	if !cmd.IsSet(config.FlagSlot) {
		return 0, fmt.Errorf("%w: %s requires --%s: the radio's slot, 0 to %d",
			ErrUsage, cmd.Name, config.FlagSlot, wnc.MaxRadioSlot())
	}

	slot := cmd.Int(config.FlagSlot)
	if slot < 0 || slot > wnc.MaxRadioSlot() {
		return 0, fmt.Errorf("%w: --%s %d: accepted values are 0 to %d",
			ErrUsage, config.FlagSlot, slot, wnc.MaxRadioSlot())
	}

	return slot, nil
}

// runRadioAdmin sets one radio's admin state. It cannot reuse runAPAction: this is the one
// write whose RPC has no ap-name arm, so it needs the address the resolve returns, and the
// radio read must land before the dry-run line so it can name the band and before the prompt
// so the band is checked rather than consented to.
func runRadioAdmin(ctx context.Context, cmd *cli.Command, v adminVerb) error {
	name, err := requireAPName(cmd)
	if err != nil {
		return err
	}

	slot, err := slotArg(cmd)
	if err != nil {
		return err
	}

	target, client, err := actionPrologue(ctx, cmd)
	if err != nil {
		return err
	}

	mac, err := resolveAP(ctx, client, target, name)
	if err != nil {
		return err
	}

	radio, err := radioForAdmin(ctx, client, target, name, mac, slot)
	if err != nil {
		return err
	}

	label := fmt.Sprintf("slot %d (%s GHz) of %s on %s", slot, radio.BandLabel, name, target.Name)

	return confirmAndAct(ctx, cmd, writeWording{
		question: v.title + " %s? [y/N]: ",
		dryRun:   "%s: would " + v.name + "\n",
		sent:     "%s: " + v.name + " sent\n",
		failed:   v.gerund + " %s",
	}, func(ctx context.Context) error {
		return client.SetRadioAdminState(ctx, mac, slot, radio.Type, v.on)
	}, label)
}

// radioForAdmin reads the one radio the write will name and refuses every state the RPC cannot
// express. The band number is the radio type's and the band the prompt names is the served band's,
// so four things are refused here: a slot that is not a radio, a type the RPC takes no number for,
// a served band with no label to prompt with, and a pair the must clause forbids.
func radioForAdmin(
	ctx context.Context, client *wnc.Client, target config.Target, name, mac string, slot int,
) (*wnc.RadioAdmin, error) {
	// radio-oper-data is keyed on wtp-mac, so a map row carrying no address would key the read on
	// the empty string and the 404 would be reported as a slot the access point does not hold.
	if mac == "" {
		return nil, fmt.Errorf("%s reports no radio address for %s", target.Name, name)
	}

	radio, err := client.RadioBySlot(ctx, mac, slot)
	if err != nil {
		if cause, _ := wnc.Classify(err); cause == wnc.CauseNotFound {
			return nil, fmt.Errorf("%s holds no radio in slot %d of %s", target.Name, slot, name)
		}

		return nil, fmt.Errorf("reading radio-oper-data from %s: %s", target.Name, wnc.Message(err))
	}

	if radio == nil {
		return nil, fmt.Errorf("%s holds no radio in slot %d of %s", target.Name, slot, name)
	}

	if !radio.IsRadio() {
		return nil, fmt.Errorf("slot %d of %s on %s is a remote-LAN port, not a radio",
			slot, name, target.Name)
	}

	if radio.BandWire == 0 {
		if radio.Type == "" {
			return nil, fmt.Errorf("%s reports no radio type for slot %d of %s", target.Name, slot, name)
		}

		return nil, fmt.Errorf(
			"%s reports slot %d of %s as radio type %s, which this RPC has no band number for",
			target.Name, slot, name, radio.Type)
	}

	// The label is the prompt's and not the wire's: a radio whose served band the controller does
	// not name cannot be checked against the Band column of wnc show overview before the write.
	if radio.BandLabel == "" {
		if radio.Band == "" {
			return nil, fmt.Errorf("%s reports no band for slot %d of %s", target.Name, slot, name)
		}

		return nil, fmt.Errorf("%s reports slot %d of %s on an unknown band (%s)",
			target.Name, slot, name, radio.Band)
	}

	if !radio.SlotAllowed() {
		return nil, fmt.Errorf(
			"%s reports slot %d of %s as radio type %s, and accepts that radio on slot %s only",
			target.Name, slot, name, radio.Type, joinInts(radio.AllowedSlots()))
	}

	return radio, nil
}

func joinInts(v []int) string {
	out := make([]string, 0, len(v))
	for _, n := range v {
		out = append(out, strconv.Itoa(n))
	}

	return strings.Join(out, " or ")
}
