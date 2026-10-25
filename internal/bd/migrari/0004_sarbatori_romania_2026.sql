-- F07-T03, D014: sărbători legale România 2026, cu date mobile ortodoxe.
-- Surse și precizări: docs/decizii/D014-sarbatori-2026.md.
-- 17 sărbători nominale în 16 date; 1 iunie are două denumiri.
-- Păstrăm orice înregistrare preexistentă pentru aceeași dată.
-- +goose Up

INSERT INTO sarbatori (data, denumire, creat_la, modificat_la) VALUES
    ('2026-01-01', '{"ro":"Anul Nou","en":"New Year''s Day"}', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    ('2026-01-02', '{"ro":"A doua zi de Anul Nou","en":"Second day of New Year"}', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    ('2026-01-06', '{"ro":"Botezul Domnului – Boboteaza","en":"Epiphany"}', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    ('2026-01-07', '{"ro":"Soborul Sfântului Proroc Ioan Botezătorul","en":"Saint John the Baptist"}', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    ('2026-01-24', '{"ro":"Ziua Unirii Principatelor Române","en":"Unification of the Romanian Principalities"}', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    ('2026-04-10', '{"ro":"Vinerea Mare (calendar ortodox)","en":"Good Friday (Orthodox calendar)"}', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    ('2026-04-12', '{"ro":"Prima zi de Paști (calendar ortodox)","en":"Easter Sunday (Orthodox calendar)"}', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    ('2026-04-13', '{"ro":"A doua zi de Paști (calendar ortodox)","en":"Easter Monday (Orthodox calendar)"}', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    ('2026-05-01', '{"ro":"Ziua Muncii","en":"Labour Day"}', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    ('2026-05-31', '{"ro":"Prima zi de Rusalii (calendar ortodox)","en":"Pentecost Sunday (Orthodox calendar)"}', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    ('2026-06-01', '{"ro":"Ziua Copilului / A doua zi de Rusalii (calendar ortodox)","en":"Children''s Day / Pentecost Monday (Orthodox calendar)"}', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    ('2026-08-15', '{"ro":"Adormirea Maicii Domnului","en":"Dormition of the Mother of God"}', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    ('2026-11-30', '{"ro":"Sfântul Apostol Andrei","en":"Saint Andrew"}', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    ('2026-12-01', '{"ro":"Ziua Națională a României","en":"Romanian National Day"}', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    ('2026-12-25', '{"ro":"Prima zi de Crăciun","en":"Christmas Day"}', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    ('2026-12-26', '{"ro":"A doua zi de Crăciun","en":"Second day of Christmas"}', strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
ON CONFLICT (data) DO NOTHING;
