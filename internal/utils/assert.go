package utils

func Assert(b bool) {
	if !b {
		panic("NOT OK")
	}
}
