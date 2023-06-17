package scenario

import (
	"encoding/json"
	"github.com/pkg/errors"
	"reflect"
)

type Artifacts struct {
	artifacts map[string]any
}

func NewArtifacts(items ...map[string]any) (*Artifacts, error) {
	ret := &Artifacts{
		artifacts: map[string]any{},
	}
	for _, item := range items {
		for key, value := range item {
			err := ret.Add(key, value)
			if err != nil {
				return nil, err
			}
		}
	}

	return ret, nil
}

func isZeroValue(value reflect.Value) bool {
	return value.Kind() != reflect.Bool && value.IsZero()
}

func (a *Artifacts) Add(key string, value any) error {
	if key == "" {
		return errors.Errorf("key(%v) is empty", key)
	}

	if isZeroValue(reflect.ValueOf(value)) {
		return errors.Errorf("key(%v)'s value is empty", key)
	}

	if a.isExistArtifact(key) {
		return errors.Errorf("key(%v) is already exists", key)
	}

	a.artifacts[key] = value
	return nil
}

func (a *Artifacts) Fill(item any) error {
	marshal, _ := json.Marshal(a.artifacts)
	err := json.Unmarshal(marshal, item)
	if err != nil {
		return errors.WithStack(err)
	}

	fields := reflect.ValueOf(item).Elem()

	for i := 0; i < fields.NumField(); i++ {
		if isZeroValue(fields.Field(i)) {
			return errors.Errorf("field(%v.%v) is zero", reflect.TypeOf(item).Elem().Name(), fields.Type().Field(i).Name)
		}
	}

	return nil
}

func (a *Artifacts) IsExists(keys ...string) bool {
	for _, key := range keys {
		if !a.isExistArtifact(key) {
			return false
		}
	}
	return true
}

func (a *Artifacts) isExistArtifact(key string) bool {
	_, ok := a.artifacts[key]
	return ok
}
