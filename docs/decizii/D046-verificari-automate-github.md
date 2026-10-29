# D046 · Verificări automate pe GitHub (CI)

Stare: acceptată · Data: 2026-10-29 · Decisă de: Ioan Tuns (lista de taskuri din 2026-10-29, F02-T13) · Taskuri: F02-T13

## Context

`main` a stat mult la 0.7.0 cu peste 100 de commit-uri neaduse, iar fișiere neformatate au ajuns în repo. Verificările de dinainte de commit (AGENTS.md §5) depind de disciplina fiecărui agent.

## Decizie

- GitHub Actions rulează la fiecare push și pull request, pe `windows-latest` (aplicația e pentru Windows și are cod specific): `npm ci` + `npm run build`, `gofmt -l`, `go vet ./...`, `go test ./...`.
- Se folosesc doar acțiunile oficiale GitHub (`actions/checkout`, `actions/setup-go`, `actions/setup-node`), cu versiunea exactă (§6). Go vine din `go.mod`, Node e fixat la versiunea folosită local.
- Workflow-ul are doar drept de citire (`permissions: contents: read`).

## Variante respinse

- **Doar hook-ul local** — nu acoperă agenții sau clonele care nu-l activează.
- **Runner Linux** — mai rapid, dar nu rulează codul Windows (`x/sys/windows`, dialogurile).

## Consecințe

- Un PR cu teste picate sau cod neformatat se vede înainte de merge.
- Actualizarea acțiunilor și a Node intră în verificarea dependențelor de la fiecare modificare (§6).
