package kubeconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeConfig(t *testing.T) {
	from := "kind: Config\napiVersion: v1\nusers:\n- name: test-cluster-1-user\n  user:\n    client-certificate-data: IA==\n    client-key-data: IA==\nclusters:\n- name: test-cluster-1\n  cluster:\n    certificate-authority-data: IA==\n    server: https://apiserver.test.svc.net:6443\ncontexts:\n- context:\n    cluster: test-cluster-1\n    user: test-cluster-1-user\n  name: default\ncurrent-context: default\npreferences: {}"
	to := "kind: Config\napiVersion: v1\nusers:\n- name: test-cluster-2-user\n  user:\n    client-certificate-data: IA==\n    client-key-data: IA==\nclusters:\n- name: test-cluster-2\n  cluster:\n    certificate-authority-data: IA==\n    server: https://apiserver.test.svc.net:6443\ncontexts:\n- context:\n    cluster: test-cluster-2\n    user: test-cluster-2-user\n  name: default\ncurrent-context: default\npreferences: {}"
	fromKubeconfigDir = "./from-test.yaml"
	toKubeconfigDir = "./to-test.yaml"
	overwrite = "overwrite-cluster-name"
	os.WriteFile(fromKubeconfigDir, []byte(from), 0644)
	os.WriteFile(toKubeconfigDir, []byte(to), 0644)
	assert.Nil(t, run(), "Fail to merge kubeconfig")
}
