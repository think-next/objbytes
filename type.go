package objbytes

import "unsafe"

// A Kind represents the specific kind of type that a Type represents.
// The zero Kind is not a valid kind.
type Kind uint

const (
    Invalid Kind = iota
    Bool
    Int
    Int8
    Int16
    Int32
    Int64
    Uint
    Uint8
    Uint16
    Uint32
    Uint64
    Uintptr
    Float32
    Float64
    Complex64
    Complex128
    Array
    Chan
    Func
    Interface
    Map
    Ptr
    Slice
    String
    Struct
    UnsafePointer
)

// tflag is used by an rtype to signal what extra type information is
// available in the memory directly following the rtype value.
//
// tflag values must be kept in sync with copies in:
//	cmd/compile/internal/gc/reflect.go
//	cmd/link/internal/ld/decodesym.go
//	runtime/type.go
type tflag uint8

const (
    // tflagUncommon means that there is a pointer, *uncommonType,
    // just beyond the outer type structure.
    //
    // For example, if t.Kind() == Struct and t.tflag&tflagUncommon != 0,
    // then t has uncommonType data and it can be accessed as:
    //
    //	type tUncommon struct {
    //		structType
    //		u uncommonType
    //	}
    //	u := &(*tUncommon)(unsafe.Pointer(t)).u
    tflagUncommon tflag = 1 << 0

    // tflagExtraStar means the name in the str field has an
    // extraneous '*' prefix. This is because for most types T in
    // a program, the type *T also exists and reusing the str data
    // saves binary size.
    tflagExtraStar tflag = 1 << 1

    // tflagNamed means the type has a name.
    tflagNamed tflag = 1 << 2

    // tflagRegularMemory means that equal and hash functions can treat
    // this type as a single region of t.size bytes.
    tflagRegularMemory tflag = 1 << 3
)

// rtype is the common implementation of most values.
// It is embedded in other struct types.
//
// rtype must be kept in sync with ../runtime/type.go:/^type._type.
type rtype struct {
    size       uintptr
    ptrdata    uintptr // number of bytes in the type that can contain pointers
    hash       uint32  // hash of type; avoids computation in hash tables
    tflag      tflag   // extra type information flags
    align      uint8   // alignment of variable with this type
    fieldAlign uint8   // alignment of struct field with this type
    kind       uint8   // enumeration for C
    // function for comparing objects of this type
    // (ptr to object A, ptr to object B) -> ==?
    equal     func(unsafe.Pointer, unsafe.Pointer) bool
    gcdata    *byte   // garbage collection data
    str       nameOff // string form
    ptrToThis typeOff // type for pointer to this type, may be zero
}

type nameOff int32 // offset to a name
type typeOff int32 // offset to an *rtype
type textOff int32 // offset from top of text section
