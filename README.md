# kube-utils (kut)
Some k8s dev utils.

## Install
```shell
go install github.com/Eikykun/kut@latest
```
## Usage

**Merge Kubeconfig** 
```shell
kut merge --from=config-1.yaml --to=config-2.yaml --overwrite=context-1
```
