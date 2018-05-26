package fpGo

type fnObj func(*interface{}) *interface{}

func Compose(fnList ...fnObj) fnObj {
	return func(s *interface{}) *interface{} {
		f := fnList[0]
		nextFnList := fnList[1:len(fnList)]

		if len(fnList) == 1 {
			return f(s)
		}

		return f(Compose(nextFnList...)(s))
	}
}
