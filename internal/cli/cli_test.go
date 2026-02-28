package cli

import (
	"testing"

	"github.com/roniel/todo-app/internal/config"
)

func TestRun_NoArgs(t *testing.T) {
	err := Run(nil, nil, nil, nil, config.Config{})
	if err == nil {
		t.Fatal("expected error for no args, got nil")
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	err := Run([]string{"foobar"}, nil, nil, nil, config.Config{})
	if err == nil {
		t.Fatal("expected error for unknown command, got nil")
	}
}

func TestCmdAdd_NoTitle(t *testing.T) {
	err := Run([]string{"add"}, nil, nil, nil, config.Config{})
	if err == nil {
		t.Fatal("expected error for add with no title, got nil")
	}
}

func TestCmdAdd_InvalidPriority(t *testing.T) {
	err := Run([]string{"add", "--priority", "extreme", "my task"}, nil, nil, nil, config.Config{})
	if err == nil {
		t.Fatal("expected error for invalid priority, got nil")
	}
}

func TestCmdAdd_InvalidDueDate(t *testing.T) {
	err := Run([]string{"add", "--due", "not-a-date", "my task"}, nil, nil, nil, config.Config{})
	if err == nil {
		t.Fatal("expected error for invalid due date, got nil")
	}
}

func TestCmdDone_NoID(t *testing.T) {
	err := Run([]string{"done"}, nil, nil, nil, config.Config{})
	if err == nil {
		t.Fatal("expected error for done with no ID, got nil")
	}
}

func TestCmdDone_InvalidID(t *testing.T) {
	err := Run([]string{"done", "abc"}, nil, nil, nil, config.Config{})
	if err == nil {
		t.Fatal("expected error for invalid ID, got nil")
	}
}

func TestCmdList_InvalidPriority(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "Task"}, ts, js)
	err := run(t, []string{"list", "--priority", "extreme"}, ts, js)
	if err == nil {
		t.Fatal("expected error for invalid priority, got nil")
	}
}

func TestCmdJournalAdd_NoText(t *testing.T) {
	err := Run([]string{"journal"}, nil, nil, nil, config.Config{})
	if err == nil {
		t.Fatal("expected error for journal with no text, got nil")
	}
}

func TestCmdExport_InvalidFormat(t *testing.T) {
	ts, js := newTestStores(t)
	err := run(t, []string{"export", "--format", "csv"}, ts, js)
	if err == nil {
		t.Fatal("expected error for invalid export format, got nil")
	}
}

func TestFilterTasks_All(t *testing.T) {
	tasks := filterTasks(nil, "all")
	if tasks != nil {
		t.Errorf("expected nil for nil input, got %v", tasks)
	}
}

func TestFilterTasks_Pending(t *testing.T) {
	tasks := filterTasks(nil, "pending")
	if tasks != nil {
		t.Errorf("expected nil for nil input, got %v", tasks)
	}
}

func TestFilterTasks_Done(t *testing.T) {
	tasks := filterTasks(nil, "done")
	if tasks != nil {
		t.Errorf("expected nil for nil input, got %v", tasks)
	}
}

func TestFilterTasks_Active(t *testing.T) {
	tasks := filterTasks(nil, "active")
	if tasks != nil {
		t.Errorf("expected nil for nil input, got %v", tasks)
	}
}
