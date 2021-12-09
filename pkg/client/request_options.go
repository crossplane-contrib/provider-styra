package client

import (
	"github.com/go-openapi/runtime"
)

// ForAllStatusCodes is a wildcard for all status codes
const ForAllStatusCodes = -1

// Option can be passed to a client call to modify the request
type Option func(*runtime.ClientOperation)

// DropContentTypeHeader by setting ConsumesMediaTypes to nil because some Styra APIs (i.e. DELETE /v1/secrets/<id>) will fail
// if a request contains a header Content-Type: application/json but does not provide body by default.
// NOTE: This effectively sets the header value to runtime.DefaultMediaType.
func DropContentTypeHeader(rt *runtime.ClientOperation) {
	rt.ConsumesMediaTypes = nil
}

// ReturnRawResponse overwrites the default consumer to return the raw response body in bytes
func ReturnRawResponse(rt *runtime.ClientOperation) {
	OverwriteConsumer(runtime.ByteStreamConsumer())(rt)
}

// OverwriteConsumer replaces the default consumer with another
func OverwriteConsumer(overwrite runtime.Consumer) Option {
	return OverwriteConsumerForStatus(overwrite, ForAllStatusCodes)
}

// OverwriteConsumerForStatus overwrites the default consumer with another for a given status code
func OverwriteConsumerForStatus(overwrite runtime.Consumer, forStatusCode int) Option {
	return func(rt *runtime.ClientOperation) {
		rt.Reader = &overwriteConsumerReader{
			requestReader: rt.Reader,
			consumer:      overwrite,
			forStatusCode: forStatusCode,
		}
	}
}

type overwriteConsumerReader struct {
	// the original request reader
	requestReader runtime.ClientResponseReader

	// the consumer that should be used
	consumer runtime.Consumer

	// the response status code for which the consumer should be overwritten
	forStatusCode int
}

// ReadResponse uses the given consumer instead of the default
func (r *overwriteConsumerReader) ReadResponse(resp runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	if r.forStatusCode == ForAllStatusCodes || resp.Code() == r.forStatusCode {
		return r.requestReader.ReadResponse(resp, r.consumer)
	}

	return r.requestReader.ReadResponse(resp, consumer)
}

// WrapperConsumer wraps the original consumer with another
type WrapperConsumer func(original runtime.Consumer) runtime.Consumer

// WrapConsumerForStatusCode wraps the original consumer for another with a given status code
func WrapConsumerForStatusCode(wrapper WrapperConsumer, forStatusCode int) Option {
	return func(rt *runtime.ClientOperation) {
		rt.Reader = &wrapConsumerReader{
			requestReader: rt.Reader,
			wrapper:       wrapper,
			forStatusCode: forStatusCode,
		}
	}
}

type wrapConsumerReader struct {
	// the original request reader
	requestReader runtime.ClientResponseReader

	// the consumer that should be used
	wrapper WrapperConsumer

	// the response status code for which the consumer should be overwritten
	forStatusCode int
}

// ReadResponse wraps the default consumer with the given
func (r *wrapConsumerReader) ReadResponse(resp runtime.ClientResponse, original runtime.Consumer) (interface{}, error) {
	if r.forStatusCode == ForAllStatusCodes || resp.Code() == r.forStatusCode {
		wrapperConsumer := r.wrapper(original)
		return r.requestReader.ReadResponse(resp, wrapperConsumer)
	}

	return r.requestReader.ReadResponse(resp, original)
}
