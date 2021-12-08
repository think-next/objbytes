package objbytes

import (
    "reflect"
    "unsafe"
)

func Marshal(v interface{}) ([]byte, error) {

    rv := reflect.ValueOf(v)
    switch rv.Type().Kind() {
    case reflect.Struct:

    case reflect.Ptr:
    case reflect.String:
    case reflect.Slice:
    default:

    }

    return nil, nil
}

func StructMarsh(data reflect.Value) ([]byte, error) {
    var buffer []byte

    //  deep copy
    rv := *(*Value)(unsafe.Pointer(&data))
    bh := (*SliceHeader)(unsafe.Pointer(&buffer))
    bh.Data = uintptr(rv.ptr)
    bh.Len = int(data.Type().Size())
    bh.Cap = bh.Len

    for i := 0; i < data.NumField(); i++ {
        switch data.Field(i).Type().Kind() {
        case reflect.Struct:
        case reflect.Slice:
        case reflect.String:
        case reflect.Ptr:
            Marshal(data.Field(i))
        }
    }

    return buffer, nil
}
