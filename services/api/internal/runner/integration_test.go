//go:build integration

package runner

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/testcontainers/testcontainers-go/modules/k3s"
)

var testK8sClient kubernetes.Interface

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Start K3s container
	fmt.Println("Starting K3s testcontainer...")
	container, err := k3s.Run(ctx, "rancher/k3s:v1.31.2-k3s1")
	if err != nil {
		fmt.Fprintf(os.Stderr, "K3s start failed: %v\n", err)
		os.Exit(1)
	}
	defer container.Terminate(ctx)

	// Get kubeconfig
	kubeconfig, err := container.GetKubeConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Create K8s client from kubeconfig
	restCfg, err := clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse kubeconfig: %v\n", err)
		os.Exit(1)
	}

	testK8sClient, err = kubernetes.NewForConfig(restCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create K8s client: %v\n", err)
		os.Exit(1)
	}

	// Wait for K3s to be ready
	fmt.Println("Waiting for K3s to be ready...")
	for i := 0; i < 30; i++ {
		_, err := testK8sClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	fmt.Println("K3s ready")

	os.Exit(m.Run())
}

func TestInteg_K8sRunner_CreateJob(t *testing.T) {
	ctx := context.Background()

	r := &KubernetesRunner{
		client: testK8sClient,
		config: Config{
			AgentImage:    "alpine:3.21", // lightweight image for testing
			Namespace:     "default",
			CPURequest:    "50m",
			CPULimit:      "100m",
			MemoryRequest: "32Mi",
			MemoryLimit:   "64Mi",
		},
	}

	err := r.Run(ctx, RunOptions{
		ProjectID: "integ-proj",
		RunID:     "integ-run-12345678",
		MaxSteps:  5,
	})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Verify Job exists in K3s
	jobs, err := testK8sClient.BatchV1().Jobs("default").List(ctx, metav1.ListOptions{
		LabelSelector: "run-id=integ-run-12345678",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs.Items) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs.Items))
	}

	job := jobs.Items[0]
	t.Logf("Job created: %s", job.Name)

	// Verify labels
	if job.Labels["app"] != "decisionbox-agent" {
		t.Errorf("label app = %q", job.Labels["app"])
	}
	if job.Labels["project-id"] != "integ-proj" {
		t.Errorf("label project-id = %q", job.Labels["project-id"])
	}

	// Verify container image
	if job.Spec.Template.Spec.Containers[0].Image != "alpine:3.21" {
		t.Errorf("image = %q", job.Spec.Template.Spec.Containers[0].Image)
	}
}

func TestInteg_K8sRunner_CancelJob(t *testing.T) {
	ctx := context.Background()

	r := &KubernetesRunner{
		client: testK8sClient,
		config: Config{
			AgentImage:    "alpine:3.21",
			Namespace:     "default",
			CPURequest:    "50m",
			CPULimit:      "100m",
			MemoryRequest: "32Mi",
			MemoryLimit:   "64Mi",
		},
	}

	runID := "integ-cancel-12345"
	err := r.Run(ctx, RunOptions{
		ProjectID: "cancel-proj",
		RunID:     runID,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Cancel
	err = r.Cancel(ctx, runID)
	if err != nil {
		t.Fatalf("Cancel failed: %v", err)
	}

	// Verify Job is deleted
	time.Sleep(time.Second) // give K3s a moment
	jobs, _ := testK8sClient.BatchV1().Jobs("default").List(ctx, metav1.ListOptions{
		LabelSelector: "run-id=" + runID,
	})
	if len(jobs.Items) != 0 {
		t.Errorf("expected 0 jobs after cancel, got %d", len(jobs.Items))
	}
}

func TestInteg_K8sRunner_MultipleParallelJobs(t *testing.T) {
	ctx := context.Background()

	r := &KubernetesRunner{
		client: testK8sClient,
		config: Config{
			AgentImage:    "alpine:3.21",
			Namespace:     "default",
			CPURequest:    "50m",
			CPULimit:      "100m",
			MemoryRequest: "32Mi",
			MemoryLimit:   "64Mi",
		},
	}

	// Create 3 parallel runs
	for i := 0; i < 3; i++ {
		err := r.Run(ctx, RunOptions{
			ProjectID: "parallel-proj",
			RunID:     fmt.Sprintf("parallel-run-%d-abc", i),
		})
		if err != nil {
			t.Fatalf("Run %d failed: %v", i, err)
		}
	}

	// Verify all 3 exist
	jobs, _ := testK8sClient.BatchV1().Jobs("default").List(ctx, metav1.ListOptions{
		LabelSelector: "project-id=parallel-proj",
	})
	if len(jobs.Items) != 3 {
		t.Errorf("expected 3 parallel jobs, got %d", len(jobs.Items))
	}
}
