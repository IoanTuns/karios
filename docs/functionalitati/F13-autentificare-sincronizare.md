# F13 · Autentificare și sincronizare

Etapă: Ulterior · Stare: amânat · Depinde de: F06

## Scop

Utilizatorul se autentifică printr-un API GraphQL (de exemplu contul Saleor), iar datele se pot sincroniza. Prima versiune doar se pregătește pentru asta.

## Taskuri

| ID | Task | Stare | Note |
| --- | --- | --- | --- |
| F13-T01 | Pregătire în V1: `current_user()`, adresa API în `config.toml`, coloanele de sincronizare | de făcut | Coloanele intră în F06-T02 |
| F13-T02 | Login cu `tokenCreate`, reînnoire cu `tokenRefresh` (`genqlient`) | amânat | |
| F13-T03 | Tokenuri în Windows Credential Manager (`go-keyring`) | amânat | Excepție de la portabilitate, vezi D002 |
| F13-T04 | Sincronizarea datelor | amânat | |
