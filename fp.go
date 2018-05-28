package fpGo

import "sync"

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

func YieldIO(target MonadIODef) *interface{} {
	var result *interface{} = nil

	var wg sync.WaitGroup
	wg.Add(1)
	target.SubscribeOn(nil).Subscribe(Subscription{
		OnNext: func(in *interface{}) {
			result = in
			wg.Done()
		},
	})
	wg.Wait()

	return result
}
