-- Migration: Add Phoenician Mediterranean Empire date template group and templates
-- Run this manually on your local database to add the Phoenician historical period

-- Insert the Phoenician template group
INSERT INTO date_template_groups (name, description, display_order, name_en, name_ru, description_en, description_ru)
VALUES (
  'Phoenician Mediterranean Empire',
  'Phoenician maritime dominance and Carthaginian power',
  5.5,
  'Phoenician Mediterranean Empire',
  'Финикийская средиземноморская империя',
  'Phoenician maritime dominance and Carthaginian power',
  'Финикийское морское господство и могущество Карфагена'
);

-- Insert the Phoenician period templates
-- Note: group_id will be auto-assigned based on the insert above
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order, name_en, name_ru, description_en, description_ru)
SELECT 
  (SELECT id FROM date_template_groups WHERE name_en = 'Phoenician Mediterranean Empire'),
  'Early Phoenician City-States', 
  'Rise of Tyre, Sidon, Byblos as independent trading powers', 
  '1500-01-01', 'BC', '1200-01-01', 'BC', 1, 
  'Early Phoenician City-States', 'Ранние финикийские города-государства', 
  'Rise of Tyre, Sidon, Byblos as independent trading powers', 'Возвышение Тира, Сидона, Библа как независимых торговых держав'
UNION ALL
SELECT 
  (SELECT id FROM date_template_groups WHERE name_en = 'Phoenician Mediterranean Empire'),
  'Phoenician Renaissance', 
  'Maritime expansion after Bronze Age collapse', 
  '1200-01-01', 'BC', '814-01-01', 'BC', 2,
  'Phoenician Renaissance', 'Финикийский Ренессанс',
  'Maritime expansion after Bronze Age collapse', 'Морская экспансия после катастрофы бронзового века'
UNION ALL
SELECT 
  (SELECT id FROM date_template_groups WHERE name_en = 'Phoenician Mediterranean Empire'),
  'Age of Colonization', 
  'Founding of Mediterranean colonies and trade networks', 
  '1100-01-01', 'BC', '650-01-01', 'BC', 3,
  'Age of Colonization', 'Эпоха колонизации',
  'Founding of Mediterranean colonies and trade networks', 'Основание средиземноморских колоний и торговых сетей'
UNION ALL
SELECT 
  (SELECT id FROM date_template_groups WHERE name_en = 'Phoenician Mediterranean Empire'),
  'Carthaginian Ascendancy', 
  'Rise of Carthage as independent Mediterranean superpower', 
  '650-01-01', 'BC', '264-01-01', 'BC', 4,
  'Carthaginian Ascendancy', 'Возвышение Карфагена',
  'Rise of Carthage as independent Mediterranean superpower', 'Возвышение Карфагена как независимой средиземноморской сверхдержавы'
UNION ALL
SELECT 
  (SELECT id FROM date_template_groups WHERE name_en = 'Phoenician Mediterranean Empire'),
  'Persian Period', 
  'Phoenician homeland under Persian Empire rule', 
  '539-01-01', 'BC', '332-01-01', 'BC', 5,
  'Persian Period', 'Персидский период',
  'Phoenician homeland under Persian Empire rule', 'Финикийская родина под властью Персидской империи'
UNION ALL
SELECT 
  (SELECT id FROM date_template_groups WHERE name_en = 'Phoenician Mediterranean Empire'),
  'Punic Wars Era', 
  'Carthage vs Rome struggle for Mediterranean dominance', 
  '264-01-01', 'BC', '146-01-01', 'BC', 6,
  'Punic Wars Era', 'Эпоха Пунических войн',
  'Carthage vs Rome struggle for Mediterranean dominance', 'Борьба Карфагена и Рима за господство в Средиземноморье'
UNION ALL
SELECT 
  (SELECT id FROM date_template_groups WHERE name_en = 'Phoenician Mediterranean Empire'),
  'Hannibalic War', 
  'Hannibal''s legendary campaign against Rome', 
  '218-01-01', 'BC', '201-01-01', 'BC', 7,
  'Hannibalic War', 'Война Ганнибала',
  'Hannibal''s legendary campaign against Rome', 'Легендарная кампания Ганнибала против Рима';
