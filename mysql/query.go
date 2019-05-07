package mysql

import (
	"fmt"
	"github.com/hhy5861/go-common/common"
	"reflect"

	"github.com/jinzhu/gorm"
)

type (
	Service struct {
	}

	QueryResult struct {
		Total   int64       `json:"total"`
		HasNext bool        `json:"hasNext"`
		Size    int64       `json:"size"`
		Page    int64       `json:"page"`
		Data    interface{} `json:"data"`
	}

	QueryPage struct {
		Page *int64 `json:"page" form:"page" binding:"omitempty,numeric"`
		Size *int64 `json:"size" form:"size" binding:"omitempty,numeric"`
	}
)

var (
	tools *common.Tools
)

func init() {
	tools = common.NewTools()
}

func (s *Service) QueryResList(
	query *gorm.DB,
	tableName string, list interface{},
	page, size int64) (*QueryResult, error) {

	var total int64

	queryCount := query
	queryCount.Table(tableName).Count(&total)
	if total > 0 {
		offset := tools.SetOffset(s.GetPage(page), s.GetSize(size))
		err = query.Table(tableName).Limit(s.GetSize(size)).Offset(offset).Find(list).Error
		if err != nil {
			return nil, err
		}
	}

	result := &QueryResult{
		Data:    list,
		HasNext: total > (s.GetPage(page) * s.GetSize(size)),
		Total:   total,
		Page:    s.GetPage(page),
		Size:    s.GetSize(size),
	}

	return result, err
}

const (
	defaultSize = 100
)

func (s *Service) PageQueryReflect(
	query *gorm.DB,
	tableName string, list interface{},
	page, size int64,
	class interface{},
	name string) (*QueryResult, error) {

	var total int64
	queryCount := query
	queryCount.Table(tableName).Count(&total)
	if total > 0 {
		offset := tools.SetOffset(s.GetPage(page), s.GetSize(size))
		err := query.Table(tableName).Limit(s.GetSize(size)).Offset(offset).Find(list).Error
		if err != nil {
			return nil, err
		}
	}

	result := &QueryResult{
		Data:    list,
		HasNext: total > (s.GetPage(page) * s.GetSize(size)),
		Total:   total,
		Page:    s.GetPage(page),
		Size:    s.GetSize(size),
	}

	if class != nil && name != "" {
		r, err := s.Calls(class, name, list)
		if err != nil {
			return nil, err
		}

		result.Data = r[0].Interface()
	}

	return result, err
}

func (s *Service) PageResult(
	list interface{},
	total, page, size int64) (*QueryResult, error) {

	return &QueryResult{
		Data:    list,
		HasNext: total > (s.GetPage(page) * s.GetSize(size)),
		Total:   total,
		Page:    s.GetPage(page),
		Size:    s.GetSize(size),
	}, nil
}

func (s *Service) Calls(
	myClass interface{},
	name string,
	params ...interface{}) ([]reflect.Value, error) {

	myClassValue := reflect.ValueOf(myClass)
	m := myClassValue.MethodByName(name)

	if !m.IsValid() {
		err := fmt.Errorf("method not found param name: %s", name)
		return nil, err
	}

	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}

	return m.Call(in), nil
}

func (s *Service) GetPage(page int64) int64 {
	if page <= 0 {
		page = 1
	}

	return page
}

func (s *Service) GetSize(size int64) int64 {
	if size <= 0 {
		size = defaultSize
	}

	return size
}
