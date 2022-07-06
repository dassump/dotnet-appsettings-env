# dotnet-appsettings-env

Convert .NET appsettings.json file to Kubernetes, Docker and Docker-Compose environment variables.


## Getting started

1. Download a pre-compiled binary from the release [page](https://github.com/dassump/dotnet-appsettings-env/releases).
2. Run `dotnet-appsettings-env --help`

```shell
$ dotnet-appsettings-env --help
dotnet-appsettings-env (dev)

Convert .NetCore appsettings.json file to Kubernetes, Docker and Docker-Compose environment variables.
https://github.com/dassump/dotnet-appsettings-env

Usage of dotnet-appsettings-env:
  -file string
        Path to file appsettings.json (default "./appsettings.json")
  -separator string
        Separator character (default "__")
  -type string
        Output to Kubernetes (k8s) / Docker (docker) / Docker Compose (compose) (default "k8s")
```


## Examples

### appsettings.json

```json
{
    "ApiClientId": "*",
    "ApiClientSecret": "*",
    "ApiGateway": "*",
    "Scope": "*",
    "Middlewares": [
        {
            "Name": "api/Auth",
            "Url": "*"
        },
        {
            "Name": "api/Registration",
            "Url": "*"
        }
    ],
    "HttpManager": {
        "IgnoreCertificateValidation": true,
        "AllowAutoRedirect": true
    },
    "Logging": {
        "Enabled": true,
        "IncludeScopes": false,
        "Level": "Information",
        "Debug": {
            "LogLevel": {
                "Default": "Warning"
            }
        },
        "Console": {
            "LogLevel": {
                "Default": "Warning"
            }
        }
    },
    "Serilog": {
        "Using": [
            "Serilog.Sinks.File"
        ],
        "MinimumLevel": "Debug",
        "WriteTo": [
            {
                "Name": "Console"
            },
            {
                "Name": "File",
                "Args": {
                    "path": "Logs/Api.log",
                    "rollingInterval": "Day",
                    "fileSizeLimitBytes": "52428800",
                    "rollOnFileSizeLimit": "true",
                    "retainedFileCountLimit": "100",
                    "outputTemplate": "{Timestamp:yyyy-MM-dd HH:mm:ss.fff zzz} [{Level:u3}] {Message:lj}{NewLine}{Exception}"
                }
            }
        ]
    }
}
```


### Kubernetes

```shell
$ dotnet-appsettings-env -type k8s
- name: "ApiClientId"
  value: "*"
- name: "ApiClientSecret"
  value: "*"
- name: "ApiGateway"
  value: "*"
- name: "HttpManager__AllowAutoRedirect"
  value: "true"
- name: "HttpManager__IgnoreCertificateValidation"
  value: "true"
- name: "Logging__Console__LogLevel__Default"
  value: "Warning"
- name: "Logging__Debug__LogLevel__Default"
  value: "Warning"
- name: "Logging__Enabled"
  value: "true"
- name: "Logging__IncludeScopes"
  value: "false"
- name: "Logging__Level"
  value: "Information"
- name: "Middlewares__0__Name"
  value: "api/Auth"
- name: "Middlewares__0__Url"
  value: "*"
- name: "Middlewares__1__Name"
  value: "api/Registration"
- name: "Middlewares__1__Url"
  value: "*"
- name: "Scope"
  value: "*"
- name: "Serilog__MinimumLevel"
  value: "Debug"
- name: "Serilog__Using__0"
  value: "Serilog.Sinks.File"
- name: "Serilog__WriteTo__0__Name"
  value: "Console"
- name: "Serilog__WriteTo__1__Args__fileSizeLimitBytes"
  value: "52428800"
- name: "Serilog__WriteTo__1__Args__outputTemplate"
  value: "{Timestamp:yyyy-MM-dd HH:mm:ss.fff zzz} [{Level:u3}] {Message:lj}{NewLine}{Exception}"
- name: "Serilog__WriteTo__1__Args__path"
  value: "Logs/Api.log"
- name: "Serilog__WriteTo__1__Args__retainedFileCountLimit"
  value: "100"
- name: "Serilog__WriteTo__1__Args__rollingInterval"
  value: "Day"
- name: "Serilog__WriteTo__1__Args__rollOnFileSizeLimit"
  value: "true"
- name: "Serilog__WriteTo__1__Name"
  value: "File"
```


### Docker

```shell
$ dotnet-appsettings-env -type docker
ApiClientId="*"
ApiClientSecret="*"
ApiGateway="*"
HttpManager__AllowAutoRedirect="true"
HttpManager__IgnoreCertificateValidation="true"
Logging__Console__LogLevel__Default="Warning"
Logging__Debug__LogLevel__Default="Warning"
Logging__Enabled="true"
Logging__IncludeScopes="false"
Logging__Level="Information"
Middlewares__0__Name="api/Auth"
Middlewares__0__Url="*"
Middlewares__1__Name="api/Registration"
Middlewares__1__Url="*"
Scope="*"
Serilog__MinimumLevel="Debug"
Serilog__Using__0="Serilog.Sinks.File"
Serilog__WriteTo__0__Name="Console"
Serilog__WriteTo__1__Args__fileSizeLimitBytes="52428800"
Serilog__WriteTo__1__Args__outputTemplate="{Timestamp:yyyy-MM-dd HH:mm:ss.fff zzz} [{Level:u3}] {Message:lj}{NewLine}{Exception}"
Serilog__WriteTo__1__Args__path="Logs/Api.log"
Serilog__WriteTo__1__Args__retainedFileCountLimit="100"
Serilog__WriteTo__1__Args__rollingInterval="Day"
Serilog__WriteTo__1__Args__rollOnFileSizeLimit="true"
Serilog__WriteTo__1__Name="File"
```


### Docker Compose

```shell
$ dotnet-appsettings-env -type compose
ApiClientId: "*"
ApiClientSecret: "*"
ApiGateway: "*"
HttpManager__AllowAutoRedirect: "true"
HttpManager__IgnoreCertificateValidation: "true"
Logging__Console__LogLevel__Default: "Warning"
Logging__Debug__LogLevel__Default: "Warning"
Logging__Enabled: "true"
Logging__IncludeScopes: "false"
Logging__Level: "Information"
Middlewares__0__Name: "api/Auth"
Middlewares__0__Url: "*"
Middlewares__1__Name: "api/Registration"
Middlewares__1__Url: "*"
Scope: "*"
Serilog__MinimumLevel: "Debug"
Serilog__Using__0: "Serilog.Sinks.File"
Serilog__WriteTo__0__Name: "Console"
Serilog__WriteTo__1__Args__fileSizeLimitBytes: "52428800"
Serilog__WriteTo__1__Args__outputTemplate: "{Timestamp:yyyy-MM-dd HH:mm:ss.fff zzz} [{Level:u3}] {Message:lj}{NewLine}{Exception}"
Serilog__WriteTo__1__Args__path: "Logs/Api.log"
Serilog__WriteTo__1__Args__retainedFileCountLimit: "100"
Serilog__WriteTo__1__Args__rollingInterval: "Day"
Serilog__WriteTo__1__Args__rollOnFileSizeLimit: "true"
Serilog__WriteTo__1__Name: "File"
```


## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/dassump/dotnet-appsettings-env.