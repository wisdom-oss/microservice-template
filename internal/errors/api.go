package errors

// This file contains all errors that are generated by the API during the
// handling of requests.
// The errors are then imported and sent back by using the .Emit
//
//
// Example:
//
//   var ErrExample = types.ServiceError{
//   	Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#section-15.5.5",
//		Status: http.StatusNotFound,
//		Title:  "Route Not Found",
//		Detail: "The requested path does not exist in this microservice. Please check the documentation and your request",
//   }