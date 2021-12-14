package objbytes

import (
    "fmt"
    "reflect"
    "unsafe"
)

func Marshal(v interface{}) ([]byte, error) {

    rv := reflect.ValueOf(v)
    switch rv.Type().Kind() {
    case reflect.Struct:
        return marshal(v)
    case reflect.Ptr:
    case reflect.String:
    case reflect.Slice:
    default:

    }

    return nil, nil
}

func marshal(data interface{}) ([]byte, error) {

    objBytes := newObjBytes()
    t := reflect.TypeOf(data)
    switch t.Kind() {
    case reflect.Struct:
        _ = objBytes.MarshalStruct(data)
    }

    return objBytes.data, nil
}

type ObjBytes struct {
    data        []byte
    pairsOffset []PairOffset
}

func newObjBytes() *ObjBytes {
    return new(ObjBytes)
}

func (o *ObjBytes) MarshalStruct(obj interface{}) error {
    t := reflect.TypeOf(obj)
    if t.Kind() != reflect.Struct {
        panic("the kind of obj must be struct")
    }

    var buffer []byte
    v := reflect.ValueOf(obj)
    rv := *(*Value)(unsafe.Pointer(&v))
    bh := (*SliceHeader)(unsafe.Pointer(&buffer))
    bh.Data = uintptr(rv.ptr)
    bh.Len = int(t.Size())
    bh.Cap = bh.Len

    o.data = buffer
    o.iteratorStruct(v)
    o.joinOverHead()
    return nil
}

func (o *ObjBytes) iteratorStruct(obj reflect.Value) error {

    for i := 0; i < obj.NumField(); i++ {
        t := obj.Field(i).Type()
        switch t.Kind() {
        case reflect.Struct:

        case reflect.Slice:
            o.appendSlice(uint32(obj.Type().Field(i).Offset), obj.Field(i))
        case reflect.String:
            o.appendString(uint32(obj.Type().Field(i).Offset), obj.Field(i).String())
        }
    }

    return nil
}

func (o *ObjBytes) appendString(offset uint32, v string) error {
    toOffset := len(o.data)
    o.data = append(o.data, []byte(v)...)
    o.pairsOffset = append(o.pairsOffset, PairOffset{
        From: offset,
        To:   uint32(toOffset),
    })

    return nil
}

func (o *ObjBytes) appendSlice(offset uint32, v reflect.Value) error {
    var buffer []byte
    rv := *(*Value)(unsafe.Pointer(&v))
    bh := (*SliceHeader)(unsafe.Pointer(&buffer))
    bh.Data = (*SliceHeader)(rv.ptr).Data
    bh.Len = int(v.Type().Elem().Size()) * v.Len()
    bh.Cap = bh.Len

    toOffset := len(o.data)
    o.data = append(o.data, buffer...)
    // fmt.Println("output：", o.data)

    o.pairsOffset = append(o.pairsOffset, PairOffset{
        From: offset,
        To:   uint32(toOffset),
    })

    return nil
}

func (o *ObjBytes) joinOverHead() error {
    pairsCount := len(o.pairsOffset)
    header := Header{
        Magic:     Magic,
        PairCount: uint32(pairsCount),
    }

    h := *(*[HeaderSize]byte)(unsafe.Pointer(&header))
    message := h[:]
    overHeadSize := HeaderSize

    // append pairs as byte array
    for i := 0; i < pairsCount; i++ {
        v := *(*uint64)(unsafe.Pointer(&o.pairsOffset[i]))
        b := *(*[unsafe.Sizeof(uint64(0))]byte)(unsafe.Pointer(&v))
        message = append(message, b[:]...)

        overHeadSize += uint64(unsafe.Sizeof(uint64(0)))
    }

    o.data = append(message, o.data...)

    // fix field offset
    for i := 0; i < pairsCount; i++ {
        fmt.Println("pair", i, o.pairsOffset[i], len(o.data))
        from := overHeadSize + uint64(o.pairsOffset[i].From)
        to := overHeadSize + uint64(o.pairsOffset[i].To)
        a := *(*[PtrSize]byte)(unsafe.Pointer(&to))
        copy(message[from:], a[:])
    }

    return nil
}

func (o *ObjBytes) memAlign() {
    if l :=
}
