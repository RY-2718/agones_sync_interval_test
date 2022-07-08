package agones_sync_interval

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	allocationv1 "agones.dev/agones/pkg/apis/allocation/v1"
	e2eframework "agones.dev/agones/test/e2e/framework"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	framework *e2eframework.Framework
)

const namespacePrefix = "sync-interval-test"

func TestMain(m *testing.M) {
	fw, err := e2eframework.NewFromFlags()
	if err != nil {
		log.Fatalf("failed to init e2e framework: %v", err)
	}
	framework = fw

	var exitCode int
	defer func() {
		_ = framework.CleanUp(namespacePrefix)
		os.Exit(exitCode)
	}()
	exitCode = m.Run()
}

func TestCreateFleetAutoscaler(t *testing.T) {
	fmt.Println("getting k8s config...")

	namespace := fmt.Sprintf("%s-%s", namespacePrefix, uuid.Must(uuid.NewRandom()).String())
	createNamespace(t, namespace)
	defer deleteNamespace(t, namespace)

	ctx := context.Background()

	t.Run("allocation with slow resync interval must fail", func(t *testing.T) {
		fleet, err := framework.AgonesClient.AgonesV1().Fleets(namespace).Create(ctx, newTestFleet(namespace, "slow"), metav1.CreateOptions{})
		require.NoError(t, err)
		_, err = framework.AgonesClient.AutoscalingV1().FleetAutoscalers(namespace).Create(ctx, newTestFleetAutoscaler(namespace, fleet, "slow", 60), metav1.CreateOptions{})
		require.NoError(t, err)

		// 準備待ち
		framework.AssertFleetCondition(t, fleet, e2eframework.FleetReadyCount(fleet.Spec.Replicas))
		// 初期値使い切る
		for i := 0; int32(i) < fleet.Spec.Replicas; i++ {
			alloc, err := framework.AgonesClient.AllocationV1().GameServerAllocations(namespace).Create(ctx, newGameServerAllocation(fleet.Name), metav1.CreateOptions{})
			assert.NoError(t, err)
			assert.Equal(t, allocationv1.GameServerAllocationAllocated, alloc.Status.State)
		}

		time.Sleep(20 * time.Second)
		alloc, err := framework.AgonesClient.AllocationV1().GameServerAllocations(namespace).Create(ctx, newGameServerAllocation(fleet.Name), metav1.CreateOptions{})
		assert.NoError(t, err)
		assert.Equal(t, allocationv1.GameServerAllocationUnAllocated, alloc.Status.State)
	})

	t.Run("allocation with fast resync interval must succeed", func(t *testing.T) {
		fleet, err := framework.AgonesClient.AgonesV1().Fleets(namespace).Create(ctx, newTestFleet(namespace, "fast"), metav1.CreateOptions{})
		require.NoError(t, err)
		_, err = framework.AgonesClient.AutoscalingV1().FleetAutoscalers(namespace).Create(ctx, newTestFleetAutoscaler(namespace, fleet, "fast", 15), metav1.CreateOptions{})
		require.NoError(t, err)

		// 準備待ち
		framework.AssertFleetCondition(t, fleet, e2eframework.FleetReadyCount(fleet.Spec.Replicas))
		// 初期値使い切る
		for i := 0; int32(i) < fleet.Spec.Replicas; i++ {
			alloc, err := framework.AgonesClient.AllocationV1().GameServerAllocations(namespace).Create(ctx, newGameServerAllocation(fleet.Name), metav1.CreateOptions{})
			assert.NoError(t, err)
			assert.Equal(t, allocationv1.GameServerAllocationAllocated, alloc.Status.State)
		}

		time.Sleep(20 * time.Second)
		alloc, err := framework.AgonesClient.AllocationV1().GameServerAllocations(namespace).Create(ctx, newGameServerAllocation(fleet.Name), metav1.CreateOptions{})
		assert.NoError(t, err)
		assert.Equal(t, allocationv1.GameServerAllocationAllocated, alloc.Status.State)
	})
}

func createNamespace(t *testing.T, namespace string) {
	t.Helper()
	assert.NoError(t, framework.CreateNamespace(namespace))
}

func deleteNamespace(t *testing.T, namespace string) {
	t.Helper()
	assert.NoError(t, framework.DeleteNamespace(namespace))
}
