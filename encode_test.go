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
        A: 1,
    }

    // TODO 不能执行两次 reflect.ValueOf
    b, _ := Marshal(c)
    t.Log(b)

    header := *(*Header)(unsafe.Pointer(&b[0]))
    fmt.Println(header.Magic)
    fmt.Println(header.PairCount)

    // 反射不出来数组结构，因为数组的长度需要是一个常量
    var offsets []uint64
    temp := (*SliceHeader)(unsafe.Pointer(&offsets))

    temp.Data = uintptr(unsafe.Pointer(&b[unsafe.Sizeof(emptyHeader)]))
    temp.Len = int(header.PairCount)
    temp.Cap = int(header.PairCount)
    fmt.Println(offsets)
    for _, v := range offsets {

        // TODO 这里存在问题
        offset := *(*PairOffset)(unsafe.Pointer(&v))
        fmt.Println(offset.From)
        fmt.Println(offset.To)
    }

    // 结构体映射
    fmt.Println(int(unsafe.Sizeof(emptyHeader)))
    fmt.Println(int(unsafe.Sizeof(uint64(0))) * int(header.PairCount))
    // obj := (*Value)(unsafe.Pointer(&b[int(unsafe.Sizeof(emptyHeader))+int(unsafe.Sizeof(uint64(0)))*int(header.PairCount)]))

    ca := &Case{}
    elem := reflect.ValueOf(ca).Elem()
    elemStruct := (*Value)(unsafe.Pointer(&elem))
    elemStruct.ptr = unsafe.Pointer(&b[int(unsafe.Sizeof(emptyHeader))+int(unsafe.Sizeof(uint64(0)))*int(header.PairCount)])

    reflect.ValueOf(ca).Elem().Set(elem)
    fmt.Println(ca)

    // aa := []byte{60, 111, 98, 106, 98, 121, 116, 101, 115, 46, 67, 97, 115, 101, 32, 86, 97, 108, 117, 101, 62}
    // fmt.Println("奇怪", string(aa))
}
