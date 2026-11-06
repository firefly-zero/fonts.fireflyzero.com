package generator

type Family struct {
	Name    string
	Author  string
	License string
	URL     string
}

var families = map[string]Family{
	"bios": {
		Name:    "bios",
		Author:  "IBM",
		License: "",
	},
	"cga": {
		Name:    "cga",
		Author:  "IBM",
		License: "",
	},
	"clean": {
		Name:    "clean",
		Author:  "Dale Schumacher",
		License: "MIT",
	},
	"cozette": {
		Name:    "cozette",
		Author:  "Ines",
		License: "MIT",
		URL:     "https://github.com/the-moonwitch/Cozette",
	},
	"creep": {
		Name:    "creep",
		Author:  "romeovs",
		License: "MIT",
		URL:     "https://github.com/romeovs/creep",
	},
	"ega": {
		Name:    "ega",
		Author:  "IBM",
		License: "",
	},
	"haxor": {
		Name:    "haxor",
		Author:  "ucw.cz",
		License: "GPLv2+",
		URL:     "https://github.com/metan-ucw/fonts",
	},
	"helvetica": {
		Name:    "helvetica",
		Author:  "Adobe",
		License: "",
	},
	"ibm-437": {
		Name:    "ibm-437",
		Author:  "IBM",
		License: "",
	},
	"knxt": {
		Name:    "knxt",
		Author:  "Konstantin Bytensky",
		License: "",
	},
	"mda": {
		Name:    "mda",
		Author:  "",
		License: "",
	},
	"peep": {
		Name:    "peep",
		Author:  "Zevv",
		License: "",
	},
	"scientifica": {
		Name:    "scientifica",
		Author:  "romeovs",
		License: "",
	},
	"spleen": {
		Name:    "spleen",
		Author:  "Frederic Cambus",
		License: "BSD 2-Clause",
		URL:     "https://github.com/fcambus/spleen",
	},
	"tahoma": {
		Name:    "tahoma",
		Author:  "Larry Snyder",
		License: "",
	},
	"tom-thumb": {
		Name:    "tom-thumb",
		Author:  "Raccoon",
		License: "MIT",
		URL:     "https://robey.lag.net/2010/01/23/tiny-monospace-font.html",
	},
	"ucs": {
		Name:    "ucs",
		Author:  "",
		License: "public domain",
		URL:     "https://www.cl.cam.ac.uk/~mgk25/ucs-fonts.html",
	},
	"vga": {
		Name:    "vga",
		Author:  "IBM",
		License: "",
	},
	"zx": {
		Name:    "zx-spectrum",
		Author:  "",
		License: "",
	},
}
