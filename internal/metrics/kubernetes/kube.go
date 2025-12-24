package kubernetes

import (
	"path/filepath"

	"github.com/techtacles/sysmonitoring/internal/logging"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var ExplicitKubeconfigPath string

var logtag string = "kubernetes"

type KubeInfo struct {
	DeploymentStats            []DeploymentInfo
	NamespaceStats             []NamespaceInfo
	PersistentVolumeStats      []PersistentVolumeInfo
	PersistentVolumeClaimStats []PersistentVolumeClaimInfo
	NodeStats                  []NodeInfo
	PodStats                   []PodInfo
	ServiceStats               []ServiceInfo
}

func (k *KubeInfo) Collect() error {
	deploy, err := getDeploymentInfo()
	if err != nil {
		return err
	}
	k.DeploymentStats = deploy

	namespace, err := getNamespaceInfo()
	if err != nil {
		return err
	}
	k.NamespaceStats = namespace

	pv, err := getPersistentVolumesInfo()
	if err != nil {
		return err
	}
	k.PersistentVolumeStats = pv

	pvc, err := getPersistentVolumeClaimsInfo()
	if err != nil {
		return err
	}
	k.PersistentVolumeClaimStats = pvc

	node, err := getNodesInfo()
	if err != nil {
		return err
	}
	k.NodeStats = node

	pod, err := getPodsInfo()
	if err != nil {
		return err
	}
	k.PodStats = pod

	service, err := getServicesInfo()
	if err != nil {
		return err
	}
	k.ServiceStats = service
	logging.Info(logtag, "successfully collected kubernetes metrics")

	return nil
}

func GetKubeConfigPath() string {
	if ExplicitKubeconfigPath != "" {
		return ExplicitKubeconfigPath
	}
	if home := homedir.HomeDir(); home != "" {
		return filepath.Join(home, ".kube", "config")
	}
	return ""
}

func getClientset() *kubernetes.Clientset {
	path := GetKubeConfigPath()

	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		logging.Error(logtag, "error getting config", err)
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logging.Error(logtag, "error getting clientset", err)
		panic(err.Error())
	}
	return clientset
}
