package common

import "time"

func (t *Tools) String(v string) *string {
	return &v
}

func (t *Tools) StringValue(v *string) string {
	if v != nil {
		return *v
	}

	return ""
}

func (t *Tools) StringSlice(src []string) []*string {
	dst := make([]*string, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = &(src[i])
	}

	return dst
}

func (t *Tools) StringValueSlice(src []*string) []string {
	dst := make([]string, len(src))
	for i := 0; i < len(src); i++ {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}

	return dst
}

func (t *Tools) StringMap(src map[string]string) map[string]*string {
	dst := make(map[string]*string)
	for k, val := range src {
		v := val
		dst[k] = &v
	}

	return dst
}

func (t *Tools) StringValueMap(src map[string]*string) map[string]string {
	dst := make(map[string]string)
	for k, val := range src {
		if val != nil {
			dst[k] = *val
		}
	}

	return dst
}

func (t *Tools) Bool(v bool) *bool {
	return &v
}

func (t *Tools) BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}

	return false
}

func (t *Tools) BoolSlice(src []bool) []*bool {
	dst := make([]*bool, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = &(src[i])
	}

	return dst
}

func (t *Tools) BoolValueSlice(src []*bool) []bool {
	dst := make([]bool, len(src))
	for i := 0; i < len(src); i++ {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}
	return dst
}

func (t *Tools) BoolMap(src map[string]bool) map[string]*bool {
	dst := make(map[string]*bool)
	for k, val := range src {
		v := val
		dst[k] = &v
	}

	return dst
}

func (t *Tools) BoolValueMap(src map[string]*bool) map[string]bool {
	dst := make(map[string]bool)
	for k, val := range src {
		if val != nil {
			dst[k] = *val
		}
	}

	return dst
}

func (t *Tools) Uint(v uint) *uint {
	return &v
}

func (t *Tools) UintValue(v *uint) uint {
	if v != nil {
		return *v
	}

	return 0
}

func (t *Tools) Int(v int) *int {
	return &v
}

func (t *Tools) IntValue(v *int) int {
	if v != nil {
		return *v
	}

	return 0
}

func (t *Tools) IntSlice(src []int) []*int {
	dst := make([]*int, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = &(src[i])
	}

	return dst
}

func (t *Tools) IntValueSlice(src []*int) []int {
	dst := make([]int, len(src))
	for i := 0; i < len(src); i++ {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}

	return dst
}

func (t *Tools) IntMap(src map[string]int) map[string]*int {
	dst := make(map[string]*int)
	for k, val := range src {
		v := val
		dst[k] = &v
	}

	return dst
}

func (t *Tools) IntValueMap(src map[string]*int) map[string]int {
	dst := make(map[string]int)
	for k, val := range src {
		if val != nil {
			dst[k] = *val
		}
	}

	return dst
}

func (t *Tools) Int64(v int64) *int64 {
	return &v
}

func (t *Tools) Int64Value(v *int64) int64 {
	if v != nil {
		return *v
	}

	return 0
}

func (t *Tools) Int64Slice(src []int64) []*int64 {
	dst := make([]*int64, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = &(src[i])
	}

	return dst
}

func (t *Tools) Int64ValueSlice(src []*int64) []int64 {
	dst := make([]int64, len(src))
	for i := 0; i < len(src); i++ {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}

	return dst
}

func (t *Tools) Int64Map(src map[string]int64) map[string]*int64 {
	dst := make(map[string]*int64)
	for k, val := range src {
		v := val
		dst[k] = &v
	}

	return dst
}

func (t *Tools) Int64ValueMap(src map[string]*int64) map[string]int64 {
	dst := make(map[string]int64)
	for k, val := range src {
		if val != nil {
			dst[k] = *val
		}
	}

	return dst
}

func (t *Tools) Float64(v float64) *float64 {
	return &v
}

func (t *Tools) Float64Value(v *float64) float64 {
	if v != nil {
		return *v
	}

	return 0
}

func (t *Tools) Float64Slice(src []float64) []*float64 {
	dst := make([]*float64, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = &(src[i])
	}

	return dst
}

func (t *Tools) Float64ValueSlice(src []*float64) []float64 {
	dst := make([]float64, len(src))
	for i := 0; i < len(src); i++ {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}

	return dst
}

func (t *Tools) Float64Map(src map[string]float64) map[string]*float64 {
	dst := make(map[string]*float64)
	for k, val := range src {
		v := val
		dst[k] = &v
	}
	return dst
}

func (t *Tools) Float64ValueMap(src map[string]*float64) map[string]float64 {
	dst := make(map[string]float64)
	for k, val := range src {
		if val != nil {
			dst[k] = *val
		}
	}

	return dst
}

func (t *Tools) Time(v time.Time) *time.Time {
	return &v
}

func TimeValue(v *time.Time) time.Time {
	if v != nil {
		return *v
	}

	return time.Time{}
}

func (t *Tools) SecondsTimeValue(v *int64) time.Time {
	if v != nil {
		return time.Unix((*v / 1000), 0)
	}

	return time.Time{}
}

func (t *Tools) MillisecondsTimeValue(v *int64) time.Time {
	if v != nil {
		return time.Unix(0, (*v * 1000000))
	}

	return time.Time{}
}

func (t *Tools) TimeUnixMilli(tm time.Time) int64 {
	return tm.UnixNano() / int64(time.Millisecond/time.Nanosecond)
}

func (t *Tools) TimeSlice(src []time.Time) []*time.Time {
	dst := make([]*time.Time, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = &(src[i])
	}

	return dst
}

func (t *Tools) TimeValueSlice(src []*time.Time) []time.Time {
	dst := make([]time.Time, len(src))
	for i := 0; i < len(src); i++ {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}

	return dst
}

func (t *Tools) TimeMap(src map[string]time.Time) map[string]*time.Time {
	dst := make(map[string]*time.Time)
	for k, val := range src {
		v := val
		dst[k] = &v
	}

	return dst
}

func (t *Tools) TimeValueMap(src map[string]*time.Time) map[string]time.Time {
	dst := make(map[string]time.Time)
	for k, val := range src {
		if val != nil {
			dst[k] = *val
		}
	}

	return dst
}
