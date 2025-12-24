package kubernetes

import (
	"context"
	"time"

	"github.com/techtacles/sysmonitoring/internal/logging"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type NamespaceInfo struct {
	Name         string
	CreationTime time.Time
}

type ServiceInfo struct {
	Name           string
	Namespace      string
	CreationTime   time.Time
	ClusterIp      string
	ClusterIps     []string
	ExternalIps    []string
	ExternalName   string
	LoadBalancerIp string
	PortName       string
	NodePort       int32
	Protocol       corev1.Protocol
	Port           int32
	TargetPort     intstr.IntOrString
	Type           corev1.ServiceType
}

type PodInfo struct {
	Name              string
	CreationTimestamp time.Time
	Namespace         string
	HostIP            string
	PodIP             string
	Phase             corev1.PodPhase
	ContainerImage    string
	ContainerName     string
	ContainerPort     []corev1.ContainerPort
	VolumeName        string
	SchedulerName     string
	NodeName          string
}

type NodeInfo struct {
	Name              string
	CreationTimestamp time.Time
	VolumesAttached   []corev1.AttachedVolume
	VolumesInUse      []corev1.UniqueVolumeName
	Addressed         []corev1.NodeAddress
	Unschedulable     bool
	PodCIDRs          []string
}

type PersistentVolumeInfo struct {
	Name             string
	CreationTime     time.Time
	Capacity         corev1.ResourceList
	AccessModes      []corev1.PersistentVolumeAccessMode
	ReclaimPolicy    corev1.PersistentVolumeReclaimPolicy
	Status           corev1.PersistentVolumePhase
	StorageClassName string
	VolumeMode       *corev1.PersistentVolumeMode
}

type PersistentVolumeClaimInfo struct {
	Name             string
	Namespace        string
	CreationTime     time.Time
	Status           corev1.PersistentVolumeClaimPhase
	AccessModes      []corev1.PersistentVolumeAccessMode
	StorageClassName string
	VolumeMode       corev1.PersistentVolumeMode
	Capacity         corev1.ResourceList
}

func getNamespaceInfo() ([]NamespaceInfo, error) {
	clientset := getClientset()
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(),
		metav1.ListOptions{})

	if err != nil {
		logging.Error(logtag, "error getting namespaces", err)
		panic(err.Error())
	}

	results := make([]NamespaceInfo, 0, len(namespaces.Items))

	for _, namespace := range namespaces.Items {
		results = append(results, NamespaceInfo{Name: namespace.Name,
			CreationTime: namespace.CreationTimestamp.Time})

	}
	return results, nil
}

func getServicesInfo() ([]ServiceInfo, error) {
	clientset := getClientset()
	services, err := clientset.CoreV1().Services(metav1.NamespaceAll).List(context.Background(),
		metav1.ListOptions{})

	if err != nil {
		logging.Error(logtag, "error getting services", err)
		panic(err.Error())
	}

	results := make([]ServiceInfo, 0, len(services.Items))

	for _, service := range services.Items {
		if len(service.Name) == 0 {
			continue
		}

		commonInfo := ServiceInfo{
			Name:           service.Name,
			CreationTime:   service.CreationTimestamp.Time,
			Namespace:      service.Namespace,
			ClusterIp:      service.Spec.ClusterIP,
			ClusterIps:     service.Spec.ClusterIPs,
			ExternalIps:    service.Spec.ExternalIPs,
			ExternalName:   service.Spec.ExternalName,
			LoadBalancerIp: service.Spec.LoadBalancerIP,
			Type:           service.Spec.Type,
		}

		if len(service.Spec.Ports) == 0 {
			results = append(results, commonInfo)
			continue
		}

		for _, port := range service.Spec.Ports {
			info := commonInfo
			info.PortName = port.Name
			info.NodePort = port.NodePort
			info.Protocol = port.Protocol
			info.Port = port.Port
			info.TargetPort = port.TargetPort
			results = append(results, info)
		}
	}
	return results, nil
}

func getPodsInfo() ([]PodInfo, error) {
	clientset := getClientset()
	pods, err := clientset.CoreV1().Pods(metav1.NamespaceAll).List(context.Background(),
		metav1.ListOptions{})

	if err != nil {
		logging.Error(logtag, "error getting pods", err)
		panic(err.Error())
	}

	results := make([]PodInfo, 0, len(pods.Items))

	for _, pod := range pods.Items {
		commonInfo := PodInfo{Name: pod.Name,
			CreationTimestamp: pod.CreationTimestamp.Time,
			Namespace:         pod.Namespace,
			HostIP:            pod.Status.HostIP,
			PodIP:             pod.Status.PodIP,
			Phase:             pod.Status.Phase,
			SchedulerName:     pod.Spec.SchedulerName,
			NodeName:          pod.Spec.NodeName}

		if len(pod.Spec.Containers) == 0 {
			results = append(results, commonInfo)
			continue
		}
		if len(pod.Spec.Volumes) == 0 {
			results = append(results, commonInfo)
			continue
		}
		for _, container := range pod.Spec.Containers {
			info := commonInfo
			info.ContainerName = container.Name
			info.ContainerImage = container.Image
			info.ContainerPort = container.Ports
			results = append(results, info)
		}
		for _, volume := range pod.Spec.Volumes {
			info := commonInfo
			info.VolumeName = volume.Name
			results = append(results, info)
		}
		results = append(results, commonInfo)
	}

	return results, nil
}

func getNodesInfo() ([]NodeInfo, error) {
	clientset := getClientset()
	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})

	if err != nil {
		logging.Error(logtag, "error getting nodes", err)
		panic(err.Error())
	}

	results := make([]NodeInfo, 0, len(nodes.Items))
	for _, node := range nodes.Items {
		results = append(results, NodeInfo{Name: node.Name,
			CreationTimestamp: node.CreationTimestamp.Time,
			VolumesAttached:   node.Status.VolumesAttached,
			VolumesInUse:      node.Status.VolumesInUse,
			Addressed:         node.Status.Addresses,
			Unschedulable:     node.Spec.Unschedulable,
			PodCIDRs:          node.Spec.PodCIDRs})
	}
	return results, nil
}

func getPersistentVolumesInfo() ([]PersistentVolumeInfo, error) {
	clientset := getClientset()
	pvs, err := clientset.CoreV1().PersistentVolumes().List(context.Background(), metav1.ListOptions{})

	if err != nil {
		logging.Error(logtag, "error getting persistent volumes", err)
		panic(err.Error())
	}

	results := make([]PersistentVolumeInfo, 0, len(pvs.Items))
	for _, pv := range pvs.Items {
		results = append(results, PersistentVolumeInfo{
			Name:             pv.Name,
			CreationTime:     pv.CreationTimestamp.Time,
			Capacity:         pv.Spec.Capacity,
			AccessModes:      pv.Spec.AccessModes,
			ReclaimPolicy:    pv.Spec.PersistentVolumeReclaimPolicy,
			Status:           pv.Status.Phase,
			StorageClassName: pv.Spec.StorageClassName,
			VolumeMode:       pv.Spec.VolumeMode,
		})
	}
	return results, nil
}

func getPersistentVolumeClaimsInfo() ([]PersistentVolumeClaimInfo, error) {
	clientset := getClientset()
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims(metav1.NamespaceAll).List(context.Background(), metav1.ListOptions{})

	if err != nil {
		logging.Error(logtag, "error getting persistent volume claims", err)
		panic(err.Error())
	}

	results := make([]PersistentVolumeClaimInfo, 0, len(pvcs.Items))
	for _, pvc := range pvcs.Items {
		results = append(results, PersistentVolumeClaimInfo{
			Name:             pvc.Name,
			Namespace:        pvc.Namespace,
			CreationTime:     pvc.CreationTimestamp.Time,
			Status:           pvc.Status.Phase,
			AccessModes:      pvc.Spec.AccessModes,
			StorageClassName: *pvc.Spec.StorageClassName,
			VolumeMode:       *pvc.Spec.VolumeMode,
			Capacity:         pvc.Status.Capacity,
		})
	}
	return results, nil
}
