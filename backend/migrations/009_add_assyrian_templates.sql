-- +goose Up
-- Add Assyrian Empire date templates to First Empires group

-- Insert the Assyrian period templates into the First Empires group
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order, name_en, name_ru, description_en, description_ru)
SELECT 
  (SELECT id FROM date_template_groups WHERE name_en = 'First Empires'),
  'Old Assyrian Trade Empire', 
  'Assyrian merchant colonies and trade networks across Anatolia', 
  '2025-01-01'::date, 'BC', '1364-01-01'::date, 'BC', 2.1, 
  'Old Assyrian Trade Empire', 'Староассирийская торговая империя', 
  'Assyrian merchant colonies and trade networks across Anatolia', 'Ассирийские торговые колонии и торговые сети в Анатолии'
UNION ALL
SELECT 
  (SELECT id FROM date_template_groups WHERE name_en = 'First Empires'),
  'Middle Assyrian Rise', 
  'Territorial expansion and military conquests of Mitanni and neighbors', 
  '1363-01-01'::date, 'BC', '912-01-01'::date, 'BC', 2.2,
  'Middle Assyrian Rise', 'Возвышение Среднеассирийской империи',
  'Territorial expansion and military conquests of Mitanni and neighbors', 'Территориальная экспансия и военные завоевания Митанни и соседей'
UNION ALL
SELECT 
  (SELECT id FROM date_template_groups WHERE name_en = 'First Empires'),
  'Neo-Assyrian Early Expansion', 
  'Revival and renewed expansion toward Levant and Anatolia', 
  '911-01-01'::date, 'BC', '745-01-01'::date, 'BC', 2.3,
  'Neo-Assyrian Early Expansion', 'Ранняя неоассирийская экспансия',
  'Revival and renewed expansion toward Levant and Anatolia', 'Возрождение и возобновленная экспансия в Левант и Анатолию'
UNION ALL
SELECT 
  (SELECT id FROM date_template_groups WHERE name_en = 'First Empires'),
  'Neo-Assyrian Golden Age', 
  'Peak of Assyrian power under great kings Tiglath-Pileser III through Ashurbanipal', 
  '745-01-01'::date, 'BC', '627-01-01'::date, 'BC', 2.4,
  'Neo-Assyrian Golden Age', 'Золотой век Неоассирийской империи',
  'Peak of Assyrian power under great kings Tiglath-Pileser III through Ashurbanipal', 'Пик ассирийской мощи при великих царях от Тиглатпаласара III до Ашшурбанипала'
UNION ALL
SELECT 
  (SELECT id FROM date_template_groups WHERE name_en = 'First Empires'),
  'Fall of Assyria', 
  'Rapid collapse and destruction of the Neo-Assyrian Empire', 
  '627-01-01'::date, 'BC', '609-01-01'::date, 'BC', 2.5,
  'Fall of Assyria', 'Падение Ассирии',
  'Rapid collapse and destruction of the Neo-Assyrian Empire', 'Быстрый крах и разрушение Неоассирийской империи';

-- +goose Down
-- Remove Assyrian templates from First Empires group
DELETE FROM date_templates 
WHERE group_id = (SELECT id FROM date_template_groups WHERE name_en = 'First Empires')
AND name_en IN (
  'Old Assyrian Trade Empire',
  'Middle Assyrian Rise',
  'Neo-Assyrian Early Expansion',
  'Neo-Assyrian Golden Age',
  'Fall of Assyria'
);
