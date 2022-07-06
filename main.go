package main

import (
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

	content   map[string]interface{}
	variables [][]string

	info       string = "%s (%s)\n\n%s\n%s\n\n"
	usage      string = "Usage of %s:\n"
	docker     string = "%s=%s\n"
	compose    string = "\"%s\": \"%s\"\n"
	kubernetes string = "- name: \"%s\"\n  value: \"%s\"\n"
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

	err = json.Unmarshal(file_bytes, &content)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	parser(content, nil)

	sort.Slice(variables[:], func(i, j int) bool {
		for key := range variables[i] {
			if variables[i][key] == variables[j][key] {
				continue
			}
			return variables[i][key] < variables[j][key]
		}
		return false
	})

	for _, value := range variables {
		switch output {
		case "docker":
			fmt.Printf(docker, value[0], value[1])
		case "compose":
			fmt.Printf(compose, value[0], value[1])
		default:
			fmt.Printf(kubernetes, value[0], value[1])
		}
	}
}

func parser(data map[string]interface{}, root []string) {
	for key, value := range data {
		keys := append(root, key)

		switch value.(type) {
		case []interface{}:
			for key, value := range value.([]interface{}) {
				switch value.(type) {
				case map[string]interface{}:
					parser(value.(map[string]interface{}), append(keys, fmt.Sprint(key)))
				default:
					variables = append(variables, []string{
						fmt.Sprintf("%s__%d", strings.Join(keys, separator), key),
						fmt.Sprint(value),
					})
				}

			}
		case map[string]interface{}:
			parser(value.(map[string]interface{}), keys)
		default:
			variables = append(variables, []string{
				strings.Join(keys, separator),
				fmt.Sprint(value),
			})
		}
	}
}
