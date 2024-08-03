package main

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// List of known insecure images or patterns
var insecureImages = []string{
	"ubuntu:16.04", // Sample insecure image
	"alpine:3.6",   // SAMple insecure image
}

var versionCheckList = map[string]string{
	"kube-apiserver":          "v1.24.0", // Sample version to check against
	"kube-scheduler":          "v1.24.0", // Sample version to check against
	"kube-controller-manager": "v1.24.0", // Sample version to check against
}

func scanRBACRoles(clientset *kubernetes.Clientset, issues *[]string) {
	roles, err := clientset.RbacV1().Roles("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		*issues = append(*issues, fmt.Sprintf("Failed to list roles: %v", err))
		return
	}

	for _, role := range roles.Items {
		checkRole(&role, issues)
	}
}

func checkRole(role *rbacv1.Role, issues *[]string) {
	for _, rule := range role.Rules {
		if len(rule.Verbs) == 0 {
			*issues = append(*issues, fmt.Sprintf("Warning: ROle %s has empty verbs.", role.Name))
		}
	}
}

func scanPods(clientset *kubernetes.Clientset, issues *[]string) {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		*issues = append(*issues, fmt.Sprintf("Failed to list pods: %v", err))
		return
	}

	for _, pod := range pods.Items {
		checkPod(&pod, issues)
	}
}

func checkPod(pod *corev1.Pod, issues *[]string) {
	for _, container := range pod.Spec.Containers {
		if isInsecureImage(container.Image) {
			*issues = append(*issues, fmt.Sprintf("Warning: Pod%s has insecure image %s.", pod.Name, container.Image))
		}
	}
}

func isInsecureImage(image string) bool {
	for _, insecureImage := range insecureImages {
		if strings.Contains(image, insecureImage) {
			return true
		}
	}
	return false
}

func checkSoftwareVersions(issues *[]string) {
	// This is a mock function. Replace with actual logic to retrieve versions.
	currentVersions := map[string]string{
		"kube-apiserver": "v1.22.0", // Sample current version
		"kube-scheduler": "v1.23.0", // Sample current version
	}

	for component, currentVersion := range currentVersions {
		if minVersion, ok := versionCheckList[component]; ok {
			if compareVersions(currentVersion, minVersion) < 0 {
				*issues = append(*issues, fmt.Sprintf("Warning: %s is outdated. Current version: %s, Minimum required version: %s.", component, currentVersion, minVersion))
			}
		}
	}
}

// Compare versions, returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	// A Simple version comparison logic (assuming semantic versioning)
	return strings.Compare(v1, v2)
}
