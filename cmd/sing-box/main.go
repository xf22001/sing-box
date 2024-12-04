//go:build !generate

package boxmain

import "github.com/sagernet/sing-box/log"

func main() {
	if err := mainCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
