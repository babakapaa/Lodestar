package study

import "net/http"

func dayRange(r *http.Request) (Day, Day, error) {
	q := r.URL.Query()
	from, err := ParseDay(q.Get("from"))
	if err != nil {
		return Day{}, Day{}, err
	}
	to, err := ParseDay(q.Get("to"))
	if err != nil {
		return Day{}, Day{}, err
	}
	return from, to, nil
}
