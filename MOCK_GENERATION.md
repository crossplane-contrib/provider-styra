# Client Mock Generation

[go-mock](https://github.com/golang/mock) is used to generate mocks of the Styra client.

## Install

Follow the [installation instructions](https://github.com/golang/mock#installation) to get the latest version.

## Generate mocks

The following example shows how to generate mocks for the `systems` API:

    MOCK_API="systems"
    MOCK_INTERFACE="ClientService"
    
    mockgen -package $MOCK_API -destination pkg/client/mock/$MOCK_API/mock.go github.com/mistermx/styra-go-client/pkg/client/$MOCK_API $MOCK_INTERFACE
