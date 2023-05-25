// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 Authors of KubeArmor

// Package version checks the current CLI version and if there's a need to update it
package version

import (
	"context"
	"fmt"
	"runtime"

	"github.com/accuknox/accuknox-cli/k8s"
	"github.com/accuknox/accuknox-cli/selfupdate"
	"github.com/fatih/color"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PrintVersion handler for accuknox-cli version
func PrintVersion(c *k8s.Client) error {
	fmt.Printf("accuknox-cli version %s %s/%s BuildDate=%s\n", selfupdate.GitSummary, runtime.GOOS, runtime.GOARCH, selfupdate.BuildDate)
	latest, latestVer := selfupdate.IsLatest(selfupdate.GitSummary)
	if !latest {
		color.HiMagenta("update available version " + latestVer)
		color.HiMagenta("use [accuknox-cli selfupdate] to update to latest")
	}
	kubearmorVersion, err := getKubeArmorVersion(c)
	if err != nil {
		return nil
	}
	if kubearmorVersion == "" {
		fmt.Printf("kubearmor not running\n")
		return nil
	}
	fmt.Printf("kubearmor image (running) version %s\n", kubearmorVersion)
	return nil
}

func getKubeArmorVersion(c *k8s.Client) (string, error) {
	pods, err := c.K8sClientset.CoreV1().Pods("kube-system").List(context.Background(), metav1.ListOptions{LabelSelector: "kubearmor-app=kubearmor"})
	if err != nil {
		return "", err
	}
	if len(pods.Items) > 0 {
		image := pods.Items[0].Spec.Containers[0].Image
		return image, nil
	}
	return "", nil
}
