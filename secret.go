package k8s

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// DockerRegistrySecretData 定义 dockerconfigjson 的结构
type DockerRegistrySecretData struct {
	Auths map[string]DockerAuthEntry `json:"auths"`
}

type DockerAuthEntry struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Auth     string `json:"auth"`
}

// CreateDockerRegistrySecret 创建镜像仓库 Secret
func CreateDockerRegistrySecret(clientset *kubernetes.Clientset, namespace, name, server, username, password string) error {
	secretClient := clientset.CoreV1().Secrets(namespace)

	// 检查 Secret 是否存在
	if _, err := secretClient.Get(context.TODO(), name, metav1.GetOptions{}); err == nil {
		fmt.Printf("✅ Secret %q 已存在于命名空间 %q。\n", name, namespace)
		return nil
	}

	// 构造 dockerconfigjson 数据结构
	auth := DockerRegistrySecretData{
		Auths: map[string]DockerAuthEntry{
			server: {
				Username: username,
				Password: password,
				Auth:     basicAuth(username, password), // Base64 编码
			},
		},
	}

	jsonData, err := json.Marshal(auth)
	if err != nil {
		return fmt.Errorf("序列化 dockerconfigjson 失败: %v", err)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Type: corev1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			".dockerconfigjson": jsonData,
		},
	}

	_, err = secretClient.Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("创建 Secret 失败: %v", err)
	}

	fmt.Printf("✅ Docker 仓库 Secret %q 已在命名空间 %q 创建成功。\n", name, namespace)
	return nil
}

// basicAuth 拼接 Base64 授权字符串
func basicAuth(username, password string) string {
	raw := fmt.Sprintf("%s:%s", username, password)
	return base64.StdEncoding.EncodeToString([]byte(raw))
}
