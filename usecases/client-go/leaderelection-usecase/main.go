package main

import (
	"context"
	"fmt"
	"os"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/util/workqueue"
)

func main() {
	// 创建一个 Kubernetes 客户端
	client := fake.NewSimpleClientset()

	// 定义锁的配置
	lockName := "example-lock"
	lockNamespace := "default"
	id := os.Getenv("POD_NAME") // 使用 Pod 名称作为唯一标识符
	if id == "" {
		id = "example-id"
	}

	// 创建一个 ResourceLock
	lock := &resourcelock.LeaseLock{
		LeaseMeta: v1.ObjectMeta{
			Name:      lockName,
			Namespace: lockNamespace,
		},
		Client: client.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id,
		},
	}

	// 定义 Leader Election 的回调
	callbacks := leaderelection.LeaderCallbacks{
		OnStartedLeading: func(ctx context.Context) {
			fmt.Println("I am the leader now!")
			// 在这里执行你的主节点逻辑
			runLeaderTasks(ctx)
		},
		OnStoppedLeading: func() {
			fmt.Println("I am no longer the leader.")
			// 在这里处理失去领导权的逻辑
		},
		OnNewLeader: func(identity string) {
			if identity == id {
				fmt.Println("I just became the leader!")
			} else {
				fmt.Printf("New leader elected: %s\n", identity)
			}
		},
	}

	// 配置 Leader Election
	lec := leaderelection.LeaderElectionConfig{
		Lock:            lock,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks:       callbacks,
		ReleaseOnCancel: true,
	}

	// 创建 LeaderElector
	leaderElector, err := leaderelection.NewLeaderElector(lec)
	if err != nil {
		fmt.Printf("Failed to create leader elector: %v\n", err)
		return
	}

	// 启动 Leader Election
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Starting leader election...")
	leaderElector.Run(ctx)
}

// runLeaderTasks 是主节点的任务逻辑
func runLeaderTasks(ctx context.Context) {
	queue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[struct{}]())
	defer queue.ShutDown()

	// 模拟一个任务
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping leader tasks...")
				return
			default:
				fmt.Println("Running leader tasks...")
				time.Sleep(5 * time.Second)
			}
		}
	}()
}
