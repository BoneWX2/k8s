package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClient 创建 K8s 客户端
func NewClient(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig) //连接k8s集群
	if err != nil {
		return nil, fmt.Errorf("加载 kubeconfig 失败: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建客户端失败: %v", err)
	}

	return clientset, nil
}

// EnsureNamespace 确保命名空间存在，不存在则创建
func EnsureNamespace(clientset *kubernetes.Clientset, name string) error {
	_, err := clientset.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err == nil {
		fmt.Printf("✅ 命名空间 %q 已存在。\n", name)
		return nil
	}

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("创建命名空间失败: %v", err)
	}

	fmt.Printf("✅ 命名空间 %q 创建成功。\n", name)
	return nil
}
