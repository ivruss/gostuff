package errors

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	InternalErrorMsg    = "Internal Error"
	IncorrectPayloadMsg = "Incorrect Payload"
)

type ErrorResult struct {
	Code int
	Msg  string
	Err  error
}

func (err *ErrorResult) Error() string {
	return err.Err.Error()
}

func WrapInternalError(err error) error {
	return &ErrorResult{Code: 500, Msg: "Internal Error", Err: err}
}

func WrapNotFoundError(err error) error {
	return &ErrorResult{Code: 404, Msg: "Content Not Found", Err: err}
}

func WrapBadRequestError(err error) error {
	return &ErrorResult{Code: 400, Msg: "Incorrect Payload", Err: err}
}

func FromError(err error) (*ErrorResult, bool) {
	var result *ErrorResult
	if ok := errors.As(err, &result); !ok {
		return nil, false
	}
	return result, true
}

func FromErrorToGRPC(err error) error {
	var result *ErrorResult
	if ok := errors.As(err, &result); !ok {
		return status.Errorf(codes.Internal, InternalErrorMsg)
	}

	fmt.Println(err)
	switch result.Code {
	case 400:
		return status.Errorf(codes.InvalidArgument, "%s", result.Msg)
	case 404:
		return status.Errorf(codes.NotFound, "%s", result.Msg)
	default:
		return status.Errorf(codes.Internal, "%s", result.Msg)
	}
}
