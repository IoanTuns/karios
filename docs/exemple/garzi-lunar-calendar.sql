SELECT
  CAST(strftime('%d', c.data) AS INTEGER) AS data,
  CASE strftime('%w', c.data)
    WHEN '0' THEN 'Duminică'
    WHEN '1' THEN 'Luni'
    WHEN '2' THEN 'Marți'
    WHEN '3' THEN 'Miercuri'
    WHEN '4' THEN 'Joi'
    WHEN '5' THEN 'Vineri'
    WHEN '6' THEN 'Sâmbătă'
  END AS ziua,
  COALESCE(g.angajat, '') AS nume_si_prenume,
  CASE WHEN g.inceput IS NULL THEN '' ELSE substr(g.inceput, 12, 5) || ' - ' || substr(g.sfarsit, 12, 5)
    || CASE WHEN substr(g.sfarsit, 1, 10) <> c.data
       THEN ' (' || strftime('%d.%m.%Y', g.sfarsit) || ')'
       ELSE '' END END AS observatii,
  c.data AS data_calendar
FROM calendar c
LEFT JOIN garzi g ON g.data = c.data
WHERE c.data >= :de_la AND c.data < :pana_la
ORDER BY c.data, g.inceput, g.angajat, g.marca, g.rowid;
