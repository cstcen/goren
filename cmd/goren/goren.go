package main

import (
	"bytes"
	"flag"
	"fmt"
	"git.tenvine.cn/backend/goren/v2/internal/gen"
	"os"
	"os/exec"
	"path"
)

func errExit(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

var (
	flagOutputFile string
	flagConfigFile string
	flagPackage    string
)

func main() {
	//flag.StringVar(&flagOutputFile, "o", "goren/api/gen/api.go", "Where to output generated code, stdout is default")
	//flag.StringVar(&flagConfigFile, "config", "goren/goren-config.yaml", "a YAML config file that controls oapi-codegen behavior")
	flag.StringVar(&flagPackage, "package", "api", "package name")

	flag.Parse()

	genOutputDir := fmt.Sprintf("goren/%s/gen", flagPackage)
	flagOutputFile = fmt.Sprintf("%s/%s.gen.go", genOutputDir, flagPackage)

	//if output, err := exec.Command("go", "get", "-u", "github.com/deepmap/oapi-codegen/cmd/oapi-codegen/v2").CombinedOutput(); err != nil {
	//	errExit("error go get cmd/oapi-codegen, output: %s\n", output)
	//}

	if output, err := exec.Command("go", "install", "github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest").CombinedOutput(); err != nil {
		errExit("error go install cmd/oapi-codegen, output: %s\n", output)
	}

	var moduleName string
	if output, err := exec.Command("go", "list", "-m").CombinedOutput(); err != nil {
		errExit("error go list -m, output: %s\n", output)
	} else {
		moduleName = string(bytes.TrimSpace(output))
	}

	dir, _ := path.Split(flagOutputFile)
	if _, err := os.Stat(flagOutputFile); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			errExit("error mkdir all: %s\n", err.Error())
		}
	}
	packageNameApi := flagPackage
	packageNameApiGen := packageNameApi + "gen"
	ops := gen.Params{
		ModuleName:        moduleName,
		PackageNameApi:    packageNameApi,
		PackageNameApiGen: packageNameApiGen,
	}
	if err := gen.Generate(flagOutputFile, ops); err != nil {
		errExit("error generate main: %s\n", err)
	}

	filenameSchemas := "https://doc.xk5.com/specs/externalref/schemas.yaml"
	if output, err := exec.Command("oapi-codegen", "-config", path.Join(genOutputDir, "schemas", "goren-config-schemas.yaml"), filenameSchemas).CombinedOutput(); err != nil {
		errExit("failed to exec cmd oapi-codegen, name: %s, output: %s\n", filenameSchemas, output)
	}
	filenameParameters := "https://doc.xk5.com/specs/externalref/parameters.yaml"
	if output, err := exec.Command("oapi-codegen", "-config", path.Join(genOutputDir, "parameters", "goren-config-parameters.yaml"), filenameParameters).CombinedOutput(); err != nil {
		errExit("failed to exec cmd oapi-codegen, name: %s, output: %s\n", filenameParameters, output)
	}
	filenameResponses := "https://doc.xk5.com/specs/externalref/responses.yaml"
	if output, err := exec.Command("oapi-codegen", "-config", path.Join(genOutputDir, "responses", "goren-config-responses.yaml"), filenameResponses).CombinedOutput(); err != nil {
		errExit("failed to exec cmd oapi-codegen, name: %s, output: %s\n", filenameResponses, output)
	}

	filenameOpenapi := flag.Arg(0)
	if output, err := exec.Command("oapi-codegen", "-config", path.Join(genOutputDir, "goren-config.yaml"), "-o", flagOutputFile, "-package", packageNameApiGen, filenameOpenapi).CombinedOutput(); err != nil {
		errExit("failed to exec cmd oapi-codegen, name: %s, output: %s\n", filenameOpenapi, output)
	}
}
