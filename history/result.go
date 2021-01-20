package history

const ResultStateSuccess = "success"
const ResultStateFailed = "failed"
const ResultStateUnapplied = "unapplied"

var SuccessResult = Result{
	State: ResultStateSuccess,
}

var FailedResult = Result{
	State: ResultStateFailed,
}

var UnappliedResult = Result{
	State: ResultStateUnapplied,
}

type Result struct {
	State string
}

func (r Result) IsSuccess() bool {
	return r.State == ResultStateSuccess
}

func (r Result) IsFailure() bool {
	return r.State == ResultStateFailed
}

func (r Result) IsUnapplied() bool {
	return r.State == ResultStateUnapplied
}
