package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var (
	app         = "dotnet-appsettings-env"
	version     = "dev"
	description = "Convert .NET appsettings.json file to Kubernetes, Docker, Docker-Compose and Bicep environment variables."
	site        = "https://github.com/dassump/dotnet-appsettings-env"

	file      = flag.String("file", "./appsettings.json", "Path to file appsettings.json")
	output    = flag.String("type", "k8s", "Output to Kubernetes (k8s) / Docker (docker) / Docker Compose (compose) / Bicep (bicep)")
	separator = flag.String("separator", "__", "Separator character(s)")

	comments = `(?m:\/\*[\s\S]*?\*\/|([^:]|^)\/\/.*$)`
	format   = map[string]string{
		"k8s":     "- name: %q\n  value: %q\n",
		"docker":  "%s=%q\n",
		"compose": "%s: %q\n",
		"bicep":   "{\nname: '%s'\nvalue: '%s'\n}\n",
	}

	objects   = make(map[string]any)
	variables = make(map[string]string)
)

func init() {
	log.SetFlags(log.Lmsgprefix)
	log.SetPrefix("Error: ")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s (%s)\n\n%s\n%s\n\n", app, version, description, site)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if _, ok := format[*output]; !ok {
		log.Fatalln(*output, "is not valid type")
	}

	if len(*separator) < 1 {
		log.Fatalln("separator cannot be an empty string")
	}
}

func main() {
	content, err := os.ReadFile(*file)
	if err != nil {
		log.Fatalln(err)
	}

	if content, _, err = transform.Bytes(
		unicode.BOMOverride(unicode.UTF8.NewDecoder()), content,
	); err != nil {
		log.Fatalln(err)
	}

	content = regexp.MustCompile(comments).ReplaceAll(content, nil)

	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.UseNumber()

	if err := decoder.Decode(&objects); err != nil {
		switch err := err.(type) {
		case *json.SyntaxError:
			newline := []byte("\n")
			line := 1 + bytes.Count(content[:err.Offset], newline)
			column := int(err.Offset) - bytes.LastIndex(content[:err.Offset], newline) - len(newline)
			near := int64(60)

			before := err.Offset - near
			if err.Offset-near < 0 {
				before = 0
			}

			after := err.Offset + near
			if err.Offset+near > int64(len(content)) {
				after = int64(len(content))
			}

			log.Fatalf(
				"%s in %s\n\n... line %d, column %d\n%s >>> %s <<< %s\n...\n",
				err, *file, line, column,
				content[before:err.Offset-1], content[err.Offset-1:err.Offset], content[err.Offset:after],
			)

		default:
			log.Fatalf("%s in %s\n", err, *file)
		}
	}

	parser(objects, variables, nil)

	keys := make([]string, 0, len(variables))
	for key := range variables {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
	})

	for _, key := range keys {
		fmt.Printf(format[*output], key, variables[key])
	}
}

func parser(in map[string]any, out map[string]string, root []string) {
	for key, value := range in {
		keys := append(root, key)

		switch any(value).(type) {
		case []any:
			for key, value := range value.([]any) {
				switch any(value).(type) {
				case map[string]any:
					parser(value.(map[string]any), out, append(keys, fmt.Sprint(key)))

				default:
					out[fmt.Sprintf("%s%s%d", strings.Join(keys, *separator), *separator, key)] = fmt.Sprint(value)
				}
			}

		case map[string]any:
			parser(value.(map[string]any), out, keys)

		default:
			out[strings.Join(keys, *separator)] = fmt.Sprint(value)
		}
	}
}
