/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package storageclass

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/storage"
	"k8s.io/kubernetes/pkg/apis/storage/validation"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/registry/generic"
	"k8s.io/kubernetes/pkg/runtime"
	apistorage "k8s.io/kubernetes/pkg/storage"
	"k8s.io/kubernetes/pkg/util/validation/field"
)

// storageClassStrategy implements behavior for StorageClass objects
type storageClassStrategy struct {
	runtime.ObjectTyper
	api.NameGenerator
}

// Strategy is the default logic that applies when creating and updating
// StorageClass objects via the REST API.
var Strategy = storageClassStrategy{api.Scheme, api.SimpleNameGenerator}

func (storageClassStrategy) NamespaceScoped() bool {
	return false
}

// ResetBeforeCreate clears the Status field which is not allowed to be set by end users on creation.
func (storageClassStrategy) PrepareForCreate(ctx api.Context, obj runtime.Object) {
	_ = obj.(*storage.StorageClass)
}

func (storageClassStrategy) Validate(ctx api.Context, obj runtime.Object) field.ErrorList {
	storageClass := obj.(*storage.StorageClass)
	return validation.ValidateStorageClass(storageClass)
}

// Canonicalize normalizes the object after validation.
func (storageClassStrategy) Canonicalize(obj runtime.Object) {
}

func (storageClassStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForUpdate sets the Status fields which is not allowed to be set by an end user updating a PV
func (storageClassStrategy) PrepareForUpdate(ctx api.Context, obj, old runtime.Object) {
	_ = obj.(*storage.StorageClass)
	_ = old.(*storage.StorageClass)
}

func (storageClassStrategy) ValidateUpdate(ctx api.Context, obj, old runtime.Object) field.ErrorList {
	errorList := validation.ValidateStorageClass(obj.(*storage.StorageClass))
	return append(errorList, validation.ValidateStorageClassUpdate(obj.(*storage.StorageClass), old.(*storage.StorageClass))...)
}

func (storageClassStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	cls, ok := obj.(*storage.StorageClass)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not of type StorageClass")
	}
	return labels.Set(cls.ObjectMeta.Labels), StorageClassToSelectableFields(cls), nil
}

// MatchStorageClass returns a generic matcher for a given label and field selector.
func MatchStorageClasses(label labels.Selector, field fields.Selector) apistorage.SelectionPredicate {
	return apistorage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// StorageClassToSelectableFields returns a label set that represents the object
func StorageClassToSelectableFields(storageClass *storage.StorageClass) fields.Set {
	return generic.ObjectMetaFieldsSet(&storageClass.ObjectMeta, false)
}
