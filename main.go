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
	vars         [][]string
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
	b, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = json.Unmarshal(b, &content)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	parser(content, nil)

	sort.Slice(vars[:], func(i, j int) bool {
		for x := range vars[i] {
			if vars[i][x] == vars[j][x] {
				continue
			}
			return vars[i][x] < vars[j][x]
		}
		return false
	})

	for _, v := range vars {
		switch output {
		case "docker":
			fmt.Printf("%s=%s\n", v[0], v[1])
		case "compose":
			fmt.Printf("\"%s\": \"%s\"\n", v[0], v[1])
		default:
			fmt.Printf("- name: \"%s\"\n", v[0])
			fmt.Printf("  value: \"%s\"\n", v[1])
		}
	}
}

func parser(m map[string]interface{}, root []string) {
	for k, v := range m {
		key := append(root, k)

		switch v.(type) {
		case []interface{}:
			for k, v := range v.([]interface{}) {
				switch v.(type) {
				case map[string]interface{}:
					parser(v.(map[string]interface{}), append(key, fmt.Sprint(k)))
				default:
					vars = append(vars, []string{
						fmt.Sprintf("%s__%d", strings.Join(key, separator), k),
						fmt.Sprint(v),
					})
				}

			}
		case map[string]interface{}:
			parser(v.(map[string]interface{}), key)
		default:
			vars = append(vars, []string{
				strings.Join(key, separator),
				fmt.Sprint(v),
			})
		}
	}
}
