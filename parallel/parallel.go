package parallel

import "github.com/jonbodner/channel-talk-2/queue"

type Evaluator func(interface{}) (interface{}, error)

type response struct {
	result interface{}
	err error
}

func FanOut(evaluators []Evaluator, data interface{}) ([]interface{}, []error) {
	responses := launch(evaluators, data)
	out, errs := gather(responses, len(evaluators))
	return out, errs
}

func launch(evaluators []Evaluator, data interface{}) queue.Queue {
	responses := queue.MakeInfiniteQueue()
	for _, v := range evaluators {
		go func(e Evaluator) {
			result, err := e(data)
			responses.Put(response{result, err})
		}(v)
	}
	return responses
}

func gather(responses queue.Queue, count int) ([]interface{}, []error) {
	out := make([]interface{}, 0, count)
	errs := make([]error, 0, count)
	for i := 0; i < count; i++ {
		r, _ := responses.Get()
		resp := r.(response)
		if resp.err != nil {
			errs = append(errs,resp.err)
		} else {
			out = append(out, resp.result)
		}
	}
	responses.Close()
	return out, errs
}

