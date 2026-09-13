package main

import (
	"bytes"
	"testing"
)

func TestExecuteHelpAndVersion(t *testing.T) {
	for _, args := range [][]string{nil, {"--help"}, {"help"}} {
		var out bytes.Buffer
		if code := execute(args, &out); code != 0 || out.Len() == 0 {
			t.Errorf("execute(%v) = code %d, output %q", args, code, out.String())
		}
	}
	for _, arg := range []string{"--version", "-v", "-V"} {
		var out bytes.Buffer
		if code := execute([]string{arg}, &out); code != 0 || out.String() != version+"\n" {
			t.Errorf("execute(%q) = code %d, output %q", arg, code, out.String())
		}
	}
}

func TestPlainText(t *testing.T) {
	got := plainText("<div>Hello&nbsp;<b>world</b> &#128522;</div>")
	if got != "Hello world 😊" {
		t.Fatalf("plainText = %q", got)
	}
}

func TestListQuery(t *testing.T) {
	got := buildWIQL(listOptions{state: "Active, New", assignedTo: "Cory", limit: 20})
	want := "SELECT [System.Id], [System.Title], [System.State], [System.AssignedTo] FROM WorkItems WHERE ([System.State] = 'Active' OR [System.State] = 'New') AND [System.AssignedTo] CONTAINS 'Cory' ORDER BY [System.ChangedDate] DESC"
	if got != want {
		t.Fatalf("query = %q, want %q", got, want)
	}
}
