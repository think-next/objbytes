package objbytes

import (
    "reflect"
    "unsafe"
)

func Unmarshal(b []byte, v interface{}) error {
    // fmt.Println("buffer length", len(b))
    header := *(*Header)(unsafe.Pointer(&b[0]))

    var pairs []uint64
    temp := (*SliceHeader)(unsafe.Pointer(&pairs))
    temp.Data = uintptr(unsafe.Pointer(&b[unsafe.Sizeof(emptyHeader)]))
    temp.Len = int(header.PairCount)
    temp.Cap = int(header.PairCount)

    for _, v := range pairs {

        offset := *(*PairOffset)(unsafe.Pointer(&v))
        to := *(*[PtrSize]byte)(unsafe.Pointer(&b[offset.To]))
        copy(b[offset.From:], to[:])
    }

    elem := reflect.ValueOf(v).Elem()
    elemStruct := (*Value)(unsafe.Pointer(&elem))
    elemStruct.ptr = unsafe.Pointer(&b[int(unsafe.Sizeof(emptyHeader))+int(unsafe.Sizeof(uint64(0)))*int(header.PairCount)])

    reflect.ValueOf(v).Elem().Set(elem)
    return nil
}
