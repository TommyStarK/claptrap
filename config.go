package main

type baseHTTPRequest struct {
	http     string   `yaml:"http"`
	url      string   `yaml:"url"`
	method   string   `yaml:"method"`
	headers  []string `yaml:"headers"`
	params   []string `yaml:"params"`
	queries  []string `yaml:"queries"`
	jsonBody string   `yaml:"json_body"`
}

type config struct {
	path   string `yaml:"path"`
	events struct {
		new    []baseHTTPRequest `yaml:"new"`
		update []baseHTTPRequest `yaml:"update"`
		remove []baseHTTPRequest `yaml:"remove"`
	} `yaml:"events"`
}
