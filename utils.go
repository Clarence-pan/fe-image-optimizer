package main

func panicIf(x interface{}) {
	if x != nil {
		panic(x)
	}
}
