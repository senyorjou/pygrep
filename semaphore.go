package main

func (s semaphore) Acquire() {
	e := false
	s <- e
}

func (s semaphore) Release() {
	<-s
}
