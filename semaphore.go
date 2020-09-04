package main

type semaphore chan bool

func (s semaphore) Acquire() {
	e := false
	s <- e
}

func (s semaphore) Release() {
	<-s
}
