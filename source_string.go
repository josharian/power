// Code generated by "stringer -type=Source"; DO NOT EDIT.

package power

import "strconv"

const _Source_name = "UnknownBatteryACUPS"

var _Source_index = [...]uint8{0, 7, 14, 16, 19}

func (i Source) String() string {
	if i >= Source(len(_Source_index)-1) {
		return "Source(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Source_name[_Source_index[i]:_Source_index[i+1]]
}
