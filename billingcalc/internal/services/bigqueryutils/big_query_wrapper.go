package bigqueryutils

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type BigQueryWrapper interface {
	ExecuteQuery(ctx context.Context, query string) (ResultIterator, error)
	Close() error
}

type bigQueryWrapper struct {
	client *bigquery.Client
}

// ResultIterator iterates over BigQuery query results
// This interface exists primarily for testing purposes
type ResultIterator interface {
	Next(dst interface{}) error
}

// RowIterator implements ResultIterator and wraps bigquery.RowIterator
type RowIterator struct {
	iterator *bigquery.RowIterator
}

// NewRowIterator returns a new RowIterator
func NewRowIterator(itr *bigquery.RowIterator) ResultIterator {
	return &RowIterator{itr}
}

// Next is a proxy for bigquery.RowIterator.Next
func (r *RowIterator) Next(dst interface{}) error {
	err := r.iterator.Next(dst)
	if errors.Is(err, iterator.Done) {
		return IteratorDoneError{}
	}
	return err //nolint:wrapcheck
}

// IteratorDoneError is an error that indicates that a ResultIterator is done
// This is returned whenever ResultIterator.Next is called with no more results
type IteratorDoneError struct{}

// Error implements the Error interface, the message does not need to be useful
func (r IteratorDoneError) Error() string {
	return "Done iterating."
}

// NewBigQueryWrapper creates a new BigQueryWrapper
func NewBigQueryWrapper(ctx context.Context, projectID string) (BigQueryWrapper, error) {
	if len(projectID) == 0 {
		return nil, errors.New("Failed to create a new BigQuery Client. The projectID must be set.")
	}

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		msg := fmt.Sprintf("Failed to create a new BigQuery Client: \n %v \n", err)
		return nil, errors.New(msg)
	}

	return &bigQueryWrapper{client}, nil
}

func (b *bigQueryWrapper) ExecuteQuery(ctx context.Context, query string) (ResultIterator, error) {
	q := b.client.Query(query)
	iterator, err := q.Read(ctx)
	if err != nil {
		msg := fmt.Sprintf("Failed to submit a query for execution: %v", err)
		return nil, errors.New(msg)
	}

	return iterator, nil
}

func (b *bigQueryWrapper) Close() error {
	err := b.client.Close()
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
