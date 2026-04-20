package main

type Fixtures struct {
	Path         string   `json:"path"`
	Main         string   `json:"main"`
	Port         int      `json:"port"`
	Signal_types []string `json:"signal_types"`
}

var Tests map[int]string = map[int]string{
	1: "PackageTest",
	2: "PortTest",
	3: "SignalTest",
}

var FileMap map[string]string = map[string]string{
	"dir":  "fixtures",
	"json": "fixtures.json",
}
