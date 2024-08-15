package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Mock exec.Command for testing purposes
var mockExecCommand = func(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// Helper process for mocking exec.Command in tests
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args[3:]

	if strings.HasPrefix(args[0], "lpstat") {
		if strings.Contains(args[1], "-d") {
			// Simulate output of `lpstat -d`
			os.Stdout.WriteString("printer Test_Printer_1 is idle. enabled since Thu 01 Jan 1970 00:00:00 UTC")
		}
	} else if strings.HasPrefix(args[0], "ssh") {
		// Simulate output of SSH command
		os.Stdout.WriteString("printer Test_Remote_Printer_1 is idle. enabled since Thu 01 Jan 1970 00:00:00 UTC")
	} else if strings.HasPrefix(args[0], "scp") {
		// Simulate successful file transfer
		os.Exit(0)
	} else if strings.HasPrefix(args[0], "lp") {
		// Simulate successful print command
		os.Exit(0)
	}

	os.Exit(0)
}

// Test FindPrinterName function
func TestFindPrinterName(t *testing.T) {
	execCommand = mockExecCommand
	printerName, err := FindPrinterName()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if printerName != "Test_Printer_1" {
		t.Fatalf("Expected printer name 'Test_Printer_1', got %v", printerName)
	}
}

// Test FindPrinterOverSSH function
func TestFindPrinterOverSSH(t *testing.T) {
	execCommand = mockExecCommand
	printerName, err := FindPrinterOverSSH("test_host", "test_user", "test_pass")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if printerName != "Test_Remote_Printer_1" {
		t.Fatalf("Expected printer name 'Test_Remote_Printer_1', got %v", printerName)
	}
}

// Test PrinterLocally function
func TestPrinterLocally(t *testing.T) {
	execCommand = mockExecCommand
	p := &Printer{Name: "Test_Printer_1"}

	err := p.PrinterLocally("/tmp/test_file.txt")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

// Test PrintRemotely function
func TestPrintRemotely(t *testing.T) {
	execCommand = mockExecCommand
	p := &Printer{Name: "Test_Remote_Printer_1", Remote: true, Host: "test_host", Username: "test_user"}

	err := p.PrintRemotely("/tmp/test_file.txt")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

// Test ParseOutput function
func TestParseOutput(t *testing.T) {
	output := "printer Test_Printer_1 is idle. enabled since Thu 01 Jan 1970 00:00:00 UTC"
	expected := "Test_Printer_1"
	result := ParseOutput(output)

	if result != expected {
		t.Fatalf("Expected %v, got %v", expected, result)
	}
}
