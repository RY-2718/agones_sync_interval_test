package agones_sync_interval

import (
	"fmt"

	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	allocationv1 "agones.dev/agones/pkg/apis/allocation/v1"
	autoscalingv1 "agones.dev/agones/pkg/apis/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	containerName  = "simple-game-server"
	containerImage = "gcr.io/agones-images/simple-game-server:0.13"
	fleetBaseName  = "simple-game-server"
	fasBaseName    = "simple-game-server-autoscaler"
)

func newTestGameServer(namespace string) *agonesv1.GameServer {
	return &agonesv1.GameServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      containerName,
			Namespace: namespace,
		},
		Spec: agonesv1.GameServerSpec{
			Ports: []agonesv1.GameServerPort{
				{
					Name:          "default",
					ContainerPort: 7654,
				},
			},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  containerName,
							Image: containerImage,
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("20m"),
									corev1.ResourceMemory: resource.MustParse("64Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("20m"),
									corev1.ResourceMemory: resource.MustParse("64Mi"),
								},
							},
						},
					},
				},
			},
		},
		Status: agonesv1.GameServerStatus{},
	}
}

func newTestFleet(namespace string, name string) *agonesv1.Fleet {
	gss := newTestGameServer(namespace).Spec
	return &agonesv1.Fleet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", fleetBaseName, name),
			Namespace: namespace,
		},
		Spec: agonesv1.FleetSpec{
			Replicas: 2,
			Template: agonesv1.GameServerTemplateSpec{Spec: gss},
		},
	}
}

func newTestFleetAutoscaler(namespace string, fleet *agonesv1.Fleet, name string, interval int32) *autoscalingv1.FleetAutoscaler {
	fm := fleet.ObjectMeta
	return &autoscalingv1.FleetAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", fasBaseName, name),
			Namespace: namespace,
		},
		Spec: autoscalingv1.FleetAutoscalerSpec{
			FleetName: fm.Name,
			Policy: autoscalingv1.FleetAutoscalerPolicy{
				Type: autoscalingv1.BufferPolicyType,
				Buffer: &autoscalingv1.BufferPolicy{
					MaxReplicas: 10,
					MinReplicas: 0,
					BufferSize:  intstr.IntOrString{IntVal: 2},
				},
			},
			Sync: &autoscalingv1.FleetAutoscalerSync{
				Type: autoscalingv1.FixedIntervalSyncType,
				FixedInterval: autoscalingv1.FixedIntervalSync{
					Seconds: interval,
				},
			},
		},
	}
}

func newGameServerAllocation(fleetName string) *allocationv1.GameServerAllocation {
	return &allocationv1.GameServerAllocation{
		Spec: allocationv1.GameServerAllocationSpec{
			Selectors: []allocationv1.GameServerSelector{
				{
					LabelSelector: metav1.LabelSelector{
						MatchLabels: map[string]string{
							agonesv1.FleetNameLabel: fleetName,
						},
					},
				},
			},
		},
	}
}
