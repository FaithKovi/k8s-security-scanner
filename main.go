package main

import (
	"flag"
	"log"
	"time"

	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	// Define a flag for output format
	format := flag.String("format", "json", "Output format: json or yaml")
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

	printReport(reportContent)
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
