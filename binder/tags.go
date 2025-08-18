package binder

import (
	"reflect"
	
	"BinGo/enum"
)

func hasTags(dst any) map[enum.Tag]struct{} {
	tags := make(map[enum.Tag]struct{}, len(enum.Tags.Values()))
	hasTagsRecursive(reflect.TypeOf(dst).Elem(), tags)
	return tags
}

func hasTagsRecursive(t reflect.Type, tags map[enum.Tag]struct{}) {
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct && field.Type.String() != "time.Time" {
			hasTagsRecursive(field.Type, tags)
			return
		}
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			hasTagsRecursive(field.Type.Elem(), tags)
			return
		}
		for _, tag := range enum.Tags.Values() {
			if _, ok := tags[tag]; !ok {
				if _, exists := field.Tag.Lookup(tag.String()); exists {
					tags[tag] = struct{}{}
				}
			}
		}
	}
}

func (b *dataBind) getTags() {
	getTagsRecursive(reflect.TypeOf(b.DataDist).Elem(), b.tag, b.data)
}

func getTagsRecursive(t reflect.Type, tag enum.Tag, data map[string]any) {
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct && field.Type.String() != "time.Time" {
			getTagsRecursive(field.Type, tag, data)
			return
		}
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			getTagsRecursive(field.Type.Elem(), tag, data)
			return
		}
		if val, ok := field.Tag.Lookup(tag.String()); ok {
			data[val] = struct{}{}
		}
	}
}
