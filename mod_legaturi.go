//go:build bindings

package main

// generareLegaturi: Wails compileaza aplicatia cu tag-ul `bindings` si o ruleaza doar ca sa
// genereze frontend/wailsjs (I020, F02-T33). In acest mod nu se pregateste folderul de date si
// nu se deschide baza: generarea nu atinge niciun karios/ si nu depinde de istoricul migrarilor.
const generareLegaturi = true
