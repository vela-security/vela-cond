package cond

import "fmt"

type Section struct {
	err    error
	not    bool
	method op
	raw    string
	keys   []string
	data   []string
}

func (s *Section) invalid(format string, v ...interface{}) {
	s.err = fmt.Errorf(format, v...)
}

func (s *Section) trim(offset *int, n int) {
	for i := *offset; i < n; i++ {
		if ch := s.raw[i]; ch != ' ' {
			*offset = i
			return
		}
	}
}

func (s *Section) Ok() bool {
	return s.err == nil
}

func (s *Section) withA(offset *int, n int) {
	s.trim(offset, n)
	sep := *offset

	for i := *offset; i < n; i++ {
		ch := s.raw[i]
		switch ch {
		case ',':
			s.keys = append(s.keys, s.raw[sep:i])
			sep = i

		case ' ':
			if sep != i {
				s.keys = append(s.keys, s.raw[sep:i])
			}
			*offset = i
			return
		}
	}

}

func (s *Section) withB(offset *int, n int) {
	if !s.Ok() {
		return
	}

	s.trim(offset, n)
	sep := *offset

	if sep+3 > n {
		s.invalid("not found method")
		return
	}

	if s.raw[sep] == '!' {
		s.not = true
		sep++
	}

	em := s.raw[sep : sep+2]
	switch em {
	case "eq":
		s.method = eq
	case "re":
		s.method = re
	case "cn":
		s.method = cn
	case "in":
		s.method = in
	case "lt":
		s.method = lt
	case "le":
		s.method = le
	case "ge":
		s.method = ge
	case "gt":
		s.method = gt

	default:
		s.invalid("invalid method %s", em)
		return
	}

	*offset = sep + 2
}

func (s *Section) withC(offset *int, n int) {
	if !s.Ok() {
		return
	}

	s.trim(offset, n)
	sep := *offset
	for i := *offset; i < n; i++ {
		ch := s.raw[i]

		if ch != ',' {
			continue
		}

		if s.raw[i-1] == '\\' {
			continue
		}

		if s.raw[sep] == ',' {
			s.data = append(s.data, s.raw[sep+1:i])
		} else {
			s.data = append(s.data, s.raw[sep:i])
		}

		sep = i
	}

	//single value
	if sep == *offset {
		s.data = append(s.data, s.raw[sep:])
		return
	}

	//last value
	if sep != n-1 {
		s.data = append(s.data, s.raw[sep+1:])
	}
}

func (s *Section) Match(v string) bool {
	n := len(s.data)
	for i := 0; i < n; i++ {
		if s.method.Do(v, s.data[i]) {
			return true
		}
	}

	return false
}

func (s *Section) Call(ov *option) (bool, error) {
	if ov.peek == nil {
		return false, fmt.Errorf("invalid peek function")
	}

	if !s.Ok() {
		return false, s.err
	}

	n := len(s.keys)
	for i := 0; i < n; i++ {
		if !s.Match(ov.peek(s.keys[i])) {
			continue
		}
		return s.not != true, s.err
	}

	return s.not != false, s.err
}

//Compile
//aaa eq abc,eee,fff => Section{not:false , keys: []string{aaa} , method: eq , data: []string{abc , eee , ff}}
//aaa !eq abc,eee,fff => Section{not:true, keys: []string{aaa} , method: eq , data: []string{abc , eee , ff}}
func Compile(raw string) *Section {
	section := &Section{
		raw:    raw,
		method: oop,
	}

	n := len(raw)
	if n < 6 {
		section.err = fmt.Errorf("too short")
		return section
	}

	offset := 0
	section.withA(&offset, n)
	section.withB(&offset, n)
	section.withC(&offset, n)

	if len(section.data) == 0 {
		section.invalid("not found data")
	}

	return section
}
