# OpenAPI Codegen

Most OpenAPI tools can work with multi-file OpenAPI definitions and resolve $refs dynamically.

If you specifically need to get a single resolved file, Swagger Codegen can do this. Codegen has a CLI version (used in
the examples below), a Maven plugin (usage example) and a Docker image.

The input file (-i argument of the CLI) can be a local file or a URL.

Note: Line breaks are added for readability.

OpenAPI 3.0 example
Use Codegen 3.x to resolve OpenAPI 3.0 files:

```
java -jar swagger-codegen-cli-3.0.35.jar generate
     -l openapi-yaml
     -i ./path/to/openapi.yaml
     -o ./OUT_DIR
     -DoutputFile=output.yaml
```

-l openapi-yaml outputs YAML, -l openapi outputs JSON.

-DoutputFile is optional, the default file name is openapi.yaml / openapi.json.

OpenAPI 2.0 example
Use Codegen 2.x to resolve OpenAPI 2.0 files (swagger: '2.0'):

```
java -jar swagger-codegen-cli-2.4.28.jar generate
     -l swagger-yaml
     -i ./path/to/openapi.yaml
     -o ./OUT_DIR
     -DoutputFile=output.yaml
```

-l swagger-yaml outputs YAML, -l swagger outputs JSON.

-DoutputFile is optional, the default file name is swagger.yaml / swagger.json.