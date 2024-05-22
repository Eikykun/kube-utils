# kube-utils
Some k8s dev utils.

## Install
```shell
go install github.com/Eikykun/kube-utils@latest
```
## Usage

**Merge Kubeconfig** 
```shell
kube-utils merge --from=config-1.yaml --to=config-2.yaml --overwrite=context-1
```
