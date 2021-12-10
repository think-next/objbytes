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

    case reflect.Ptr:
    case reflect.String:
    case reflect.Slice:
    default:

    }

    return nil, nil
}

func StructMarsh(data interface{}, buffer []byte, pairs []PairOffset) ([]byte, error) {

    t := reflect.TypeOf(data)
    switch t.Kind() {
    case reflect.Struct:

    }

    rv := *(*Value)(unsafe.Pointer(&data))
    bh := (*SliceHeader)(unsafe.Pointer(&buffer))
    bh.Data = uintptr(rv.ptr)
    bh.Len = int(data.Type().Size())
    bh.Cap = bh.Len

    //
    for i := 0; i < data.NumField(); i++ {

        switch data.Field(i).Type().Kind() {
        case reflect.Struct:

        case reflect.Slice:

        case reflect.String:

            // 获取这个字段的offset，通过offset来知道它byte数组中的位置
            srcOffset := data.Type().Field(i).Offset
            destRelativeOffset := len(buffer)
            buffer = append(buffer, []byte(data.Field(i).String())...)

            fmt.Println("buffer ①", buffer)
            // 已经加到了数组中，下来就是要声明一个byte头
            pairs = append(pairs, PairOffset{
                From: uint32(srcOffset),
                To:   uint32(destRelativeOffset),
            })
        }
    }
    /*
       | magic｜         uint32
       | pairsCount｜    uint32
       | paris|          uint64 * pairsCount
    */
    pairsCount := len(pairs)
    header := Header{
        Magic:     1,
        PairCount: uint32(pairsCount),
    }

    // 将 header 转换为 []byte 数组，把数据存放到 byte 的数据段里面了
    b := *(*[HeaderSize]byte)(unsafe.Pointer(&header))

    message := make([]byte, HeaderSize)
    copy(message, b[:])

    // 将 pairs 转换为 uint64 的数组
    for i := 0; i < len(pairs); i++ {
        v := *(*uint64)(unsafe.Pointer(&pairs[i]))

        // 数组的声明必须是一个常量
        b := *(*[unsafe.Sizeof(uint64(0))]byte)(unsafe.Pointer(&v))
        message = append(message, b[:]...)
    }

    headerOffset := len(message)
    fmt.Println("xxxxx", headerOffset)

    // 把数据元素加进去，过程中需要调整指针的位置
    message = append(message, buffer...)
    for i := 0; i < len(pairs); i++ {
        // TODO 应该把 from 和 to 的值也给改了才对
        from := int(pairs[i].From) + headerOffset
        to := int(pairs[i].To) + headerOffset

        newAddr := uintptr(unsafe.Pointer(&message[to]))
        buf := *(*[8]byte)(unsafe.Pointer(&newAddr))
        copy(message[from:], buf[:])
    }

    return message, nil
}
