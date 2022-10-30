package main

const (
	Candidate int = 0
	Follower  int = 1
	Leader    int = 2
)

type Server struct {
	state int
}

func ServerConstructor() Server {
	var s Server
	s.state = Candidate
	return s
}
