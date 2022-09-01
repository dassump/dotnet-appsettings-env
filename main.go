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
	app         string = "dotnet-appsettings-env"
	version     string = "dev"
	description string = "Convert .NET appsettings.json file to Kubernetes, Docker, Docker-Compose and Bicep environment variables."
	site        string = "https://github.com/dassump/dotnet-appsettings-env"
	info        string = "%s (%s)\n\n%s\n%s\n\n"
	usage       string = "Usage of %s:\n"

	file       string
	file_name  string = "file"
	file_value string = "./appsettings.json"
	file_usage string = "Path to file appsettings.json"
	file_error string = "%s\n"

	output       string
	output_name  string = "type"
	output_value string = "k8s"
	output_usage string = "Output to Kubernetes (k8s) / Docker (docker) / Docker Compose (compose) / Bicep (bicep)"
	output_error string = "Type %s is not valid\n"

	separator       string
	separator_name  string = "separator"
	separator_value string = "__"
	separator_usage string = "Separator character"

	content              map[string]any
	content_comments     string = `(?m:\/\*[\s\S]*?\*\/|([^:]|^)\/\/.*$)`
	content_error        string = "%s in %s\n"
	content_syntax_error string = "%s in %s\n\n... line %d, column %d\n%s >>> %s <<< %s\n...\n"
	content_syntax_near  int64  = 60

	variables = map[string]string{}
	format    = map[string]string{
		"k8s":     "- name: %q\n  value: %q\n",
		"docker":  "%s=%q\n",
		"compose": "%s: %q\n",
		"bicep":   "{\nname: '%s'\nvalue: '%s'\n}\n",
	}
)

func init() {
	log.SetFlags(log.Lmsgprefix)
	log.SetPrefix("Error: ")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), info, app, version, description, site)
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&file, file_name, file_value, file_usage)
	flag.StringVar(&output, output_name, output_value, output_usage)
	flag.StringVar(&separator, separator_name, separator_value, separator_usage)

	flag.Parse()

	if _, ok := format[output]; !ok {
		log.Fatalf(output_error, output)
	}
}

func main() {
	file_bytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf(file_error, err)
	}

	_, err = transform.NewReader(
		bytes.NewReader(file_bytes),
		unicode.BOMOverride(unicode.UTF8.NewDecoder()),
	).Read(file_bytes)
	if err != nil {
		log.Fatalf(file_error, err)
	}

	file_bytes = regexp.MustCompile(content_comments).ReplaceAll(file_bytes, nil)

	decoder := json.NewDecoder(bytes.NewReader(file_bytes))
	decoder.UseNumber()

	if err := decoder.Decode(&content); err != nil {
		switch err := err.(type) {
		case *json.SyntaxError:
			new_line := []byte("\n")
			line := 1 + bytes.Count(file_bytes[:err.Offset], new_line)
			column := int(err.Offset) - bytes.LastIndex(file_bytes[:err.Offset], new_line) - len(new_line)

			near_before := err.Offset - content_syntax_near
			if err.Offset-content_syntax_near < 0 {
				near_before = 0
			}

			near_after := err.Offset + content_syntax_near
			if err.Offset+content_syntax_near > int64(len(file_bytes)) {
				near_after = int64(len(file_bytes))
			}

			log.Fatalf(
				content_syntax_error,
				err, file, line, column,
				file_bytes[near_before:err.Offset-1], file_bytes[err.Offset-1:err.Offset], file_bytes[err.Offset:near_after],
			)
		default:
			log.Fatalf(content_error, err, file)
		}
	}

	parser(content, variables, nil)

	keys := make([]string, 0, len(variables))
	for key := range variables {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
	})

	for _, key := range keys {
		fmt.Printf(format[output], key, variables[key])
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
