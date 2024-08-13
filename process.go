package hush

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"sync"
)

// processField is the main function for processing individual fields.
// It handles different types of fields (struct, pointer, slice, map, etc.) and applies the appropriate masking.
func (ht *hushType) processValue(ctx context.Context, fieldName string, field reflect.StructField, value reflect.Value, opts *hushOptions) ([][]string, error) {
	hushTag := field.Tag.Get("hush")

	if field.PkgPath != "" && !opts.includePrivate {
		return nil, nil // Skip unexported fields when not including private fields
	}

	switch value.Kind() {
	case reflect.Struct:
		return ht.processStruct(ctx, value, fieldName, opts)
	case reflect.Ptr:
		if value.IsNil() {
			return [][]string{{fieldName, "nil"}}, nil
		}
		return ht.processValue(ctx, fieldName, reflect.StructField{}, value.Elem(), opts)
	case reflect.Slice, reflect.Array:
		return ht.processSliceOrArray(ctx, fieldName, value, opts)
	case reflect.Map:
		return ht.processMap(ctx, fieldName, value, opts)
	default:
		return ht.processSimpleField(fieldName, field, value, hushTag, opts)
	}
}

// processStruct handles the processing of struct fields.
func (ht *hushType) processStruct(ctx context.Context, rv reflect.Value, prefix string, opts *hushOptions) ([][]string, error) {
	data := make([][]string, 0, rv.NumField())
	errChan := make(chan error, rv.NumField())
	var wg sync.WaitGroup
	var mu sync.Mutex

	t := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := t.Field(i)
		value := rv.Field(i)

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

			result, err := ht.processValue(ctx, fieldName, field, value, opts)
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

	var firstErr error
	for err := range errChan {
		firstErr = err
		break
	}

	if firstErr != nil {
		return nil, firstErr
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i][0] < data[j][0]
	})

	return data, nil
}

// processSimpleField handles the processing of simple (non-composite) fields.
func (ht *hushType) processSimpleField(fieldName string, field reflect.StructField, value reflect.Value, hushTag string, opts *hushOptions) ([][]string, error) {
	if field.PkgPath != "" {
		// This is an unexported field
		if opts.includePrivate {
			return processNonComposite(fieldName, value, hushTag, opts)
		}
		return nil, nil // Skip unexported fields when not including private fields
	}

	// For exported fields, use Interface() as before
	return processNonComposite(fieldName, value, hushTag, opts)
}

func processNonComposite(fieldName string, value reflect.Value, hushTag string, opts *hushOptions) ([][]string, error) {
	convertedString := convertNonCompositeToString(value)
	return processString(fieldName, convertedString, hushTag, opts), nil
}

// processString applies the masking function to string values if needed.
func processString(fieldName, value, hushTag string, opts *hushOptions) [][]string {
	if opts.hushType != "" {
		hushTag = string(opts.hushType)
	}

	if hushTag == string(TagHide) {
		value = HiddenValue
	} else if hushTag == string(TagMask) && opts.maskFunc != nil {
		value = opts.maskFunc(value)
	}

	if fieldName == "" {
		return [][]string{{value}}
	}

	return [][]string{{fieldName, value}}
}

func convertNonCompositeToString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%v", value)
	case reflect.String:
		return value.String()
	default:
		return ""
	}
}

// processSliceOrArray handles the processing of slice or array fields.
func (ht *hushType) processSliceOrArray(ctx context.Context, fieldName string, value reflect.Value, opts *hushOptions) ([][]string, error) {
	var result [][]string
	for i := 0; i < value.Len(); i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		elemFieldName := fmt.Sprintf("%s[%d]", fieldName, i)
		elemValue := value.Index(i)

		elemResult, err := ht.processValue(ctx, elemFieldName, reflect.StructField{}, elemValue, opts)
		if err != nil {
			return nil, err
		}
		result = append(result, elemResult...)
	}
	return result, nil
}

// processMap handles the processing of map fields.
func (ht *hushType) processMap(ctx context.Context, fieldName string, value reflect.Value, opts *hushOptions) ([][]string, error) {
	result := make([][]string, 0, value.Len())
	for _, key := range value.MapKeys() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		keyStr := fmt.Sprintf("%v", key.Interface())
		mapFieldName := fieldName + "[" + keyStr + "]"
		fieldValue := value.MapIndex(key)

		processedValue, err := ht.processValue(ctx, mapFieldName, reflect.StructField{}, fieldValue, opts)
		if err != nil {
			return nil, err
		}
		result = append(result, processedValue...)
	}

	// Sort the result to ensure consistent order
	sort.Slice(result, func(i, j int) bool {
		return result[i][0] < result[j][0]
	})

	return result, nil
}
