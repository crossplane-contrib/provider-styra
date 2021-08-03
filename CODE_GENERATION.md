# Code Generation

[go-swagger](https://github.com/go-swagger/go-swagger) is used in this repo to generate the client code needed to connect to the Styra API.

To generate the code [install go-swagger](https://goswagger.io/install.html) and run:

```bash
swagger generate client -f doc/styra-swagger.json --skip-validation -c pkg/client -m pkg/client/models --default-scheme=https
```

The generated code might contain some errors regarding integer to string conversion. This is not an issue with swagger but rather the OpenAPI spec being malformed. In this case it should be checked that all integer types (especially query parameters) haven an appropriate type definition:

```yaml 
    type: "integer"
    format: "int32"
```
