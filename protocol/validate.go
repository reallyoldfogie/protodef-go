package protocol

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateJSONSchema validates a JSON document against a JSON schema file.
// If ajv-cli is available via npx, it will be used to support advanced regex patterns.
// Otherwise, falls back to gojsonschema (which has limited regex support).
func ValidateJSONSchema(schemaPath string, jsonData []byte) (bool, error) {
	// Try ajv first if available (supports negative lookahead regex)
	if valid, err := validateWithAJV(schemaPath, jsonData); err == nil {
		return valid, nil
	}

	// Fallback to gojsonschema
	absPath, err := filepath.Abs(schemaPath)
	if err != nil {
		return false, err
	}
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + absPath)
	documentLoader := gojsonschema.NewBytesLoader(jsonData)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}

// validateWithAJV uses ajv-cli via npx for validation if available
func validateWithAJV(schemaPath string, jsonData []byte) (bool, error) {
	// Check if npx is available
	if _, err := exec.LookPath("npx"); err != nil {
		return false, fmt.Errorf("npx not available")
	}

	// Create temp file for JSON data
	tmpFile, err := os.CreateTemp("", "validate-*.json")
	if err != nil {
		return false, err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(jsonData); err != nil {
		tmpFile.Close()
		return false, err
	}
	tmpFile.Close()

	// Get absolute paths
	absSchemaPath, err := filepath.Abs(schemaPath)
	if err != nil {
		return false, err
	}

	// Run ajv validate with strict=false to allow older schema patterns
	cmd := exec.Command("npx", "-y", "ajv-cli", "validate", "--strict=false", "-s", absSchemaPath, "-d", tmpFile.Name())
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		// Check if it's a validation failure or actual error
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				// Validation failed (invalid data)
				return false, nil
			}
		}
		// Some other error occurred
		if strings.Contains(stderr.String(), "not found") {
			return false, fmt.Errorf("ajv-cli not available")
		}
		return false, fmt.Errorf("ajv error: %v - %s", err, stderr.String())
	}

	return true, nil
}
