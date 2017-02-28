package cli

import (
	"text/template"

	"fmt"

	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
	"github.com/olekukonko/tablewriter"
)

// statusOutput represents the formattable source of the status command output.
type statusOutput struct {
	Entry  tracking.Entry
	Status tracking.Status
}

// StatusCommand creates a command to view the status of the current timer.
func StatusCommand(gateway tracking.Gateway) console.Command {
	var format string
	var hash string

	configure := func(def *console.Definition) {
		def.AddOption(
			parameters.NewStringValue(&format),
			"-f, --format=FORMAT",
			"Format string, uses Go templates.",
		)

		def.AddArgument(
			parameters.NewStringValue(&hash),
			"[HASH]",
			"A short or long hash for an entry.",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		status, err := gateway.FindOrCreateStatus()
		if err != nil {
			return err
		}

		if status.Entry == "" {
			output.Println("status: No timer to check the status of")
			return nil
		}

		if hash == "" {
			hash = status.Entry
		}

		entry, err := gateway.FindEntry(hash)
		if err != nil && err != state.ErrNilResult {
			return err
		}

		if err == state.ErrNilResult {
			output.Printf("status: No entry with hash '%s'\n", hash)
			return nil
		}

		dateFormat := "3:04:05PM (2006-01-02)"

		if entry.IsRunning {
			// If we're viewing the status of the currently active entry, we should get make sure
			// that it's duration is up-to-date.
			entry.UpdateDuration()
		}

		if format != "" {
			out := statusOutput{}
			out.Entry = entry
			out.Status = status

			tmpl := template.Must(template.New("status").Parse(format))
			tmpl.Execute(output.Writer, out)

			// Always end with a new line...
			output.Println()
		} else {
			table := tablewriter.NewWriter(output.Writer)
			table.SetHeader([]string{
				"Date",
				"Hash",
				"Created",
				"Updated",
				"Note",
				"Duration",
				"Running",
			})
			table.Append([]string{
				entry.Timesheet,
				entry.ShortHash(),
				entry.Created.Format(dateFormat),
				entry.Updated.Format(dateFormat),
				entry.Note,
				entry.Duration.String(),
				fmt.Sprintf("%t", entry.IsRunning),
			})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.Render()
		}

		return nil
	}

	return console.Command{
		Name:        "status",
		Description: "View the current status.",
		Configure:   configure,
		Execute:     execute,
	}
}
