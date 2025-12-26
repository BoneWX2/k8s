(# k8s 客户端与 HTTP 服务

这是一个轻量的示例项目，用于展示如何从本地 kubeconfig 创建 Kubernetes 客户端，并通过 HTTP 接口暴露资源查询（目前实现 `/pods` 列表接口）。README 包含项目请求逻辑、配置查找顺序、如何运行与示例。

## 项目结构（相关文件）
- `client/client.go`：Kubernetes 客户端封装，提供 `CreateClientFromConfig`、`NewFromKubeconfig` 与 `ListPods`。
- `server/server.go`：启动 HTTP 服务、创建客户端并注册路由（`/pods`）。
- `http/pod.go`：HTTP handler（如果存在）用于处理 `/pods` 请求并调用 `client`。
- `main.go`：程序入口，仅负责解析参数并调用 `server.Start`（启动服务）。

## 配置查找顺序
客户端创建时将按下列顺序查找 kubeconfig：

1. 函数入参 `configPath`（来自 `--kubeconfig` 参数）
2. 环境变量 `KUBECONFIG`
3. 项目相对路径 `./kubeconfig/config.yaml`（如果该文件存在）
4. 用户主目录下的 `~/.kube/config`（如果存在）
5. 回退到集群内配置（in-cluster）

这样你可以把配置放在项目的 `kubeconfig/config.yaml` 中并直接运行，或者通过 `--kubeconfig` 或 `KUBECONFIG` 指定路径。

## 请求逻辑（高层）
当接收到 HTTP 请求（例如 GET /pods）时，处理流程如下：

1. `server` 创建或接收一个 `KubeClient`（通过 `CreateClientFromConfig` 使用上文的查找逻辑）。
2. HTTP handler（`/pods`）调用 `KubeClient.ListPods(namespace)` 来获取 Pod 列表。
3. 将 Pod 列表按 JSON 格式返回给客户端；若发生错误，返回 500 并包含错误信息（中文提示）。

## HTTP 路由
- GET /pods?ns=<namespace>
	- 参数：`ns`（可选，默认 `default`），指定命名空间；传空字符串表示所有命名空间。
	- 返回示例：
		```json
		["default/nginx-7bb7c576c4-abcde", "kube-system/coredns-5d4f6f6f6-xyz"]
		```

## 运行示例（PowerShell）

- 使用项目内的 kubeconfig：
```powershell
go run . --kubeconfig .\kubeconfig\config.yaml --ns default
```

- 使用环境变量：
```powershell
#$env:KUBECONFIG = "D:\\Go\\repo\\k8s\\kubeconfig\\config.yaml"
go run . --ns default
```

服务启动后，可以用 curl 访问（假设监听本地 8080）：
```powershell
curl http://127.0.0.1:8080/pods?ns=default
```

## 常见问题
- 找不到 kubeconfig：确认 `--kubeconfig` 路径正确，或者将配置放到 `./kubeconfig/config.yaml` 或 `~/.kube/config`，或者设置环境变量 `KUBECONFIG`。
- 权限错误：确认 kubeconfig 中的凭据有列出 Pods 的权限。

## 后续改进建议
- 支持分页或 limit 参数，以处理大量 Pod 返回。
- 使用结构化日志（例如 `klog`）替代 `fmt`，便于生产环境排查。
- 增加单元测试与集成测试，模拟 Kubernetes API 返回值。
)

