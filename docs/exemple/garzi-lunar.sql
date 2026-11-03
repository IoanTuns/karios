SELECT
  CAST(strftime('%d', g.data) AS INTEGER) AS data,
  CASE strftime('%w', g.data)
    WHEN '0' THEN 'Duminică'
    WHEN '1' THEN 'Luni'
    WHEN '2' THEN 'Marți'
    WHEN '3' THEN 'Miercuri'
    WHEN '4' THEN 'Joi'
    WHEN '5' THEN 'Vineri'
    WHEN '6' THEN 'Sâmbătă'
  END AS ziua,
  g.angajat AS nume_si_prenume,
  substr(g.inceput, 12, 5) || ' - ' || substr(g.sfarsit, 12, 5)
    || CASE WHEN substr(g.sfarsit, 1, 10) <> g.data
       THEN ' (' || strftime('%d.%m.%Y', g.sfarsit) || ')'
       ELSE '' END AS observatii
FROM garzi g
WHERE g.data >= :de_la AND g.data < :pana_la
ORDER BY g.inceput, g.angajat, g.marca, g.rowid;
