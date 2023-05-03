# accuknox-cli

**accuknox-cli** is a client tool to help manage [KubeArmor](https://github.com/kubearmor/KubeArmor) and [Discovery Engine](https://github.com/accuknox/discovery-engine).

## Installation

### Installing from Source 

Build accuknox-cli from source if you want to test the latest (pre-release) accuknox-cli version.

```
git clone https://github.com/accuknox/accuknox-cli.git
cd accuknox-cli
make install
```

## Usage

```
CLI Utility to help manage KubeArmor and Discovery Engine

KubeArmor is a container-aware runtime security enforcement system that
restricts the behavior (such as process execution, file access, and networking
operation) of containers at the system level.

Discovery Engine discovers the security posture for your workloads and auto-discovers the policy-set required to put the workload in least-permissive mode. The engine leverages the rich visibility provided by KubeArmor and Cilium to auto discover the systems and network security posture.

Usage:
  accuknox-cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  discover    Discover applicable policies
  help        Help about any command
  install     Install KubeArmor in a Kubernetes Cluster
  logs        Observe Logs from KubeArmor
  probe       Checks for supported KubeArmor features in the current environment
  profile     Profiling of logs
  recommend   Recommend Policies
  rotate-tls  Rotate webhook controller tls certificates
  selfupdate  selfupdate this cli tool
  summary     Observability from discovery engine
  sysdump     Collect system dump information for troubleshooting and error report
  uninstall   Uninstall KubeArmor from a Kubernetes Cluster
  version     Display version information
  vm          VM commands for kvmservice

Flags:
      --context string      Name of the kubeconfig context to use
  -h, --help                help for accuknox-cli
      --kubeconfig string   Path to the kubeconfig file to use

Use "accuknox-cli [command] --help" for more information about a command.
```
