package hush

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"sync"
)

const (
	// maxRecursionDepth guards against stack overflow from circular pointer references.
	maxRecursionDepth = 64

	// maxConcurrency limits the number of goroutines spawned for concurrent struct field processing.
	maxConcurrency = 64
)

// processValue is the main function for processing individual fields.
// It handles different types of fields (struct, pointer, slice, map, etc.) and applies the appropriate masking.
// inheritedTag carries the hush tag from a parent field (e.g., through pointer dereference) so it isn't lost.
func (ht *hushType) processValue(ctx context.Context, fieldName string, field reflect.StructField, value reflect.Value, opts *hushOptions, inheritedTag string, depth int) ([][]string, error) {
	if depth > maxRecursionDepth {
		return [][]string{{fieldName, "[max depth exceeded]"}}, nil
	}

	hushTag := field.Tag.Get("hush")
	if hushTag == "" {
		hushTag = inheritedTag
	}

	if field.PkgPath != "" && !opts.includePrivate {
		return nil, nil // Skip unexported fields when not including private fields
	}

	if hushTag == string(TagRemove) {
		return nil, nil
	}

	switch value.Kind() {
	case reflect.Struct:
		return ht.processStruct(ctx, value, fieldName, opts, depth+1)
	case reflect.Ptr:
		if value.IsNil() {
			return [][]string{{fieldName, "nil"}}, nil
		}
		return ht.processValue(ctx, fieldName, reflect.StructField{}, value.Elem(), opts, hushTag, depth+1)
	case reflect.Slice, reflect.Array:
		return ht.processSliceOrArray(ctx, fieldName, value, opts, hushTag, depth+1)
	case reflect.Map:
		return ht.processMap(ctx, fieldName, value, opts, hushTag, depth+1)
	default:
		return ht.processSimpleField(fieldName, field, value, hushTag, opts)
	}
}

// processStruct handles the processing of struct fields.
func (ht *hushType) processStruct(ctx context.Context, rv reflect.Value, prefix string, opts *hushOptions, depth int) ([][]string, error) {
	data := make([][]string, 0, rv.NumField())
	errChan := make(chan error, rv.NumField())
	var wg sync.WaitGroup
	var mu sync.Mutex

	sem := make(chan struct{}, maxConcurrency)

	t := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := t.Field(i)
		value := rv.Field(i)

		if !opts.includePrivate && !field.IsExported() {
			continue
		}

		wg.Add(1)
		sem <- struct{}{} // acquire semaphore slot
		go func(field reflect.StructField, value reflect.Value) {
			defer wg.Done()
			defer func() { <-sem }() // release semaphore slot

			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
			}

			fieldName := buildFieldName(prefix, field.Name, opts.separator)

			result, err := ht.processValue(ctx, fieldName, field, value, opts, "", depth)
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
func (ht *hushType) processSliceOrArray(ctx context.Context, fieldName string, value reflect.Value, opts *hushOptions, hushTag string, depth int) ([][]string, error) {
	// Handle basic types more elegantly
	if isBasicType := isBasicTypeKind(value.Type().Elem().Kind()); isBasicType {
		return ht.processBasicTypeSlice(fieldName, value, opts, hushTag)
	}

	return ht.processComplexTypeSlice(ctx, fieldName, value, opts, hushTag, depth)
}

// Helper function to determine if a kind is a basic type
func isBasicTypeKind(kind reflect.Kind) bool {
	basicTypes := map[reflect.Kind]bool{
		reflect.String:  true,
		reflect.Bool:    true,
		reflect.Int:     true,
		reflect.Int8:    true,
		reflect.Int16:   true,
		reflect.Int32:   true,
		reflect.Int64:   true,
		reflect.Uint:    true,
		reflect.Uint8:   true,
		reflect.Uint16:  true,
		reflect.Uint32:  true,
		reflect.Uint64:  true,
		reflect.Float32: true,
		reflect.Float64: true,
	}
	return basicTypes[kind]
}

// Process slices of basic types
func (ht *hushType) processBasicTypeSlice(fieldName string, value reflect.Value, opts *hushOptions, hushTag string) ([][]string, error) {
	result := make([][]string, 0, value.Len())
	for i := 0; i < value.Len(); i++ {
		elemFieldName := fmt.Sprintf("%s[%d]", fieldName, i)
		elemValue := value.Index(i)
		convertedString := convertNonCompositeToString(elemValue)
		elemResult := processString(elemFieldName, convertedString, hushTag, opts)
		result = append(result, elemResult...)
	}
	return result, nil
}

// Process slices of complex types
func (ht *hushType) processComplexTypeSlice(ctx context.Context, fieldName string, value reflect.Value, opts *hushOptions, hushTag string, depth int) ([][]string, error) {
	result := make([][]string, 0, value.Len())
	for i := 0; i < value.Len(); i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		elemFieldName := fmt.Sprintf("%s[%d]", fieldName, i)
		elemValue := value.Index(i)

		elemResult, err := ht.processValue(ctx, elemFieldName, reflect.StructField{}, elemValue, opts, hushTag, depth)
		if err != nil {
			return nil, err
		}
		result = append(result, elemResult...)
	}
	return result, nil
}

// processMap handles the processing of map fields.
func (ht *hushType) processMap(ctx context.Context, fieldName string, value reflect.Value, opts *hushOptions, hushTag string, depth int) ([][]string, error) {
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

		processedValue, err := ht.processValue(ctx, mapFieldName, reflect.StructField{}, fieldValue, opts, hushTag, depth)
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
