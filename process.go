package hush

import (
	"context"
	"fmt"
	"reflect"
	"sort"
)

// processField is the main function for processing individual fields.
// It handles different types of fields (struct, pointer, slice, map, etc.) and applies the appropriate masking.
func processField(ctx context.Context, fieldName string, field reflect.StructField, value reflect.Value, opts *hushOptions) ([][]string, error) {
	hushTag := field.Tag.Get("hush")

	// Check if the field is unexported and handle accordingly
	if field.PkgPath != "" && !opts.includePrivate {
		return nil, nil // Skip unexported fields when not including private fields
	}

	switch value.Kind() {
	case reflect.Struct:
		nestedHush := hushType{value: value, isStr: false}
		nestedOpts := &hushOptions{
			separator:      opts.separator,
			maskFunc:       opts.maskFunc,
			includePrivate: opts.includePrivate,
		}
		return nestedHush.Hush(ctx, fieldName, WithOptions(nestedOpts))
	case reflect.Ptr:
		if value.IsNil() {
			return [][]string{{fieldName, "nil"}}, nil
		}
		return processField(ctx, fieldName, reflect.StructField{}, value.Elem(), opts)
	case reflect.Slice, reflect.Array:
		return processSliceOrArray(ctx, fieldName, value, opts)
	case reflect.Map:
		return processMap(ctx, fieldName, value, opts)
	default:
		return processSimpleField(fieldName, field, value, hushTag, opts)
	}
}

// processSimpleField handles the processing of simple (non-composite) fields.
func processSimpleField(fieldName string, field reflect.StructField, value reflect.Value, hushTag string, opts *hushOptions) ([][]string, error) {
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
	if hushTag == "" {
		return [][]string{{fieldName, value}}
	}
	if hushTag == TagHide {
		return [][]string{{fieldName, HiddenValue}}
	}
	if (hushTag == TagMask && opts.maskFunc != nil) || (opts.includePrivate && opts.maskFunc != nil) {
		return [][]string{{fieldName, opts.maskFunc(value)}}
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
func processSliceOrArray(ctx context.Context, fieldName string, value reflect.Value, opts *hushOptions) ([][]string, error) {
	var result [][]string
	for i := 0; i < value.Len(); i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		elemFieldName := fmt.Sprintf("%s[%d]", fieldName, i)
		elemValue := value.Index(i)

		elemResult, err := processField(ctx, elemFieldName, reflect.StructField{}, elemValue, opts)
		if err != nil {
			return nil, err
		}
		result = append(result, elemResult...)
	}
	return result, nil
}

// processMap handles the processing of map fields.
func processMap(ctx context.Context, fieldName string, value reflect.Value, opts *hushOptions) ([][]string, error) {
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

		processedValue, err := processField(ctx, mapFieldName, reflect.StructField{}, fieldValue, opts)
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
