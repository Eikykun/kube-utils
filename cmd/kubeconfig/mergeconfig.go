package kubeconfig

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	v1 "k8s.io/client-go/tools/clientcmd/api/v1"
	"sigs.k8s.io/yaml"
)

func Cmd() *cobra.Command {
	return mergeConfigCmd
}

// mergeConfigCmd represents the mergeConfig command
var mergeConfigCmd = &cobra.Command{
	Use:   "merge",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mergeConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	mergeConfigCmd.Flags().StringVar(&fromKubeconfigDir, "from", "", "Kubeconfig file that needs to be merged")
	mergeConfigCmd.Flags().StringVar(&toKubeconfigDir, "to", "", "Merged kubeconfig file.")
	mergeConfigCmd.Flags().StringVar(&overwrite, "overwrite", "", "Overwrite context/cluster/user names all in one")
}

var (
	fromKubeconfigDir string
	toKubeconfigDir   string
	overwrite         string
)

func run() error {
	if fromKubeconfigDir == "" || toKubeconfigDir == "" {
		return fmt.Errorf("fromKubeconfigDir or toKubeconfigDir is empty")
	}
	toConfig, err := readConfigYaml(toKubeconfigDir)
	if err != nil {
		return err
	}
	fromConfig, err := readConfigYaml(toKubeconfigDir)
	if err != nil {
		return err
	}
	if err = validateSingleConfig(fromConfig); err != nil {
		return fmt.Errorf("validateSingleConfig failed: %v", err)
	}
	getter := newSingleGetter(fromConfig, overwrite)
	toConfig.Contexts = append(toConfig.Contexts, getter.Context())
	toConfig.Clusters = append(toConfig.Clusters, getter.Cluster())
	toConfig.AuthInfos = append(toConfig.AuthInfos, getter.User())
	val, _ := yaml.Marshal(toConfig)
	fmt.Printf("Merged config:\n%s\n", string(val))
	return os.WriteFile(toKubeconfigDir, val, 0644)
}

func readConfigYaml(dir string) (config *v1.Config, err error) {
	val, err := os.ReadFile(dir)
	if err != nil {
		return nil, err
	}
	return unmarshalConfig(val)
}

func unmarshalConfig(configByte []byte) (*v1.Config, error) {
	config := &v1.Config{}
	if err := yaml.Unmarshal(configByte, config); err != nil {
		return nil, err
	}
	return config, nil
}

func validateSingleConfig(config *v1.Config) error {
	if len(config.Clusters) != 1 || len(config.AuthInfos) != 1 || config.Clusters[0].Cluster.Server == "" {
		return fmt.Errorf("invalid single config")
	}
	return nil
}

func newSingleGetter(cfg *v1.Config, overwrite string) *singleGetter {
	return &singleGetter{
		overwrite: overwrite,
		cfg:       cfg,
	}
}

type singleGetter struct {
	overwrite string
	cfg       *v1.Config
}

func (g *singleGetter) User() v1.NamedAuthInfo {
	if overwrite == "" {
		return g.cfg.AuthInfos[0]
	}
	return v1.NamedAuthInfo{
		Name:     g.overwrite,
		AuthInfo: g.cfg.AuthInfos[0].AuthInfo,
	}
}

func (g *singleGetter) Cluster() v1.NamedCluster {
	if overwrite == "" {
		return g.cfg.Clusters[0]
	}
	return v1.NamedCluster{
		Name:    g.overwrite,
		Cluster: g.cfg.Clusters[0].Cluster,
	}
}

func (g *singleGetter) Context() v1.NamedContext {
	if overwrite == "" {
		return g.cfg.Contexts[0]
	}
	return v1.NamedContext{
		Name: g.overwrite,
		Context: v1.Context{
			Cluster:  g.overwrite,
			AuthInfo: g.overwrite,
		},
	}
}
