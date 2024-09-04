package feed

type FeedHTTPError struct {
	statusCode int
	message    string
}

func (f FeedHTTPError) StatusCode() int {
	return f.statusCode
}

func (f FeedHTTPError) Error() string {
	return f.message
}

func NewFeedHTTPError(statusCode int, message string) FeedHTTPError {
	return FeedHTTPError{statusCode, message}
}

var _ error = NewFeedHTTPError(0, "")
