package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"log/slog"
)

func main() {
	programLevel := new(slog.LevelVar)
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel}))
	slog.SetDefault(logger)

	integrationID := os.Getenv("INTEGRATION_ID")
	logger.Info(integrationID)
	outputDir := os.Getenv("OUTPUT_DIR")

	// get input files
	sessionToken := os.Getenv("SESSION_TOKEN")
	apiHost := os.Getenv("PENNSIEVE_API_HOST")
	apiHost2 := os.Getenv("PENNSIEVE_API_HOST2")
	integrationResponse, err := getIntegration(apiHost2, integrationID, sessionToken)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(integrationResponse))
	var integration Integration
	if err := json.Unmarshal(integrationResponse, &integration); err != nil {
		logger.ErrorContext(context.Background(), err.Error())
	}
	fmt.Println(integration)

	manifest, err := getPresignedUrls(apiHost, getPackageIds(integration.PackageIDs), sessionToken)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(manifest))
	var payload Manifest
	if err := json.Unmarshal(manifest, &payload); err != nil {
		logger.ErrorContext(context.Background(), err.Error())
	}

	// copy files into input directory
	fmt.Println(payload.Data)
	for _, d := range payload.Data {

		// Prepare the download path with directory structure
		fullFilePath, err := getDownloadPath(outputDir, d.Path, d.FileName, logger)
		if err != nil {
			continue // Skip this file if directory creation failed
		}

		cmd := exec.Command("wget", "-v", "-O", fullFilePath, d.Url) // Note: we don't set cmd.Dir because we're using absolute paths
		var out strings.Builder
		var stderr strings.Builder
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()

		if err != nil {
			logger.Error("Download failed",
				slog.String("file", d.FileName),
				slog.String("error", stderr.String()))
			continue
		}

		// Print stdout content
		stdoutContent := out.String()
		fmt.Println("Stdout output:")
		fmt.Println(stdoutContent)

		// Print or log stderr content
		stderrContent := stderr.String()
		fmt.Println("Stderr output (verbose output):")
		fmt.Println(stderrContent)

		fmt.Printf("✓ Successfully downloaded: %s\n\n", fullFilePath)
	}

}

type Packages struct {
	NodeIds []string `json:"nodeIds"`
}

type Manifest struct {
	Data []ManifestData `json:"data"`
}

type ManifestData struct {
	NodeId   string   `json:"nodeId"`
	FileName string   `json:"fileName"`
	Path     []string `json:"path"`
	Url      string   `json:"url"`
}

type Integration struct {
	Uuid          string      `json:"uuid"`
	ApplicationID int64       `json:"applicationId"`
	DatasetNodeID string      `json:"datasetId"`
	PackageIDs    []string    `json:"packageIds"`
	Params        interface{} `json:"params"`
}

func getPresignedUrls(apiHost string, packages Packages, sessionToken string) ([]byte, error) {
	url := fmt.Sprintf("%s/packages/download-manifest?api_key=%s", apiHost, sessionToken)
	b, err := json.Marshal(packages)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(b))

	payload := strings.NewReader(string(b))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "*/*")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return body, nil
}

func getPackageIds(packageIds []string) Packages {
	return Packages{
		NodeIds: packageIds,
	}
}

func getIntegration(apiHost string, integrationId string, sessionToken string) ([]byte, error) {
	url := fmt.Sprintf("%s/integrations/%s", apiHost, integrationId)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sessionToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return body, nil
}

// getDownloadPath constructs the full file path with directory structure and creates directories if needed
func getDownloadPath(outputDir string, path []string, fileName string, logger *slog.Logger) (string, error) {
	
	// Construct the target directory path based on the folder structure
	var targetDir string
	if len(path) > 0 {
		targetDir = outputDir + "/" + strings.Join(path, "/") 
	} else {
		targetDir = outputDir 
	}
	
	// Create the directory structure if it doesn't exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		logger.Error("Failed to create directory structure",
			slog.String("directory", targetDir),
			slog.String("error", err.Error()))
		return "", err 
	}
	
	// Full path where the file will be saved
	fullFilePath := targetDir + "/" + fileName 
	fmt.Printf("Downloading to: %s\n", fullFilePath)
	
	return fullFilePath, nil
}
