package main

type Kind int

const (
	KindFile Kind = iota
	KindPage
	KindGroup
	KindDir
	KindFunc
	KindMacro
)
