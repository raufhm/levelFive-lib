/*
Provide a utility
to managing and printing documents to both local and remote printers using the Common Unix Printing System (CUPS).

Key Features:
- Automatically identifies the default printer using `lpstat`.
- Supports printing to both local and remote printers.
- Utilizes SCP for securely transferring print jobs to remote machines.
- Provides an abstracted command execution mechanism for easier testing and flexibility.

References:
- CUPS (Common Unix Printing System) options:
https://opensource.apple.com/source/cups/cups-450/cups/doc/help/options.html

Usage:
- Instantiate the Printer struct with the necessary details (e.g., printer name, host details for remote printers).
- Call the `Do()` method on the Printer instance to execute the print job.
*/
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Printer represents a printer that can be used to print documents.
type Printer struct {
	Name     string
	Remote   bool
	Host     string
	Username string
	Password string
	RawFile  string
}

// CommandExecutor is a type that defines a function to execute commands.
type CommandExecutor func(name string, arg ...string) *exec.Cmd

// execCommand is a variable that holds the command execution function.
var execCommand CommandExecutor = exec.Command

func ParseOutput(output string) string {
	var printerName string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "printer") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				printerName = parts[1]
				break
			}

		}
	}
	return printerName
}

// FindPrinterName uses lpstat to find printers configured with CUPS.
func FindPrinterName() (string, error) {

	cmd := execCommand("lpstat", "-d")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute lpstat command: %w", err)
	}

	return ParseOutput(string(output)), nil
}

func FindPrinterOverSSH(remoteHost, username, password string) (string, error) {
	// Example SSH command to execute lpstat on the remote machine
	cmd := execCommand("ssh", fmt.Sprintf("%s@%s", username, remoteHost), "lpstat", "-d")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute SSH command: %w", err)
	}

	return ParseOutput(string(output)), nil
}

func (p *Printer) Do() error {
	tmpFile, err := os.CreateTemp("", "receipt-*.txt")
	if err != nil {
		return err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			return
		}
	}(tmpFile.Name())

	_, err = tmpFile.WriteString(p.RawFile)
	if err != nil {
		return err
	}

	if err = tmpFile.Close(); err != nil {
		return err
	}

	if p.Remote {
		return p.PrintRemotely(tmpFile.Name())
	}

	return p.PrinterLocally(tmpFile.Name())
}

func (p *Printer) PrinterLocally(filePath string) error {
	cmd := execCommand("lp", "-d", p.Name, filePath)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (p *Printer) PrintRemotely(filePath string) error {
	cmd := execCommand("scp", filePath, fmt.Sprintf("%s@%s:/tmp/", p.Username, p.Host))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to transfer file to remote host: %w", err)
	}

	remoteCmd := fmt.Sprintf("lp -d %s /tmp/%s", p.Name, filePath)
	cmd = execCommand("ssh", fmt.Sprintf("%s@%s", p.Username, p.Host), remoteCmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to print remotely: %w", err)
	}

	return nil
}
