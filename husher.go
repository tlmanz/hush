package hush

import (
	"context"
	"reflect"
	"sort"
	"sync"
)

// Husher is the interface that wraps the Hush method.
type Husher interface {
	Hush(ctx context.Context, prefix string, options ...Option) ([][]string, error)
}

type hushType struct {
	value reflect.Value
	isStr bool
}

// Hush processes the struct or string and returns a slice of field name-value pairs.
func (ht *hushType) Hush(ctx context.Context, prefix string, options ...Option) ([][]string, error) {
	opts := &hushOptions{
		separator:      DefaultSeparator,
		maskFunc:       defaultMaskFunc,
		includePrivate: false,
	}
	for _, o := range options {
		o(opts)
	}

	if ht.isStr {
		return [][]string{{prefix, opts.maskFunc(ht.value.String())}}, nil
	}

	// Pre-allocate the data slice with an estimated capacity
	data := make([][]string, 0, ht.value.NumField())

	// Use a buffered channel for error reporting
	errChan := make(chan error, ht.value.NumField())

	var wg sync.WaitGroup
	var mu sync.Mutex

	t := ht.value.Type()

	for i := 0; i < ht.value.NumField(); i++ {
		field := t.Field(i)
		value := ht.value.Field(i)

		if !opts.includePrivate && !field.IsExported() {
			continue
		}

		wg.Add(1)
		go func(field reflect.StructField, value reflect.Value) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
			}

			fieldName := buildFieldName(prefix, field.Name, opts.separator)

			result, err := processField(ctx, fieldName, field, value, opts)
			if err != nil {
				errChan <- err
				return
			}

			mu.Lock()
			data = append(data, result...)
			mu.Unlock()
		}(field, value)
	}

	wg.Wait()
	close(errChan)

	// Collect errors (only the first one if you prefer)
	var firstErr error
	for err := range errChan {
		firstErr = err
		break // Remove this line if you want to collect all errors
	}

	if firstErr != nil {
		return nil, firstErr
	}

	// Sort the result to ensure consistent order
	sort.Slice(data, func(i, j int) bool {
		return data[i][0] < data[j][0]
	})

	return data, nil
}