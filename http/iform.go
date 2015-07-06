package http

import (
	"net/http"
	"strconv"
)

type IFormValueGetter interface {
	HasFormValue(key string) bool
	FormValue(key string) string
	// Get value from request.Form.
	// The second returned value indicates whether the key exists.
	FormValueOk(key string) (string, bool)
	// The second returned value is the strconv parse error.
	FormValueInt(key string) (int, error)
	FormValueInt8(key string) (int8, error)
	FormValueInt16(key string) (int16, error)
	FormValueInt32(key string) (int32, error)
	FormValueInt64(key string) (int64, error)
	FormValueUint(key string) (uint, error)
	FormValueUint8(key string) (uint8, error)
	FormValueUint16(key string) (uint16, error)
	FormValueUint32(key string) (uint32, error)
	FormValueUint64(key string) (uint64, error)
	FormValueFloat32(key string) (float32, error)
	FormValueFloat64(key string) (float64, error)
	FormValueBool(key string) (bool, error)

	HasPostFormValue(key string) bool
	PostFormValue(key string) string
	PostFormValueOk(key string) (string, bool)
	PostFormValueInt(key string) (int, error)
	PostFormValueInt8(key string) (int8, error)
	PostFormValueInt16(key string) (int16, error)
	PostFormValueInt32(key string) (int32, error)
	PostFormValueInt64(key string) (int64, error)
	PostFormValueUint(key string) (uint, error)
	PostFormValueUint8(key string) (uint8, error)
	PostFormValueUint16(key string) (uint16, error)
	PostFormValueUint32(key string) (uint32, error)
	PostFormValueUint64(key string) (uint64, error)
	PostFormValueFloat32(key string) (float32, error)
	PostFormValueFloat64(key string) (float64, error)
	PostFormValueBool(key string) (bool, error)

	// Array values.
	FormValues(key string) []string
	FormValueInts(key string) ([]int, error)
	FormValueInt8s(key string) ([]int8, error)
	FormValueInt16s(key string) ([]int16, error)
	FormValueInt32s(key string) ([]int32, error)
	FormValueInt64s(key string) ([]int64, error)
	FormValueUints(key string) ([]uint, error)
	FormValueUint8s(key string) ([]uint8, error)
	FormValueUint16s(key string) ([]uint16, error)
	FormValueUint32s(key string) ([]uint32, error)
	FormValueUint64s(key string) ([]uint64, error)
	FormValueFloat32s(key string) ([]float32, error)
	FormValueFloat64s(key string) ([]float64, error)
	FormValueBools(key string) ([]bool, error)

	PostFormValues(key string) []string
	PostFormValueInts(key string) ([]int, error)
	PostFormValueInt8s(key string) ([]int8, error)
	PostFormValueInt16s(key string) ([]int16, error)
	PostFormValueInt32s(key string) ([]int32, error)
	PostFormValueInt64s(key string) ([]int64, error)
	PostFormValueUints(key string) ([]uint, error)
	PostFormValueUint8s(key string) ([]uint8, error)
	PostFormValueUint16s(key string) ([]uint16, error)
	PostFormValueUint32s(key string) ([]uint32, error)
	PostFormValueUint64s(key string) ([]uint64, error)
	PostFormValueFloat32s(key string) ([]float32, error)
	PostFormValueFloat64s(key string) ([]float64, error)
	PostFormValueBools(key string) ([]bool, error)
}

type requestFormValueGetter struct{ *http.Request }

func RequestFormValueGetter(r *http.Request) IFormValueGetter {
	// r.PostFormValue(key) will call r.ParsePostForm() if necessary.
	// And here, the error r.ParsePostForm() returns will be ignored.
	_ = r.PostFormValue("")
	return &requestFormValueGetter{r}
}

// -- FormValue* --

func (r *requestFormValueGetter) HasFormValue(key string) bool {
	_, ok := r.Form[key]
	return ok
}

func (r *requestFormValueGetter) FormValue(key string) string {
	return r.Request.FormValue(key)
}

func (r *requestFormValueGetter) FormValueOk(key string) (string, bool) {
	_, ok := r.Form[key]
	if !ok {
		return "", false
	}
	return r.FormValue(key), true
}
func (r *requestFormValueGetter) FormValueInt(key string) (int, error) {
	tmp, err := strconv.ParseInt(r.FormValue(key), 10, strconv.IntSize)
	return int(tmp), err
}
func (r *requestFormValueGetter) FormValueInt8(key string) (int8, error) {
	tmp, err := strconv.ParseInt(r.FormValue(key), 10, 8)
	return int8(tmp), err
}
func (r *requestFormValueGetter) FormValueInt16(key string) (int16, error) {
	tmp, err := strconv.ParseInt(r.FormValue(key), 10, 16)
	return int16(tmp), err
}
func (r *requestFormValueGetter) FormValueInt32(key string) (int32, error) {
	tmp, err := strconv.ParseInt(r.FormValue(key), 10, 32)
	return int32(tmp), err
}
func (r *requestFormValueGetter) FormValueInt64(key string) (int64, error) {
	tmp, err := strconv.ParseInt(r.FormValue(key), 10, 64)
	return int64(tmp), err
}
func (r *requestFormValueGetter) FormValueUint(key string) (uint, error) {
	tmp, err := strconv.ParseUint(r.FormValue(key), 10, strconv.IntSize)
	return uint(tmp), err
}
func (r *requestFormValueGetter) FormValueUint8(key string) (uint8, error) {
	tmp, err := strconv.ParseUint(r.FormValue(key), 10, 8)
	return uint8(tmp), err
}
func (r *requestFormValueGetter) FormValueUint16(key string) (uint16, error) {
	tmp, err := strconv.ParseUint(r.FormValue(key), 10, 16)
	return uint16(tmp), err
}
func (r *requestFormValueGetter) FormValueUint32(key string) (uint32, error) {
	tmp, err := strconv.ParseUint(r.FormValue(key), 10, 32)
	return uint32(tmp), err
}
func (r *requestFormValueGetter) FormValueUint64(key string) (uint64, error) {
	tmp, err := strconv.ParseUint(r.FormValue(key), 10, 64)
	return uint64(tmp), err
}
func (r *requestFormValueGetter) FormValueFloat32(key string) (float32, error) {
	tmp, err := strconv.ParseFloat(r.FormValue(key), 32)
	return float32(tmp), err
}
func (r *requestFormValueGetter) FormValueFloat64(key string) (float64, error) {
	tmp, err := strconv.ParseFloat(r.FormValue(key), 64)
	return float64(tmp), err
}
func (r *requestFormValueGetter) FormValueBool(key string) (bool, error) {
	return strconv.ParseBool(r.FormValue(key))
}

// -- PostFormValue* --

func (r *requestFormValueGetter) HasPostFormValue(key string) bool {
	_, ok := r.PostForm[key]
	return ok
}

func (r *requestFormValueGetter) PostFormValue(key string) string {
	return r.Request.PostFormValue(key)
}

func (r *requestFormValueGetter) PostFormValueOk(key string) (string, bool) {
	_, ok := r.PostForm[key]
	if !ok {
		return "", false
	}
	return r.PostFormValue(key), true
}
func (r *requestFormValueGetter) PostFormValueInt(key string) (int, error) {
	tmp, err := strconv.ParseInt(r.PostFormValue(key), 10, strconv.IntSize)
	return int(tmp), err
}
func (r *requestFormValueGetter) PostFormValueInt8(key string) (int8, error) {
	tmp, err := strconv.ParseInt(r.PostFormValue(key), 10, 8)
	return int8(tmp), err
}
func (r *requestFormValueGetter) PostFormValueInt16(key string) (int16, error) {
	tmp, err := strconv.ParseInt(r.PostFormValue(key), 10, 16)
	return int16(tmp), err
}
func (r *requestFormValueGetter) PostFormValueInt32(key string) (int32, error) {
	tmp, err := strconv.ParseInt(r.PostFormValue(key), 10, 32)
	return int32(tmp), err
}
func (r *requestFormValueGetter) PostFormValueInt64(key string) (int64, error) {
	tmp, err := strconv.ParseInt(r.PostFormValue(key), 10, 64)
	return int64(tmp), err
}
func (r *requestFormValueGetter) PostFormValueUint(key string) (uint, error) {
	tmp, err := strconv.ParseUint(r.PostFormValue(key), 10, strconv.IntSize)
	return uint(tmp), err
}
func (r *requestFormValueGetter) PostFormValueUint8(key string) (uint8, error) {
	tmp, err := strconv.ParseUint(r.PostFormValue(key), 10, 8)
	return uint8(tmp), err
}
func (r *requestFormValueGetter) PostFormValueUint16(key string) (uint16, error) {
	tmp, err := strconv.ParseUint(r.PostFormValue(key), 10, 16)
	return uint16(tmp), err
}
func (r *requestFormValueGetter) PostFormValueUint32(key string) (uint32, error) {
	tmp, err := strconv.ParseUint(r.PostFormValue(key), 10, 32)
	return uint32(tmp), err
}
func (r *requestFormValueGetter) PostFormValueUint64(key string) (uint64, error) {
	tmp, err := strconv.ParseUint(r.PostFormValue(key), 10, 64)
	return uint64(tmp), err
}
func (r *requestFormValueGetter) PostFormValueFloat32(key string) (float32, error) {
	tmp, err := strconv.ParseFloat(r.PostFormValue(key), 32)
	return float32(tmp), err
}
func (r *requestFormValueGetter) PostFormValueFloat64(key string) (float64, error) {
	tmp, err := strconv.ParseFloat(r.PostFormValue(key), 64)
	return float64(tmp), err
}
func (r *requestFormValueGetter) PostFormValueBool(key string) (bool, error) {
	return strconv.ParseBool(r.PostFormValue(key))
}

// -- FormValue*s --

func (r *requestFormValueGetter) FormValues(key string) []string {
	return r.Form[key]
}
func (r *requestFormValueGetter) FormValueInts(key string) ([]int, error) {
	vs := r.Form[key]
	results := make([]int, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseInt(v, 10, strconv.IntSize); err != nil {
			return nil, err
		} else {
			results[i] = int(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) FormValueInt8s(key string) ([]int8, error) {
	vs := r.Form[key]
	results := make([]int8, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseInt(v, 10, 8); err != nil {
			return nil, err
		} else {
			results[i] = int8(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) FormValueInt16s(key string) ([]int16, error) {
	vs := r.Form[key]
	results := make([]int16, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseInt(v, 10, 16); err != nil {
			return nil, err
		} else {
			results[i] = int16(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) FormValueInt32s(key string) ([]int32, error) {
	vs := r.Form[key]
	results := make([]int32, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseInt(v, 10, 32); err != nil {
			return nil, err
		} else {
			results[i] = int32(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) FormValueInt64s(key string) ([]int64, error) {
	vs := r.Form[key]
	results := make([]int64, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseInt(v, 10, 64); err != nil {
			return nil, err
		} else {
			results[i] = int64(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) FormValueUints(key string) ([]uint, error) {
	vs := r.Form[key]
	results := make([]uint, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseUint(v, 10, strconv.IntSize); err != nil {
			return nil, err
		} else {
			results[i] = uint(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) FormValueUint8s(key string) ([]uint8, error) {
	vs := r.Form[key]
	results := make([]uint8, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseUint(v, 10, 8); err != nil {
			return nil, err
		} else {
			results[i] = uint8(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) FormValueUint16s(key string) ([]uint16, error) {
	vs := r.Form[key]
	results := make([]uint16, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseUint(v, 10, 16); err != nil {
			return nil, err
		} else {
			results[i] = uint16(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) FormValueUint32s(key string) ([]uint32, error) {
	vs := r.Form[key]
	results := make([]uint32, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseUint(v, 10, 32); err != nil {
			return nil, err
		} else {
			results[i] = uint32(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) FormValueUint64s(key string) ([]uint64, error) {
	vs := r.Form[key]
	results := make([]uint64, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseUint(v, 10, 64); err != nil {
			return nil, err
		} else {
			results[i] = uint64(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) FormValueFloat32s(key string) ([]float32, error) {
	vs := r.Form[key]
	results := make([]float32, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseFloat(v, 32); err != nil {
			return nil, err
		} else {
			results[i] = float32(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) FormValueFloat64s(key string) ([]float64, error) {
	vs := r.Form[key]
	results := make([]float64, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseFloat(v, 64); err != nil {
			return nil, err
		} else {
			results[i] = float64(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) FormValueBools(key string) ([]bool, error) {
	vs := r.Form[key]
	results := make([]bool, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseBool(v); err != nil {
			return nil, err
		} else {
			results[i] = bool(tmp)
		}
	}
	return results, nil
}

// -- PostFormValue*s --

func (r *requestFormValueGetter) PostFormValues(key string) []string {
	return r.PostForm[key]
}
func (r *requestFormValueGetter) PostFormValueInts(key string) ([]int, error) {
	vs := r.PostForm[key]
	results := make([]int, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseInt(v, 10, strconv.IntSize); err != nil {
			return nil, err
		} else {
			results[i] = int(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) PostFormValueInt8s(key string) ([]int8, error) {
	vs := r.PostForm[key]
	results := make([]int8, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseInt(v, 10, 8); err != nil {
			return nil, err
		} else {
			results[i] = int8(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) PostFormValueInt16s(key string) ([]int16, error) {
	vs := r.PostForm[key]
	results := make([]int16, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseInt(v, 10, 16); err != nil {
			return nil, err
		} else {
			results[i] = int16(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) PostFormValueInt32s(key string) ([]int32, error) {
	vs := r.PostForm[key]
	results := make([]int32, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseInt(v, 10, 32); err != nil {
			return nil, err
		} else {
			results[i] = int32(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) PostFormValueInt64s(key string) ([]int64, error) {
	vs := r.PostForm[key]
	results := make([]int64, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseInt(v, 10, 64); err != nil {
			return nil, err
		} else {
			results[i] = int64(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) PostFormValueUints(key string) ([]uint, error) {
	vs := r.PostForm[key]
	results := make([]uint, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseUint(v, 10, strconv.IntSize); err != nil {
			return nil, err
		} else {
			results[i] = uint(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) PostFormValueUint8s(key string) ([]uint8, error) {
	vs := r.PostForm[key]
	results := make([]uint8, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseUint(v, 10, 8); err != nil {
			return nil, err
		} else {
			results[i] = uint8(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) PostFormValueUint16s(key string) ([]uint16, error) {
	vs := r.PostForm[key]
	results := make([]uint16, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseUint(v, 10, 16); err != nil {
			return nil, err
		} else {
			results[i] = uint16(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) PostFormValueUint32s(key string) ([]uint32, error) {
	vs := r.PostForm[key]
	results := make([]uint32, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseUint(v, 10, 32); err != nil {
			return nil, err
		} else {
			results[i] = uint32(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) PostFormValueUint64s(key string) ([]uint64, error) {
	vs := r.PostForm[key]
	results := make([]uint64, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseUint(v, 10, 64); err != nil {
			return nil, err
		} else {
			results[i] = uint64(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) PostFormValueFloat32s(key string) ([]float32, error) {
	vs := r.PostForm[key]
	results := make([]float32, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseFloat(v, 32); err != nil {
			return nil, err
		} else {
			results[i] = float32(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) PostFormValueFloat64s(key string) ([]float64, error) {
	vs := r.PostForm[key]
	results := make([]float64, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseFloat(v, 64); err != nil {
			return nil, err
		} else {
			results[i] = float64(tmp)
		}
	}
	return results, nil
}
func (r *requestFormValueGetter) PostFormValueBools(key string) ([]bool, error) {
	vs := r.PostForm[key]
	results := make([]bool, len(vs))
	for i, v := range vs {
		if tmp, err := strconv.ParseBool(v); err != nil {
			return nil, err
		} else {
			results[i] = bool(tmp)
		}
	}
	return results, nil
}
