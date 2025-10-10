-- +goose Up
-- Add Russian translations for date template groups and templates

-- Update date_template_groups with Russian translations
UPDATE date_template_groups SET 
    name_ru = 'Заря цивилизации',
    description_ru = 'Первые города, письменность и рождение сложного общества'
WHERE name_en = 'Dawn of Civilization';

UPDATE date_template_groups SET 
    name_ru = 'Раннединастический период',
    description_ru = 'Египетские династии и шумерские города-государства'
WHERE name_en = 'Early Dynastic Period';

UPDATE date_template_groups SET 
    name_ru = 'Первые империи',
    description_ru = 'Аккадская империя и Древнее царство Египта'
WHERE name_en = 'First Empires';

UPDATE date_template_groups SET 
    name_ru = 'Державы бронзового века',
    description_ru = 'Вавилонская империя и Среднее царство Египта'
WHERE name_en = 'Bronze Age Powers';

UPDATE date_template_groups SET 
    name_ru = 'Имперская эпоха',
    description_ru = 'Новое царство Египта и Хеттская империя'
WHERE name_en = 'Imperial Age';

UPDATE date_template_groups SET 
    name_ru = 'Катастрофа бронзового века',
    description_ru = 'Народы моря и крах цивилизаций'
WHERE name_en = 'Bronze Age Collapse';

UPDATE date_template_groups SET 
    name_ru = 'Древняя Греция',
    description_ru = 'Классический период и основные греческие события'
WHERE name_en = 'Ancient Greece';

UPDATE date_template_groups SET 
    name_ru = 'Римская империя',
    description_ru = 'Основные периоды и события Римской империи'
WHERE name_en = 'Roman Empire';

UPDATE date_template_groups SET 
    name_ru = 'Средневековье',
    description_ru = 'Ключевые события и периоды Средних веков'
WHERE name_en = 'Medieval Period';

UPDATE date_template_groups SET 
    name_ru = 'Ренессанс',
    description_ru = 'Период европейского Возрождения и события'
WHERE name_en = 'Renaissance';

UPDATE date_template_groups SET 
    name_ru = 'Новое время',
    description_ru = 'Основные события с 1500 года'
WHERE name_en = 'Modern Era';

-- Update date_templates with Russian translations

-- Dawn of Civilization templates
UPDATE date_templates SET name_ru = 'Первые города', description_ru = 'Древнейшие городские поселения в Месопотамии'
WHERE name_en = 'First Cities';

UPDATE date_templates SET name_ru = 'Изобретение письма', description_ru = 'Развитие клинописи и иероглифов'
WHERE name_en = 'Invention of Writing';

UPDATE date_templates SET name_ru = 'Объединение Египта', description_ru = 'Нармер объединяет Верхний и Нижний Египет'
WHERE name_en = 'Egyptian Unification';

-- Early Dynastic Period templates
UPDATE date_templates SET name_ru = 'Ранние династии Египта', description_ru = 'Первая и Вторая династии'
WHERE name_en = 'Egyptian Early Dynasties';

UPDATE date_templates SET name_ru = 'Шумерские города-государства', description_ru = 'Независимые месопотамские города'
WHERE name_en = 'Sumerian City-States';

UPDATE date_templates SET name_ru = 'Эпоха пирамид', description_ru = 'Третья и Четвертая династии Египта'
WHERE name_en = 'Age of Pyramids';

UPDATE date_templates SET name_ru = 'Эра Великой пирамиды', description_ru = 'Строительство пирамид в Гизе'
WHERE name_en = 'Great Pyramid Era';

-- First Empires templates
UPDATE date_templates SET name_ru = 'Аккадская империя', description_ru = 'Саргон завоевывает Месопотамию'
WHERE name_en = 'Akkadian Empire';

UPDATE date_templates SET name_ru = 'Древнее царство Египта', description_ru = 'Расцвет строительства пирамид'
WHERE name_en = 'Egyptian Old Kingdom';

UPDATE date_templates SET name_ru = 'Третья династия Ура', description_ru = 'Неошумерское возрождение'
WHERE name_en = 'Third Dynasty of Ur';

UPDATE date_templates SET name_ru = 'Климатический кризис', description_ru = 'Засуха и крах империй'
WHERE name_en = 'Climate Crisis';

-- Bronze Age Powers templates
UPDATE date_templates SET name_ru = 'Старовавилонский период', description_ru = 'Подъем Вавилона при Хаммурапи'
WHERE name_en = 'Old Babylonian Period';

UPDATE date_templates SET name_ru = 'Среднее царство Египта', description_ru = 'Воссоединенный и процветающий Египет'
WHERE name_en = 'Egyptian Middle Kingdom';

UPDATE date_templates SET name_ru = 'Кодекс Хаммурапи', description_ru = 'Первый всеобъемлющий свод законов'
WHERE name_en = 'Code of Hammurabi';

UPDATE date_templates SET name_ru = 'Период гиксосов', description_ru = 'Иноземные правители контролируют северный Египет'
WHERE name_en = 'Hyksos Period';

-- Imperial Age templates
UPDATE date_templates SET name_ru = 'Новое царство Египта', description_ru = 'Египет как имперская сверхдержава'
WHERE name_en = 'Egyptian New Kingdom';

UPDATE date_templates SET name_ru = 'Хеттская империя', description_ru = 'Анатолийская держава бросает вызов Египту'
WHERE name_en = 'Hittite Empire';

UPDATE date_templates SET name_ru = 'Правление Хатшепсут', description_ru = 'Женщина-фараон и процветание'
WHERE name_en = 'Reign of Hatshepsut';

UPDATE date_templates SET name_ru = 'Походы Тутмоса III', description_ru = 'Египетская военная экспансия'
WHERE name_en = 'Thutmose III Campaigns';

UPDATE date_templates SET name_ru = 'Эра Рамсеса II', description_ru = 'Фараон - великий строитель'
WHERE name_en = 'Ramesses II Era';

UPDATE date_templates SET name_ru = 'Битва при Кадеше', description_ru = 'Крупнейшая битва колесниц в истории'
WHERE name_en = 'Battle of Kadesh';

-- Bronze Age Collapse templates
UPDATE date_templates SET name_ru = 'Нашествия народов моря', description_ru = 'Загадочные захватчики опустошают цивилизации'
WHERE name_en = 'Sea Peoples Invasions';

UPDATE date_templates SET name_ru = 'Падение Хеттской империи', description_ru = 'Конец анатолийской сверхдержавы'
WHERE name_en = 'Fall of Hittite Empire';

UPDATE date_templates SET name_ru = 'Катастрофа бронзового века', description_ru = 'Цивилизационный кризис в Восточном Средиземноморье'
WHERE name_en = 'Bronze Age Collapse';

UPDATE date_templates SET name_ru = 'Конец Нового царства', description_ru = 'Египет теряет имперскую власть'
WHERE name_en = 'End of New Kingdom';

UPDATE date_templates SET name_ru = 'Подъем железного века', description_ru = 'Переход к новым технологиям и державам'
WHERE name_en = 'Rise of Iron Age';

-- Ancient Greece templates
UPDATE date_templates SET name_ru = 'Архаический период', description_ru = 'Ранняя греческая цивилизация'
WHERE name_en = 'Archaic Period';

UPDATE date_templates SET name_ru = 'Классический период', description_ru = 'Золотой век Афин'
WHERE name_en = 'Classical Period';

UPDATE date_templates SET name_ru = 'Греко-персидские войны', description_ru = 'Греки против Персидской империи'
WHERE name_en = 'Persian Wars';

UPDATE date_templates SET name_ru = 'Пелопоннесская война', description_ru = 'Афины против Спарты'
WHERE name_en = 'Peloponnesian War';

UPDATE date_templates SET name_ru = 'Эпоха Александра', description_ru = 'Завоевания Александра Великого'
WHERE name_en = 'Age of Alexander';

UPDATE date_templates SET name_ru = 'Эллинистический период', description_ru = 'Распространение греческой культуры'
WHERE name_en = 'Hellenistic Period';

-- Roman Empire templates
UPDATE date_templates SET name_ru = 'Римская республика', description_ru = 'От основания до Юлия Цезаря'
WHERE name_en = 'Roman Republic';

UPDATE date_templates SET name_ru = 'Правление Юлия Цезаря', description_ru = 'Диктатура Юлия Цезаря'
WHERE name_en = 'Rule of Julius Caesar';

UPDATE date_templates SET name_ru = 'Правление Августа', description_ru = 'Первый римский император Октавиан Август'
WHERE name_en = 'Rule of Augustus';

UPDATE date_templates SET name_ru = 'Династия Юлиев-Клавдиев', description_ru = 'От Августа до Нерона'
WHERE name_en = 'Julio-Claudian Dynasty';

UPDATE date_templates SET name_ru = 'Династия Флавиев', description_ru = 'Веспасиан, Тит и Домициан'
WHERE name_en = 'Flavian Dynasty';

UPDATE date_templates SET name_ru = 'Династия Антонинов', description_ru = 'Расцвет Римской империи'
WHERE name_en = 'Antonine Dynasty';

UPDATE date_templates SET name_ru = 'Кризис III века', description_ru = 'Период военной анархии'
WHERE name_en = 'Crisis of Third Century';

UPDATE date_templates SET name_ru = 'Поздняя Римская империя', description_ru = 'От Диоклетиана до падения'
WHERE name_en = 'Late Roman Empire';

-- Medieval Period templates
UPDATE date_templates SET name_ru = 'Раннее Средневековье', description_ru = 'От падения Рима до Карла Великого'
WHERE name_en = 'Early Middle Ages';

UPDATE date_templates SET name_ru = 'Каролингская империя', description_ru = 'Карл Великий и его преемники'
WHERE name_en = 'Carolingian Empire';

UPDATE date_templates SET name_ru = 'Высокое Средневековье', description_ru = 'Крестовые походы и средневековое процветание'
WHERE name_en = 'High Middle Ages';

UPDATE date_templates SET name_ru = 'Первый крестовый поход', description_ru = 'Отвоевание Иерусалима'
WHERE name_en = 'First Crusade';

UPDATE date_templates SET name_ru = 'Третий крестовый поход', description_ru = 'Ричард Львиное Сердце против Саладина'
WHERE name_en = 'Third Crusade';

UPDATE date_templates SET name_ru = 'Позднее Средневековье', description_ru = 'От Черной смерти до Ренессанса'
WHERE name_en = 'Late Middle Ages';

UPDATE date_templates SET name_ru = 'Столетняя война', description_ru = 'Англия против Франции'
WHERE name_en = 'Hundred Years War';

UPDATE date_templates SET name_ru = 'Черная смерть', description_ru = 'Пандемия бубонной чумы'
WHERE name_en = 'Black Death';

-- Renaissance templates
UPDATE date_templates SET name_ru = 'Итальянское Возрождение', description_ru = 'Рождение Ренессанса в Италии'
WHERE name_en = 'Italian Renaissance';

UPDATE date_templates SET name_ru = 'Северное Возрождение', description_ru = 'Ренессанс в Северной Европе'
WHERE name_en = 'Northern Renaissance';

UPDATE date_templates SET name_ru = 'Эпоха великих географических открытий', description_ru = 'Открытие Нового Света'
WHERE name_en = 'Age of Exploration';

UPDATE date_templates SET name_ru = 'Протестантская Реформация', description_ru = 'От Лютера до Вестфальского мира'
WHERE name_en = 'Protestant Reformation';

-- Modern Era templates
UPDATE date_templates SET name_ru = 'Эпоха Просвещения', description_ru = 'Разум и научная революция'
WHERE name_en = 'Age of Enlightenment';

UPDATE date_templates SET name_ru = 'Французская революция', description_ru = 'Революция и Наполеоновские войны'
WHERE name_en = 'French Revolution';

UPDATE date_templates SET name_ru = 'Промышленная революция', description_ru = 'Эпоха пара и модернизации'
WHERE name_en = 'Industrial Revolution';

UPDATE date_templates SET name_ru = 'Первая мировая война', description_ru = 'Великая война'
WHERE name_en = 'World War I';

UPDATE date_templates SET name_ru = 'Вторая мировая война', description_ru = 'Глобальный конфликт'
WHERE name_en = 'World War II';

-- +goose Down
-- No need to rollback translations
