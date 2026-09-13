package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/cory-evans/devops-axi/internal/az"
	"github.com/spf13/cobra"
)

type listOptions struct {
	wiql, state, assignedTo, fields string
	me                              bool
	limit                           int
}

var runAz = az.Run

func newListCommand(out io.Writer) *cobra.Command {
	o := &listOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List work items.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if o.wiql != "" && (o.state != "" || o.assignedTo != "" || o.me) {
				return fmt.Errorf("--wiql cannot be combined with --state, --assigned-to, or --me")
			}
			if o.assignedTo != "" && o.me {
				return fmt.Errorf("--assigned-to cannot be combined with --me")
			}
			query := o.wiql
			if query == "" {
				query = buildWIQL(*o)
			}
			data, err := runAz(context.Background(), "boards", "query", "--wiql", query, "--output", "json", "--only-show-errors")
			if err != nil {
				return fmt.Errorf("unable to list work items; verify Azure DevOps configuration and authentication")
			}
			if printItems(cmd.OutOrStdout(), data, *o) != 0 {
				return fmt.Errorf("Azure DevOps returned invalid work-item data")
			}
			return nil
		},
	}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.Flags().StringVar(&o.wiql, "wiql", "", "Run a raw flat WIQL query.")
	cmd.Flags().StringVar(&o.state, "state", "", "Comma-separated states, matched with OR.")
	cmd.Flags().StringVar(&o.assignedTo, "assigned-to", "", "Partial assignee display-name match.")
	cmd.Flags().BoolVar(&o.me, "me", false, "Show work items assigned to the current user.")
	cmd.Flags().IntVar(&o.limit, "limit", 20, "Maximum rows to return.")
	cmd.Flags().StringVar(&o.fields, "fields", "", "Comma-separated additional fields.")
	cmd.Example = "  devops-axi work-item list\n  devops-axi work-item list --state Active --limit 10\n  devops-axi work-item list --me"
	return cmd
}

func buildWIQL(o listOptions) string {
	where := []string{"[System.State] <> 'Closed'"}
	if o.state != "" {
		var matches []string
		for _, state := range strings.Split(o.state, ",") {
			if state = strings.TrimSpace(state); state != "" {
				matches = append(matches, "[System.State] = '"+wiqlQuote(state)+"'")
			}
		}
		if len(matches) > 0 {
			where = []string{"(" + strings.Join(matches, " OR ") + ")"}
		}
	}
	if o.assignedTo != "" {
		where = append(where, "[System.AssignedTo] CONTAINS '"+wiqlQuote(o.assignedTo)+"'")
	}
	if o.me {
		where = append(where, "[System.AssignedTo] = @Me")
	}
	return fmt.Sprintf("SELECT [System.Id], [System.Title], [System.State], [System.AssignedTo] FROM WorkItems WHERE %s ORDER BY [System.ChangedDate] DESC", strings.Join(where, " AND "))
}

func wiqlQuote(s string) string { return strings.ReplaceAll(s, "'", "''") }

type workItem struct {
	ID     any            `json:"id"`
	Fields map[string]any `json:"fields"`
}

func printItems(w io.Writer, data []byte, o listOptions) int {
	var items []workItem
	if err := json.Unmarshal(data, &items); err != nil {
		return 1
	}
	limited := len(items) > o.limit
	if limited {
		items = items[:o.limit]
	}
	fields := []string{"System.Id", "System.Title", "System.State", "System.AssignedTo"}
	for _, f := range strings.Split(o.fields, ",") {
		if f = strings.TrimSpace(f); f != "" && !contains(fields, f) {
			fields = append(fields, f)
		}
	}
	fmt.Fprintf(w, "count: %d\nlimited: %t\nitems[%d]{%s}:\n", len(items), limited, len(items), strings.Join(toonNames(fields), ","))
	for _, item := range items {
		values := make([]string, len(fields))
		for i, field := range fields {
			values[i] = toonValue(fieldValue(item, field))
		}
		fmt.Fprintf(w, "  %s\n", strings.Join(values, ","))
	}
	return 0
}

func fieldValue(item workItem, field string) any {
	if field == "System.Id" {
		return item.ID
	}
	v := item.Fields[field]
	if m, ok := v.(map[string]any); ok {
		return m["displayName"]
	}
	return v
}
func toonNames(fields []string) []string {
	out := make([]string, len(fields))
	for i, f := range fields {
		switch f {
		case "System.Id":
			out[i] = "id"
		case "System.Title":
			out[i] = "title"
		case "System.State":
			out[i] = "state"
		case "System.AssignedTo":
			out[i] = "assignedTo"
		default:
			out[i] = f
		}
	}
	return out
}
func toonValue(v any) string {
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return fmt.Sprint(v)
	}
	if strings.ContainsAny(s, ",\n\"") {
		return strconv.Quote(s)
	}
	return s
}
func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}
