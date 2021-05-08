package main

// simplize the struct, just use Node and D struct
type D struct {
	Node
	Auth     int
	d        int
	Cl       int
	Cv       float64
	Con      int
	Unvalid  int
	isDelete bool
	isGood   bool
	Address  string
	fmm      string
}

// alternate struct for dcml algorithm
type Alter struct {
	D
	isAlter bool
}
