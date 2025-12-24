package kubernetes

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/techtacles/sysmonitoring/internal/logging"
)

var namespace string = "" //for all namespaces

type DeploymentInfo struct {
	Name                string
	Namespace           string
	AvailableReplicas   int32
	ReadyReplicas       int32
	UpdatedReplicas     int32
	TerminatingReplicas *int32
	TotalReplicas       int32
	CreationTime        time.Time
}

func getDeploymentInfo() ([]DeploymentInfo, error) {

	clientset := getClientset()

	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		logging.Error(logtag, "error getting deployments", err)
		panic(err.Error())
	}

	results := make([]DeploymentInfo, 0, len(deployments.Items))
	for _, deployment := range deployments.Items {

		results = append(results, DeploymentInfo{Name: deployment.Name,
			Namespace:           deployment.Namespace,
			AvailableReplicas:   deployment.Status.AvailableReplicas,
			ReadyReplicas:       deployment.Status.ReadyReplicas,
			UpdatedReplicas:     deployment.Status.UpdatedReplicas,
			TerminatingReplicas: deployment.Status.TerminatingReplicas,
			TotalReplicas:       deployment.Status.Replicas,
			CreationTime:        deployment.CreationTimestamp.Time})
	}
	return results, nil
}
