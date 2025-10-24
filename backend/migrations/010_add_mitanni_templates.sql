-- +goose Up
-- Add Hurrian-Mitanni Kingdom date templates to Bronze Age Powers group

-- Insert Mitanni period templates into the Bronze Age Powers group
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order, name_en, name_ru, description_en, description_ru)
SELECT 
  (SELECT id FROM date_template_groups WHERE name_en = 'Bronze Age Powers'),
  'Early Mitanni Formation', 
  'Rise of Hurrian tribes and early Mitanni state formation', 
  '1600-01-01'::date, 'BC', '1500-01-01'::date, 'BC', 4.5, 
  'Early Mitanni Formation', 'Раннее формирование Митанни', 
  'Rise of Hurrian tribes and early Mitanni state formation', 'Возвышение хурритских племен и раннее формирование государства Митанни'
UNION ALL
SELECT 
  (SELECT id FROM date_template_groups WHERE name_en = 'Bronze Age Powers'),
  'Mitanni Golden Age', 
  'Peak of Mitanni power under Shaushtatar, expansion across northern Mesopotamia and Syria', 
  '1500-01-01'::date, 'BC', '1360-01-01'::date, 'BC', 4.6,
  'Mitanni Golden Age', 'Золотой век Митанни',
  'Peak of Mitanni power under Shaushtatar, expansion across northern Mesopotamia and Syria', 'Пик могущества Митанни при Шауштатаре, экспансия в северной Месопотамии и Сирии'
UNION ALL
SELECT 
  (SELECT id FROM date_template_groups WHERE name_en = 'Bronze Age Powers'),
  'Mitanni-Egyptian Alliance', 
  'Alliance period with Egypt against Hittite threat, diplomatic marriages and peak prosperity', 
  '1400-01-01'::date, 'BC', '1360-01-01'::date, 'BC', 4.7,
  'Mitanni-Egyptian Alliance', 'Митанни-египетский союз',
  'Alliance period with Egypt against Hittite threat, diplomatic marriages and peak prosperity', 'Период союза с Египтом против хеттской угрозы, дипломатические браки и пик процветания'
UNION ALL
SELECT 
  (SELECT id FROM date_template_groups WHERE name_en = 'Bronze Age Powers'),
  'Fall of Mitanni', 
  'Hittite conquest, vassalage under Suppiluliumas I, and final Assyrian annexation', 
  '1360-01-01'::date, 'BC', '1240-01-01'::date, 'BC', 4.8,
  'Fall of Mitanni', 'Падение Митанни',
  'Hittite conquest, vassalage under Suppiluliumas I, and final Assyrian annexation', 'Хеттское завоевание, вассалитет при Суппилулиумасе I и окончательная ассирийская аннексия';

-- +goose Down
-- Remove Mitanni templates from Bronze Age Powers group
DELETE FROM date_templates 
WHERE group_id = (SELECT id FROM date_template_groups WHERE name_en = 'Bronze Age Powers')
AND name_en IN (
  'Early Mitanni Formation',
  'Mitanni Golden Age',
  'Mitanni-Egyptian Alliance',
  'Fall of Mitanni'
);
