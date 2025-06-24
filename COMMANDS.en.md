# Test Guide for eks-review-cli

This document provides a list of commands to test **eks-review-cli**. Each section briefly describes the command and the available flags.

**Prerequisites**
- Compile the CLI with `go build -o eks-review`.
- Run commands from the directory where the executable is located or ensure it is in your PATH.
- Configure your `~/.kube/config` to point to an accessible Kubernetes cluster with proper permissions.

---

## Root command (`eks-review`)

- **General help**
  ```bash
  ./eks-review --help
  ```
- **Global `--verbose` / `-v` flag**
  - Enables detailed output (DEBUG logs).
  - Example usage:
    ```bash
    ./eks-review --verbose monitor status
    ./eks-review -v monitor events -n default
    ```

---

## `monitor` subcommand
Commands for monitoring and cluster visibility.

- **Help for `monitor`**
  ```bash
  ./eks-review monitor --help
  ```

### 1. `eks-review monitor status`
Shows a summary of Pods, Deployments, Services and Ingresses.

Sample commands:
```bash
./eks-review monitor status
./eks-review monitor status --namespace kube-system
./eks-review monitor status --all-namespaces
./eks-review monitor status --help
```

### 2. `eks-review monitor events`
Displays recent cluster events.

```bash
./eks-review monitor events
./eks-review monitor events --namespace default
./eks-review monitor events --type Warning
./eks-review monitor events -n kube-system -T Warning
./eks-review monitor events --help
```

### 3. `eks-review monitor nodes`
Shows detailed information about cluster nodes.

```bash
./eks-review monitor nodes
./eks-review monitor nodes --help
```

### 4. `eks-review monitor logs`
Prints logs from a pod, deployment or service.

```bash
./eks-review monitor logs --pod <pod-name>
./eks-review monitor logs --deployment <deployment-name>
./eks-review monitor logs --service <service-name>
./eks-review monitor logs --pod <pod-name> -n <namespace>
./eks-review monitor logs --pod <pod-name> -c <container-name>
./eks-review monitor logs --pod <pod-name> -f
./eks-review monitor logs --pod <pod-name> -p
./eks-review monitor logs --pod <pod-name> --grep "text"
./eks-review monitor logs --pod <pod-name> --tail 20
./eks-review monitor logs --help
```

### 5. `eks-review monitor get <resource> [name]`
Lists resources, similar to `kubectl get`.

Supported resources:
- `pods` (`po`)
- `services` (`svc`)
- `daemonsets` (`ds`)
- `jobs` (`job`)
- `cronjobs` (`cj`)
- `namespaces` (`ns`)
- `serviceaccounts` (`sa`)

Common flags:
- `-n, --namespace <namespace>`
- `-A, --all-namespaces`
- `-l, --selector <label_selector>`
- `-o, --output <format>` (`wide`, `json`, `yaml`)

*(The resource `namespaces` does not use `-n` or `-A`.)*

---

## Planned subcommands (placeholders)
These commands exist but only print a message because their full implementation is still pending.

### `eks-review security`
Audits Network Policies, RBAC, container images and Secrets.
```bash
./eks-review security
./eks-review security --help
```

### `eks-review optimize`
Identifies unused resources and autoscaling configuration.
```bash
./eks-review optimize
./eks-review optimize --help
```

### `eks-review diagnose`
Diagnoses problems in Pods, Services and Ingresses.
```bash
./eks-review diagnose
./eks-review diagnose --help
```

---

## Help
For help on any command or subcommand, run:
```bash
eks-review [command] [subcommand] --help
```

---

Use real resource names when testing. Observe the command output and error handling. The `--help` flag shows all available flags and descriptions.
