package command

import (
	"errors"
	"fmt"
	"text/template"

	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/tid/cli/display"
	"github.com/SeerUK/tid/pkg/types"
	"github.com/SeerUK/tid/pkg/util"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// StatusCommand creates a command to view the status of the current timer.
func StatusCommand(factory util.Factory, config types.Config) *console.Command {
	var format string
	var hash string

	configure := func(def *console.Definition) {
		def.AddArgument(console.ArgumentDefinition{
			Value: parameters.NewStringValue(&hash),
			Spec:  "[HASH]",
			Desc:  "A short or long hash for an entry.",
		})

		def.AddOption(console.OptionDefinition{
			Value: parameters.NewStringValue(&format),
			Spec:  "-f, --format=FORMAT",
			Desc:  "Output formatting string. Uses Go templates.",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		sysGateway := factory.BuildSysGateway()
		trGateway := factory.BuildTrackingGateway()

		hasFormat := input.HasOption([]string{"f", "format"})

		status, err := sysGateway.FindOrCreateStatus()
		if err != nil {
			return err
		}

		if hash == "" {
			hash = status.Entry
		}

		if hash == "" {
			return errors.New("status: No timer to check the status of")
		}

		entry, err := trGateway.FindEntry(hash)
		if err != nil && err != state.ErrStoreNilResult {
			return err
		}

		if err == state.ErrStoreNilResult {
			return fmt.Errorf("status: No entry with hash '%s'", hash)
		}

		if hasFormat {
			tmpl := template.Must(template.New("entry-list").Parse(format))
			tmpl.Execute(output.Writer, entry)

			output.Println()

			return nil
		}

		display.WriteEntriesTable([]types.Entry{entry}, output.Writer, config)

		return nil
	}

	return &console.Command{
		Name:        "status",
		Alias:       "st",
		Description: "View the current status.",
		Configure:   configure,
		Execute:     execute,
	}
}
