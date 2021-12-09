package objbytes

import (
    "fmt"
    "reflect"
    "testing"
    "unsafe"
)

type Case struct {
    M string
    A int
}

func TestStructMarsh(t *testing.T) {
    c := Case{
        M: "付辉",
        A: 2,
    }

    b, _ := StructMarsh(reflect.ValueOf(c), []byte{}, []OffsetInfo{})
    t.Log(b)

    header := *(*HeaderInfo)(unsafe.Pointer(&b[0]))
    fmt.Println(header.Magic)
    fmt.Println(header.ParisCount)

    // 反射不出来数组结构，因为数组的长度需要是一个常量
    var offsets []uint64
    temp := (*SliceHeader)(unsafe.Pointer(&offsets))

    temp.Data = uintptr(unsafe.Pointer(&b[unsafe.Sizeof(emptyHeaderInfo)]))
    temp.Len = int(header.ParisCount)
    temp.Cap = int(header.ParisCount)
    fmt.Println(offsets)
    for _, v := range offsets {

        // TODO 这里存在问题
        offset := *(*OffsetInfo)(unsafe.Pointer(&v))
        fmt.Println(offset.From)
        fmt.Println(offset.To)
    }

    // 结构体映射
    fmt.Println(int(unsafe.Sizeof(emptyHeaderInfo)))
    fmt.Println(int(unsafe.Sizeof(uint64(0))) * int(header.ParisCount))
    // obj := (*Value)(unsafe.Pointer(&b[int(unsafe.Sizeof(emptyHeaderInfo))+int(unsafe.Sizeof(uint64(0)))*int(header.ParisCount)]))

    ca := &Case{}
    elem := reflect.ValueOf(ca).Elem()
    elemStruct := (*Value)(unsafe.Pointer(&elem))
    elemStruct.ptr = unsafe.Pointer(&b[int(unsafe.Sizeof(emptyHeaderInfo))+int(unsafe.Sizeof(uint64(0)))*int(header.ParisCount)])

    reflect.ValueOf(ca).Elem().Set(elem)
    fmt.Println(ca)

    // aa := []byte{60, 111, 98, 106, 98, 121, 116, 101, 115, 46, 67, 97, 115, 101, 32, 86, 97, 108, 117, 101, 62}
    // fmt.Println("奇怪", string(aa))
}
