package objbytes

import (
    "reflect"
    "testing"
)

type Case struct {
    M string
    A int
    s int
}

func TestStructMarsh(t *testing.T) {
    c := Case{
        M: "M",
        A: 1,
        s: 20,
    }

    b, _ := StructMarsh(reflect.ValueOf(c))
    t.Log(b)
}
