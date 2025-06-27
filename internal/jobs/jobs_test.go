package jobs

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"resty.dev/v3"
)

// --- Mocks ---

type mockClient struct {
	jobToReturn *Job
	errToReturn error
}

func (m *mockClient) JobRefresh(_ *resty.Request, _ *resty.Response) (*Job, error) {
	return m.jobToReturn, m.errToReturn
}

func (m *mockClient) JobParser(_ *resty.Response) (*Job, error) {
	return m.jobToReturn, m.errToReturn
}

func (m *mockClient) JobStatusParser(status string) (Status, error) {
	return Status(status), nil
}

// --- Tests ---

func TestNewJobOptions_Defaults(t *testing.T) {
	opts, err := NewJobOptions()
	assert.NoError(t, err)
	assert.Equal(t, 5*time.Minute, opts.Timeout)
	assert.Equal(t, 15*time.Second, opts.PollInterval)
}

func TestNewJobOptions_CustomValues(t *testing.T) {
	opts, err := NewJobOptions(
		WithCustomTimeout(2*time.Minute),
		WithCustomPollInterval(10*time.Second),
	)
	assert.NoError(t, err)
	assert.Equal(t, 2*time.Minute, opts.Timeout)
	assert.Equal(t, 10*time.Second, opts.PollInterval)
}

func TestWithCustomPollInterval_Invalid(t *testing.T) {
	_, err := NewJobOptions(WithCustomPollInterval(0))
	assert.Error(t, err)
}

func TestSetExtractorFunc(t *testing.T) {
	var called *bool
	extractor := func(_ *resty.Response) { *called = true }
	opts, err := NewJobOptions(SetExtractorFunc(extractor))
	assert.NoError(t, err)
	assert.NotNil(t, opts.extractorFunc)
}

func TestStatus_IsTerminated(t *testing.T) {
	assert.True(t, Success.IsTerminated())
	assert.True(t, Error.IsTerminated())
	assert.True(t, Aborted.IsTerminated())
	assert.False(t, Queued.IsTerminated())
	assert.False(t, Running.IsTerminated())
}

func TestNewJobMiddleware_JobCompletesSuccessfully(t *testing.T) {
	client := resty.New()
	job := &Job{ID: "1", Status: Success}
	mock := &mockClient{jobToReturn: job}

	opts, _ := NewJobOptions()
	middleware := NewJobMiddleware(client, mock, opts)

	resp := &resty.Response{
		Request: &resty.Request{}, // Create a new request for the response
	}
	err := middleware(client, resp)
	assert.NoError(t, err)
}

func TestNewJobMiddleware_JobFails(t *testing.T) {
	client := resty.New()
	mock := &mockClient{jobToReturn: nil, errToReturn: errors.New("job failed")}

	opts, _ := NewJobOptions()
	middleware := NewJobMiddleware(client, mock, opts)

	resp := &resty.Response{
		Request: &resty.Request{}, // Create a new request for the response
	}
	err := middleware(client, resp)
	assert.Error(t, err)
}

func TestNewJobMiddleware_NilOptions(t *testing.T) {
	client := resty.New()
	mock := &mockClient{}

	middleware := NewJobMiddleware(client, mock, nil)
	resp := &resty.Response{
		Request: &resty.Request{}, // Create a new request for the response
	}
	err := middleware(client, resp)
	assert.Error(t, err)
}

func TestJobRetryCondition_ErrorOnWait(t *testing.T) {
	mock := &mockClient{}
	resp := &resty.Response{
		Request: &resty.Request{},
	}
	cond := jobRetryCondition(mock)
	shouldRetry := cond(resp, errors.New("network error"))
	assert.False(t, shouldRetry)
}

func TestJobRetryCondition_ErrorOnParse(t *testing.T) {
	mock := &mockClient{errToReturn: errors.New("parse error")}
	resp := &resty.Response{
		Request: &resty.Request{},
	}
	cond := jobRetryCondition(mock)
	shouldRetry := cond(resp, nil)
	assert.False(t, shouldRetry)
}

func TestJobRetryCondition_JobNotTerminated(t *testing.T) {
	mock := &mockClient{jobToReturn: &Job{Status: Running}}
	resp := &resty.Response{
		Request: &resty.Request{},
	}
	cond := jobRetryCondition(mock)
	shouldRetry := cond(resp, nil)
	assert.True(t, shouldRetry)
}

func TestJobRetryCondition_JobTerminated(t *testing.T) {
	mock := &mockClient{jobToReturn: &Job{Status: Success}}
	resp := &resty.Response{
		Request: &resty.Request{},
	}
	cond := jobRetryCondition(mock)
	shouldRetry := cond(resp, nil)
	assert.False(t, shouldRetry)
}

func TestNewJobMiddleware_ExtractorFuncDefined(t *testing.T) {
	client := resty.New()

	opts, err := NewJobOptions(SetExtractorFunc(func(_ *resty.Response) {}))
	assert.NoError(t, err)

	middleware := NewJobMiddleware(client, &mockClient{}, opts)
	if err := middleware(client, &resty.Response{
		Request: &resty.Request{},
	}); err != nil {
		assert.NoError(t, err)
	}
}

func TestNewJobMiddleware_ExtractorFuncCalled(t *testing.T) {
	client := resty.New()

	var called *bool
	called = new(bool)

	extractor := func(_ *resty.Response) { *called = true }

	middleware := extractorFuncMiddleware(extractor)

	resp := &resty.Response{
		Request: &resty.Request{},
	}
	// Simulate a response middleware chain
	_ = middleware(client, resp)

	assert.True(t, *called, "extractorFunc should be called by the middleware")
}
