package full

type private struct {
	a int
	b int
}

type Public struct {
	A int
	B int
}

func pass() private {
	return private{}
}

func pass2() Public {
	return Public{}
}

func pass3() private {
	// check:exhaustive
	return private{
		a: 0,
		b: 0,
	}
}

func pass4() Public {
	// check:exhaustive
	return Public{
		A: 0,
		B: 0,
	}
}

func pass5() {
	// check:exhaustive
	_ = private{
		a: 0,
		b: 0,
	}

	// check:exhaustive
	_ = Public{
		A: 0,
		B: 0,
	}
}

func fail() private {
	// check:exhaustive
	return private{} // want "private is missing fields: a, b"
}

func fail2() Public {
	// check:exhaustive
	return Public{} // want "Public is missing fields: A, B"
}

func fail3() private {
	// check:exhaustive
	return private{ // want "private is missing fields: b"
		a: 0,
	}
}

func fail4() Public {
	// check:exhaustive
	return Public{ // want "Public is missing fields: A"
		B: 0,
	}
}

func fail5() {
	// check:exhaustive
	_ = private{ // want "private is missing fields: a"
		b: 0,
	}

	// check:exhaustive
	_ = Public{ // want "Public is missing fields: B"
		A: 0,
	}
}

func fail6() {
	// check:exhaustive // want "unmatched check:exhaustive comment"
	_ = map[string]string{
		"a": "apple",
		"b": "banana",
	}
}

func fail7() {
	// check:exhaustive // want "unmatched check:exhaustive comment"
}
