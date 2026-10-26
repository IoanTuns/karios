//go:build !windows

package main

import "log"

func arataEroare(titlu, text string) {
	log.Printf("%s: %s", titlu, text)
}
