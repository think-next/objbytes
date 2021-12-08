package objbytes

import (
    "unsafe"
)

const ptrSize = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const

// Value is the reflection interface to a Go value.
//
// Not all methods apply to all kinds of values. Restrictions,
// if any, are noted in the documentation for each method.
// Use the Kind method to find out the kind of value before
// calling kind-specific methods. Calling a method
// inappropriate to the kind of type causes a run time panic.
//
// The zero Value represents no value.
// Its IsValid method returns false, its Kind method returns Invalid,
// its String method returns "<invalid Value>", and all other methods panic.
// Most functions and methods never return an invalid value.
// If one does, its documentation states the conditions explicitly.
//
// A Value can be used concurrently by multiple goroutines provided that
// the underlying Go value can be used concurrently for the equivalent
// direct operations.
//
// To compare two Values, compare the results of the Interface method.
// Using == on two Values does not compare the underlying values
// they represent.
type Value struct {
    // typ holds the type of the value represented by a Value.
    typ *rtype

    // Pointer-valued data or, if flagIndir is set, pointer to data.
    // Valid when either flagIndir is set or typ.pointers() is true.
    ptr unsafe.Pointer

    // flag holds metadata about the value.
    // The lowest bits are flag bits:
    //	- flagStickyRO: obtained via unexported not embedded field, so read-only
    //	- flagEmbedRO: obtained via unexported embedded field, so read-only
    //	- flagIndir: val holds a pointer to the data
    //	- flagAddr: v.CanAddr is true (implies flagIndir)
    //	- flagMethod: v is a method value.
    // The next five bits give the Kind of the value.
    // This repeats typ.Kind() except for method values.
    // The remaining 23+ bits give a method number for method values.
    // If flag.kind() != Func, code can assume that flagMethod is unset.
    // If ifaceIndir(typ), code can assume that flagIndir is set.
    flag

    // A method value represents a curried method invocation
    // like r.Read for some receiver r. The typ+val+flag bits describe
    // the receiver r, but the flag's Kind bits say Func (methods are
    // functions), and the top bits of the flag give the method number
    // in r's type's method table.
}

type flag uintptr

// StringHeader is the runtime representation of a string.
// It cannot be used safely or portably and its representation may
// change in a later release.
// Moreover, the Data field is not sufficient to guarantee the data
// it references will not be garbage collected, so programs must keep
// a separate, correctly typed pointer to the underlying data.
type StringHeader struct {
    Data uintptr
    Len  int
}

// SliceHeader is the runtime representation of a slice.
// It cannot be used safely or portably and its representation may
// change in a later release.
// Moreover, the Data field is not sufficient to guarantee the data
// it references will not be garbage collected, so programs must keep
// a separate, correctly typed pointer to the underlying data.
type SliceHeader struct {
    Data uintptr
    Len  int
    Cap  int
}
