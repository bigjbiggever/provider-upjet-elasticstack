package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/pipeline"

	"github.com/bigjbiggever/provider-elasticstack/config"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		panic("root directory is required to be given as argument")
	}
	rootDir := os.Args[1]
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		panic(fmt.Sprintf("cannot calculate the absolute path with %s", rootDir))
	}
	pipeline.Run(config.GetProvider(), config.GetProviderNamespaced(), absRootDir)
	if err := patchClusterManagedResourceSpec(absRootDir); err != nil {
		panic(fmt.Sprintf("cannot patch generated cluster managed resource specs: %v", err))
	}
}

func patchClusterManagedResourceSpec(rootDir string) error {
	clusterAPIsDir := filepath.Join(rootDir, "apis", "cluster")
	return filepath.WalkDir(clusterAPIsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasPrefix(d.Name(), "zz_") || !strings.HasSuffix(d.Name(), "_types.go") {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(b)
		if !strings.Contains(content, "v1.ResourceSpec `json:\",inline\"`") {
			return nil
		}
		content = strings.ReplaceAll(content, "v1.ResourceSpec `json:\",inline\"`", "v2.ManagedResourceSpec `json:\",inline\"`")
		if !strings.Contains(content, `v2 "github.com/crossplane/crossplane-runtime/v2/apis/common/v2"`) {
			content = strings.Replace(
				content,
				`v1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"`+"\n",
				`v1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"`+"\n"+`	v2 "github.com/crossplane/crossplane-runtime/v2/apis/common/v2"`+"\n",
				1,
			)
		}
		return os.WriteFile(path, []byte(content), 0o600)
	})
}
