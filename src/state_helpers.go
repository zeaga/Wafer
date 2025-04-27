package main

// Organization:
// pop/push
// value/string/float/bool
// 1/2/3

func (state *EvalState) pop1v() (a Value, ok bool) {
	a, ok = state.values.Pop()
	return
}

func (state *EvalState) pop2v() (a, b Value, ok bool) {
	b, bk := state.values.Pop()
	a, ak := state.values.Pop()
	ok = ak && bk
	return
}

func (state *EvalState) pop3v() (a, b, c Value, ok bool) {
	c, ck := state.values.Pop()
	b, bk := state.values.Pop()
	a, ak := state.values.Pop()
	ok = ak && bk && ck
	return
}

func (state *EvalState) pop1s() (a string, ok bool) {
	av, ok := state.pop1v()
	ok = ok && (av.kind == ValueText)
	if ok {
		a = av.text
	}
	return
}

func (state *EvalState) pop2s() (a, b string, ok bool) {
	av, bv, ok := state.pop2v()
	ok = ok && av.kind == ValueText && bv.kind == ValueText
	if ok {
		a, b = av.text, bv.text
	}
	return
}

func (state *EvalState) pop3s() (a, b, c string, ok bool) {
	av, bv, cv, ok := state.pop3v()
	ok = ok && av.kind == ValueText && bv.kind == ValueText && cv.kind == ValueText
	if ok {
		a, b, c = av.text, bv.text, cv.text
	}
	return
}

func (state *EvalState) pop1f() (a float64, ok bool) {
	av, ok := state.pop1v()
	ok = ok && (av.kind == ValueNumber)
	if ok {
		a = av.number
	}
	return
}

func (state *EvalState) pop2f() (a, b float64, ok bool) {
	av, bv, ok := state.pop2v()
	ok = ok && av.kind == ValueNumber && bv.kind == ValueNumber
	if ok {
		a, b = av.number, bv.number
	}
	return
}

func (state *EvalState) pop3f() (a, b, c float64, ok bool) {
	av, bv, cv, ok := state.pop3v()
	ok = ok && av.kind == ValueNumber && bv.kind == ValueNumber && cv.kind == ValueNumber
	if ok {
		a, b, c = av.number, bv.number, cv.number
	}
	return
}

func (state *EvalState) pop1b() (a bool, ok bool) {
	af, ok := state.pop1f()
	if ok {
		a = floatToBool(af)
	}
	return
}

func (state *EvalState) pop2b() (a, b, ok bool) {
	af, bf, ok := state.pop2f()
	if ok {
		a, b = floatToBool(af), floatToBool(bf)
	}
	return
}

func (state *EvalState) pop3b() (a, b, c, ok bool) {
	af, bf, cf, ok := state.pop3f()
	if ok {
		a, b, c = floatToBool(af), floatToBool(bf), floatToBool(cf)
	}
	return
}

func (state *EvalState) push1v(a Value) bool {
	state.values.Push(a)
	return true
}

func (state *EvalState) push2v(a, b Value) bool {
	state.values.Push(a)
	state.values.Push(b)
	return true
}

func (state *EvalState) push3v(a, b, c Value) bool {
	state.values.Push(a)
	state.values.Push(b)
	state.values.Push(c)
	return true
}

func (state *EvalState) push1s(a string) bool {
	return state.push1v(Value{kind: ValueText, text: a})
}

func (state *EvalState) push2s(a, b string) bool {
	return state.push2v(
		Value{kind: ValueText, text: a},
		Value{kind: ValueText, text: b},
	)
}

func (state *EvalState) push3s(a, b, c string) bool {
	return state.push3v(
		Value{kind: ValueText, text: a},
		Value{kind: ValueText, text: b},
		Value{kind: ValueText, text: c},
	)
}

func (state *EvalState) push1f(a float64) bool {
	return state.push1v(Value{kind: ValueNumber, number: a})
}

func (state *EvalState) push2f(a, b float64) bool {
	return state.push2v(
		Value{kind: ValueNumber, number: a},
		Value{kind: ValueNumber, number: b},
	)
}

func (state *EvalState) push3f(a, b, c float64) bool {
	return state.push3v(
		Value{kind: ValueNumber, number: a},
		Value{kind: ValueNumber, number: b},
		Value{kind: ValueNumber, number: c},
	)
}

func (state *EvalState) push1b(a bool) bool {
	return state.push1f(boolToFloat(a))
}

func (state *EvalState) push2b(a, b bool) bool {
	return state.push2f(boolToFloat(a), boolToFloat(b))
}

func (state *EvalState) push3b(a, b, c bool) bool {
	return state.push3f(boolToFloat(a), boolToFloat(b), boolToFloat(c))
}
