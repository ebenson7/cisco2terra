package main

type AccessList struct {
	listName         string
	listAction       string
	objectGroupName  string
	sourceGroup      string
	destinationGroup string
	sourcePort       string
	destinationPort  string
}

//WIP