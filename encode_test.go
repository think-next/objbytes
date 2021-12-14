package objbytes

import (
    "testing"
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

    cc := &Case{}
    Unmarshal(b, cc)
    t.Log(cc)
}
