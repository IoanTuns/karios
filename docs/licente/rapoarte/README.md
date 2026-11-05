# Licențe pentru generarea PDF

Fișierele din acest folder sunt copiile notificărilor pentru versiunile folosite de F09-T07. Se includ în materialele care însoțesc distribuția aplicației; fonturile sunt incluse în executabil/PDF, fără dependență de fonturi instalate în Windows.

| Componentă | Versiune | Notificare |
| --- | --- | --- |
| maroto v2 | 2.4.4 | maroto-v2.4.4.txt |
| golang.org/x/image | 0.47.0 | go-image-v0.47.0.txt |
| Fonturi Go Regular/Bold, Bigelow & Holmes | din x/image 0.47.0 | fonturi-go.txt |
| pdfcpu | 0.16.3 | pdfcpu-v0.16.3.txt |
| gofpdf | 1.4.3 | gofpdf-v1.4.3.txt |
| go-tree | 1.1.0 | go-tree-v1.1.0.txt |
| barcode | 1.1.0 | barcode-v1.1.0.txt |
| uax29/v2 | 2.7.0 | uax29-v2.7.0.txt |
| tiff | 1.0.7 | tiff-v1.0.7.txt |
| stringprep | 1.0.4 | stringprep-v1.0.4.txt |
| go-runewidth | 0.0.31 | runewidth-v0.0.31.txt |
| pdfme (ui, generator, common, schemas) si pachetele npm de productie ale interfetei | 6.2.8 | npm-interfata-pdfme-v6.2.8.txt |

Sursele sunt modulele descărcate prin Go proxy la versiunile exacte din go.mod/go.sum. Copierea acestor notificări nu introduce alte biblioteci sau programe în aplicație.

F09-T13/T14: gofpdf1.4.3 este reutilizat direct pentru pozitionarea libera; licenta MIT reverificata din modulul local, notificarea existenta pastrata. Nu se adauga biblioteca noua.

F09-T20 (D087): pdfme 6.2.8 si dependentele lui npm de productie, toate cu licente permisive (MIT, Apache-2.0, BSD-3-Clause, ISC, 0BSD, MIT si Zlib, MPL-2.0 sau Apache-2.0). Notificarile sunt generate din package-lock.json si node_modules; pachetele fara fisier propriu de licenta sunt marcate si trimit la textul standard.
