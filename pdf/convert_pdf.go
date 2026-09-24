package pdf

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ConvertToPDF(docxPath, outputDir string) error {
	candidates := []string{
		"libreoffice", // Linux / PATH / some macOS installations
		"soffice",     // Linux / macOS / PATH
		"/Applications/LibreOffice.app/Contents/MacOS/soffice", // macOS
		`C:\Program Files\LibreOffice\program\soffice.exe`,     // Windows
		`C:\Program Files (x86)\LibreOffice\program\soffice.exe`,
	}

	var lastErr error

	for _, executable := range candidates {
		if _, err := exec.LookPath(executable); err != nil {
			// Absolute paths aren't always handled by LookPath the way
			// we want, so check them directly.
			if filepath.IsAbs(executable) {
				if _, statErr := os.Stat(executable); statErr != nil {
					continue
				}
			} else {
				continue
			}
		}

		cmd := exec.Command(
			executable,
			"--headless",
			"--convert-to",
			"pdf",
			"--outdir",
			outputDir,
			docxPath,
		)

		output, err := cmd.CombinedOutput()
		if err == nil {
			return nil
		}

		lastErr = fmt.Errorf(
			"%s failed: %w: %s",
			executable,
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return fmt.Errorf(
		"LibreOffice not found or PDF conversion failed: %w",
		lastErr,
	)
}