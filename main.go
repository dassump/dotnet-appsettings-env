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
	app, version string
	file         string
	output       string
	separator    string
	content      map[string]interface{}
	variables    [][]string
)

const (
	description  string = "Convert .NET appsettings.json file to Kubernetes, Docker and Docker-Compose environment variables."
	author_name  string = "Daniel Dias de Assumpção"
	author_email string = "dassump@gmail.com"
	site         string = "https://github.com/dassump"
)

func init() {
	flag.Usage = func() {
		fmt.Printf(
			"%s %s\n\n%s\n\n%s <%s>\n%s\n\n",
			app, version, description, author_name, author_email, site,
		)

		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&file, "file", "./appsettings.json", "Path to file appsettings.json")
	flag.StringVar(&output, "type", "k8s", "Output to Kubernetes (k8s) / Docker (docker) / Docker Compose (compose)")
	flag.StringVar(&separator, "separator", "__", "Separator character")

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
			fmt.Printf("%s=%s\n", value[0], value[1])
		case "compose":
			fmt.Printf("\"%s\": \"%s\"\n", value[0], value[1])
		default:
			fmt.Printf("- name: \"%s\"\n", value[0])
			fmt.Printf("  value: \"%s\"\n", value[1])
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
