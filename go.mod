module github.com/hzakher/configure

go 1.17

require (
	github.com/hzakher/yaml/v2 v2.4.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20210912230133-d1bdfacee922 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
)

// replace gopkg.in/yaml.v2 => /home/totta/Documents/workspace/yaml
replace github.com/narisongdev/yaml => /home/totta/Documents/workspace/narisongdev/yaml
replace github.com/hzakher/yaml/v2 => /home/totta/Documents/workspace/hzakher/yaml
