package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

var (
	app         string = "dotnet-appsettings-env"
	version     string = "dev"
	description string = "Convert .NET appsettings.json file to Kubernetes, Docker and Docker-Compose environment variables."
	site        string = "https://github.com/dassump/dotnet-appsettings-env"

	file       string
	file_name  string = "file"
	file_value string = "./appsettings.json"
	file_usage string = "Path to file appsettings.json"

	output       string
	output_name  string = "type"
	output_value string = "k8s"
	output_usage string = "Output to Kubernetes (k8s) / Docker (docker) / Docker Compose (compose)"

	separator       string
	separator_name  string = "separator"
	separator_value string = "__"
	separator_usage string = "Separator character"

	info       string = "%s (%s)\n\n%s\n%s\n\n"
	usage      string = "Usage of %s:\n"
	docker     string = "%s=%s\n"
	compose    string = "\"%s\": \"%s\"\n"
	kubernetes string = "- name: \"%s\"\n  value: \"%s\"\n"

	content   = map[string]any{}
	variables = map[string]string{}
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), info, app, version, description, site)
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&file, file_name, file_value, file_usage)
	flag.StringVar(&output, output_name, output_value, output_usage)
	flag.StringVar(&separator, separator_name, separator_value, separator_usage)

	flag.Parse()
}

func main() {
	file_bytes, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	decoder := json.NewDecoder(bytes.NewReader(file_bytes))
	decoder.UseNumber()

	if err := decoder.Decode(&content); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	parser(content, variables, nil)

	keys := make([]string, 0, len(variables))
	for key := range variables {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		switch output {
		case "docker":
			fmt.Printf(docker, key, variables[key])
		case "compose":
			fmt.Printf(compose, key, variables[key])
		default:
			fmt.Printf(kubernetes, key, variables[key])
		}
	}
}

func parser(in map[string]any, out map[string]string, root []string) {
	for key, value := range in {
		keys := append(root, key)

		switch value.(type) {
		case []any:
			for key, value := range value.([]any) {
				switch value.(type) {
				case map[string]any:
					parser(value.(map[string]any), out, append(keys, fmt.Sprint(key)))
				default:
					out[fmt.Sprintf("%s__%d", strings.Join(keys, separator), key)] = fmt.Sprint(value)
				}

			}
		case map[string]any:
			parser(value.(map[string]any), out, keys)
		default:
			out[strings.Join(keys, separator)] = fmt.Sprint(value)
		}
	}
}
