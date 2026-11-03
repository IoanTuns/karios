WITH lista AS (
  SELECT rowid AS id_rand, *
  FROM garzi
  WHERE data >= :de_la AND data < :pana_la
)
SELECT
  (SELECT COUNT(*) FROM lista b
   WHERE (b.inceput, b.angajat, b.marca, b.id_rand)
      <= (g.inceput, g.angajat, g.marca, g.id_rand)) AS nr_crt,
  g.angajat AS nume,
  g.marca,
  strftime('%d.%m.%Y', g.data) AS data,
  substr(g.inceput, 12, 5) || ' - ' || substr(g.sfarsit, 12, 5)
    || CASE WHEN substr(g.sfarsit, 1, 10) <> g.data
       THEN ' (' || strftime('%d.%m.%Y', g.sfarsit) || ')'
       ELSE '' END AS ora
FROM lista g
ORDER BY g.inceput, g.angajat, g.marca, g.id_rand;
