// Regula parolei (D045 punctul 2): cel puțin 8 caractere, cu cel puțin o literă și o cifră.
// Aceeași regulă e verificată și în Go (conturi.ValideazaParola).
export function parolaValida(p: string): boolean {
    return [...p].length >= 8 && /\p{L}/u.test(p) && /\p{N}/u.test(p);
}
