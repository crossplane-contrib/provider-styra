# Client Mock Generation

[go-mock](https://github.com/golang/mock) is used to generate mocks of the Styra client.

## Generate mocks

Mocks are generated during `make generate`. Generator calls should be added to [pkg/client/mock/generate.go](pkg/client/mock/generate.go).

Mock generation for Styra should look like this:

```golang
//go:generate go run -tags generate github.com/golang/mock/mockgen -package MOCK_API -destination ./systems/mock.go github.com/mistermx/styra-go-client/pkg/client/MOCK_API ClientService
```
