// Code generated by "stringer -type=PropType -output=prop_type_string.go -linecomment"; DO NOT EDIT.

package gosln

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PTBool-1]
	_ = x[PTInt-2]
	_ = x[PTInt8-3]
	_ = x[PTInt16-4]
	_ = x[PTInt32-5]
	_ = x[PTInt64-6]
	_ = x[PTUint-7]
	_ = x[PTUint8-8]
	_ = x[PTUint16-9]
	_ = x[PTUint32-10]
	_ = x[PTUint64-11]
	_ = x[PTUintptr-12]
	_ = x[PTFloat32-13]
	_ = x[PTFloat64-14]
	_ = x[PTComplex64-15]
	_ = x[PTComplex128-16]
	_ = x[PTBytes-17]
	_ = x[PTString-18]
	_ = x[PTTime-19]
	_ = x[PTDate-20]
	_ = x[maxPropType-21]
}

const _PropType_name = "boolintint8int16int32int64uintuint8uint16uint32uint64uintptrfloat32float64complex64complex128[]bytestringtime.Timegosln.DatePropType(21)"

var _PropType_index = [...]uint8{0, 4, 7, 11, 16, 21, 26, 30, 35, 41, 47, 53, 60, 67, 74, 83, 93, 99, 105, 114, 124, 136}

func (i PropType) String() string {
	i -= 1
	if i < 0 || i >= PropType(len(_PropType_index)-1) {
		return "PropType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _PropType_name[_PropType_index[i]:_PropType_index[i+1]]
}
