package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	format := flag.String("format", "json", "Output format: json or yaml")
	outputDir := flag.String("output-dir", ".", "Directory to save the report")
	flag.Parse()

	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Failed to build config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}

	// Start scanning
	issues := scanCluster(clientset)

	report := SecurityReport{
		Timestamp: time.Now(),
		Issues:    issues,
	}

	reportContent, err := generateReport(report, *format)
	if err != nil {
		log.Fatalf("Failed to generate report: %v", err)
	}

	if err := writeReportToFile(reportContent, *outputDir, *format); err != nil {
		log.Fatalf("Failed to write report to file: %v", err)
	}

	fmt.Println("Security report generated successfully.")
}

func scanCluster(clientset *kubernetes.Clientset) []string {
	var issues []string

	// Scan RBAC Roles
	scanRBACRoles(clientset, &issues)

	// Scan Pods
	scanPods(clientset, &issues)

	// Check Software Versions
	checkSoftwareVersions(&issues)

	return issues
}
