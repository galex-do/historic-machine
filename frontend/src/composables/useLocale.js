import { ref, computed } from 'vue'

// UI Text translations
const UI_TRANSLATIONS = {
  en: {
    // Filters
    filters: 'Filters',
    historicalEvents: 'Historical Events',
    historicalPeriod: 'Historical Period',
    specificPeriod: 'Specific Period:',
    selectPeriodGroup: 'Select Period',
    from: 'From:',
    to: 'To:',
    step: 'Step:',
    defaultPeriod: 'Default (1 AD - Today)',
    customDateRange: 'Custom Date Range',
    selectSpecificPeriod: 'Select specific period...',
    noTemplatesAvailable: 'No templates available',
    loading: 'Loading',
    years: 'years',
    apply: 'Apply',
    
    // Header Navigation
    map: 'Map',
    admin: 'Admin area',
    events: 'Events',
    tags: 'Tags',
    datasets: 'Datasets',
    users: 'Users',
    
    // Authentication
    login: 'Login',
    logout: 'Logout',
    welcome: 'Welcome',
    loggingIn: 'Logging in...',
    
    // Login Modal
    loginToHistoria: 'Login to Historia ex machina',
    username: 'Username:',
    password: 'Password:',
    enterUsername: 'Enter username',
    enterPassword: 'Enter password',
    
    // Access Levels
    superBadge: 'SUPER',
    adminBadge: 'ADMIN',
    editorBadge: 'EDITOR',
    
    // Admin Events Page
    adminEventsTitle: 'Historic Events',
    adminEventsSubtitle: 'Manage and create historic events',
    createNewEvent: 'Create New Event',
    searchEventName: 'Search by name...',
    
    // Admin Tags Page
    adminTagsTitle: 'Tags Management',
    adminTagsSubtitle: 'Manage event tags and categories',
    createNewTag: 'Create New Tag',
    
    // Admin Users Page
    adminUsersTitle: 'Users Management',
    adminUsersSubtitle: 'Manage user accounts and access levels',
    createNewUser: 'Create New User',
    
    // Table Column Names
    columnName: 'Name',
    columnDescription: 'Description',
    columnDate: 'Date',
    columnLocation: 'Location',
    columnType: 'Type',
    columnTags: 'Tags',
    columnActions: 'Actions',
    columnColor: 'Color',
    columnCreated: 'Created',
    columnUsageCount: 'Usage Count',
    columnUsername: 'Username',
    columnEmail: 'Email',
    columnAccessLevel: 'Access Level',
    columnStatus: 'Status',
    columnLastActive: 'Last Active',
    
    // Month Names
    monthJan: 'Jan',
    monthFeb: 'Feb',
    monthMar: 'Mar',
    monthApr: 'Apr',
    monthMay: 'May',
    monthJun: 'Jun',
    monthJul: 'Jul',
    monthAug: 'Aug',
    monthSep: 'Sep',
    monthOct: 'Oct',
    monthNov: 'Nov',
    monthDec: 'Dec',
    
    // Era Labels
    eraBC: 'BC',
    eraAD: 'AD',
    
    // Admin Datasets Page
    adminDatasetsTitle: 'Event Datasets',
    adminDatasetsSubtitle: 'Manage uploaded event datasets. Deleting a dataset will remove all events imported from it.',
    loadingDatasets: 'Loading datasets...',
    tryAgain: 'Try Again',
    importNewDataset: 'Import New Dataset',
    importDatasetDescription: 'Upload a JSON file containing historical events to create a new dataset.',
    chooseDatasetFile: 'Choose Dataset File',
    selectedFile: 'Selected:',
    importing: 'Importing...',
    importDataset: 'Import Dataset',
    clear: 'Clear',
    or: 'or',
    createEmptyDataset: 'Create Empty Dataset',
    
    // Tag Filter
    filteredByTags: 'Filtered by tags',
    clearAllTags: 'Clear all',
    clearAll: 'Clear all',
    remove: 'Remove',
    allTags: 'All Tags',
    searchTags: 'Search tags...',
    noTagsFound: 'No tags found',
    followEvents: 'Follow',
    enableNarrativeFlow: 'Enable narrative flow - connect events chronologically',
    disableNarrativeFlow: 'Disable narrative flow',
    
    // Datasets Table
    columnFilename: 'Filename',
    columnEventCount: 'Event Count',
    columnUploadedBy: 'Uploaded By',
    columnUploadDate: 'Upload Date',
    noDatasetsFound: 'No datasets found',
    importToCreateFirst: 'Import events to create your first dataset.',
    noDescription: 'No description',
    userPrefix: 'User',
    
    // Delete Dataset Modal
    deleteDatasetTitle: 'Delete Dataset',
    deleteConfirmQuestion: 'Are you sure you want to delete the dataset',
    deleteWillRemove: 'This will permanently delete:',
    eventsImported: 'events imported from this dataset',
    datasetRecordItself: 'The dataset record itself',
    tagsPreservedNote: 'Note: Tags created during import will be preserved.',
    cancel: 'Cancel',
    deleting: 'Deleting...',
    deleteDataset: 'Delete Dataset',
    
    // Create Dataset Modal
    createEmptyDatasetTitle: 'Create Empty Dataset',
    datasetName: 'Dataset Name*',
    datasetNamePlaceholder: 'e.g., My Historical Events',
    descriptionOptional: 'Description (optional)',
    descriptionPlaceholder: 'Brief description of this dataset...',
    creating: 'Creating...',
    createDataset: 'Create Dataset',
    datasetNameRequired: 'Dataset name is required',
    
    // Dataset Error Messages
    failedToLoadDatasets: 'Failed to load datasets',
    selectValidJson: 'Please select a valid JSON file',
    invalidFileFormat: 'Invalid file format: missing events array',
    invalidJsonFormat: 'Invalid JSON file format',
    failedToImport: 'Failed to import dataset',
    mustBeLoggedIn: 'You must be logged in to create datasets',
    failedToCreate: 'Failed to create dataset',
    failedToExport: 'Failed to export dataset',
    failedToDelete: 'Failed to delete dataset',
    
    // Dataset Action Tooltips
    exportDatasetTitle: 'Export dataset:',
    deleteDatasetTitle: 'Delete dataset:',
    unknown: 'Unknown'
  },
  ru: {
    // Filters
    filters: 'Фильтры',
    historicalEvents: 'Исторические события',
    historicalPeriod: 'Исторический период',
    specificPeriod: 'Сегмент:',
    selectPeriodGroup: 'Выбрать период',
    from: 'С:',
    to: 'По:',
    step: 'Шаг:',
    defaultPeriod: 'По умолчанию (1 н.э. - Сегодня)',
    customDateRange: 'Произвольный диапазон дат',
    selectSpecificPeriod: 'Выберите конкретный период...',
    noTemplatesAvailable: 'Шаблоны недоступны',
    loading: 'Загрузка',
    years: 'лет',
    apply: 'Применить',
    
    // Header Navigation
    map: 'Карта',
    admin: 'Управление',
    events: 'События',
    tags: 'Теги',
    datasets: 'Датасеты',
    users: 'Пользователи',
    
    // Authentication
    login: 'Войти',
    logout: 'Выйти',
    welcome: 'Добро пожаловать',
    loggingIn: 'Вход...',
    
    // Login Modal
    loginToHistoria: 'Вход в Historia ex machina',
    username: 'Имя пользователя:',
    password: 'Пароль:',
    enterUsername: 'Введите имя пользователя',
    enterPassword: 'Введите пароль',
    
    // Access Levels
    superBadge: 'СУПЕР',
    adminBadge: 'АДМИН',
    editorBadge: 'РЕДАКТОР',
    
    // Admin Events Page
    adminEventsTitle: 'Исторические события',
    adminEventsSubtitle: 'Управление и создание исторических событий',
    createNewEvent: 'Создать новое событие',
    searchEventName: 'Поиск по названию...',
    
    // Admin Tags Page
    adminTagsTitle: 'Управление тегами',
    adminTagsSubtitle: 'Управление тегами и категориями событий',
    createNewTag: 'Создать новый тег',
    
    // Admin Users Page
    adminUsersTitle: 'Управление пользователями',
    adminUsersSubtitle: 'Управление учетными записями и уровнями доступа',
    createNewUser: 'Создать нового пользователя',
    
    // Table Column Names
    columnName: 'Название',
    columnDescription: 'Описание',
    columnDate: 'Дата',
    columnLocation: 'Местоположение',
    columnType: 'Тип',
    columnTags: 'Теги',
    columnActions: 'Действия',
    columnColor: 'Цвет',
    columnCreated: 'Создан',
    columnUsageCount: 'Использование',
    columnUsername: 'Имя пользователя',
    columnEmail: 'Электронная почта',
    columnAccessLevel: 'Уровень доступа',
    columnStatus: 'Статус',
    columnLastActive: 'Последняя активность',
    
    // Month Names
    monthJan: 'янв',
    monthFeb: 'фев',
    monthMar: 'мар',
    monthApr: 'апр',
    monthMay: 'мая',
    monthJun: 'июн',
    monthJul: 'июл',
    monthAug: 'авг',
    monthSep: 'сен',
    monthOct: 'окт',
    monthNov: 'ноя',
    monthDec: 'дек',
    
    // Era Labels
    eraBC: 'до н.э.',
    eraAD: 'н.э.',
    
    // Admin Datasets Page
    adminDatasetsTitle: 'Датасеты событий',
    adminDatasetsSubtitle: 'Управление загруженными датасетами событий. Удаление датасета удалит все импортированные из него события.',
    loadingDatasets: 'Загрузка датасетов...',
    tryAgain: 'Попробовать снова',
    importNewDataset: 'Импортировать новый датасет',
    importDatasetDescription: 'Загрузите JSON-файл с историческими событиями для создания нового датасета.',
    chooseDatasetFile: 'Выбрать файл датасета',
    selectedFile: 'Выбран:',
    importing: 'Импорт...',
    importDataset: 'Импортировать датасет',
    clear: 'Очистить',
    or: 'или',
    createEmptyDataset: 'Создать пустой датасет',
    
    // Tag Filter
    filteredByTags: 'Фильтр по тегам',
    clearAllTags: 'Очистить все',
    clearAll: 'Очистить все',
    remove: 'Удалить',
    allTags: 'Все теги',
    searchTags: 'Поиск тегов...',
    noTagsFound: 'Теги не найдены',
    followEvents: 'Следовать',
    enableNarrativeFlow: 'Включить нарративный поток - соединить события хронологически',
    disableNarrativeFlow: 'Отключить нарративный поток',
    
    // Datasets Table
    columnFilename: 'Имя файла',
    columnEventCount: 'Количество событий',
    columnUploadedBy: 'Загрузил',
    columnUploadDate: 'Дата загрузки',
    noDatasetsFound: 'Датасеты не найдены',
    importToCreateFirst: 'Импортируйте события для создания первого датасета.',
    noDescription: 'Без описания',
    userPrefix: 'Пользователь',
    
    // Delete Dataset Modal
    deleteDatasetTitle: 'Удалить датасет',
    deleteConfirmQuestion: 'Вы уверены, что хотите удалить датасет',
    deleteWillRemove: 'Это безвозвратно удалит:',
    eventsImported: 'событий, импортированных из этого датасета',
    datasetRecordItself: 'Сам датасет',
    tagsPreservedNote: 'Примечание: Теги, созданные при импорте, будут сохранены.',
    cancel: 'Отмена',
    deleting: 'Удаление...',
    deleteDataset: 'Удалить датасет',
    
    // Create Dataset Modal
    createEmptyDatasetTitle: 'Создать пустой датасет',
    datasetName: 'Название датасета*',
    datasetNamePlaceholder: 'например, Мои исторические события',
    descriptionOptional: 'Описание (необязательно)',
    descriptionPlaceholder: 'Краткое описание датасета...',
    creating: 'Создание...',
    createDataset: 'Создать датасет',
    datasetNameRequired: 'Название датасета обязательно',
    
    // Dataset Error Messages
    failedToLoadDatasets: 'Не удалось загрузить датасеты',
    selectValidJson: 'Пожалуйста, выберите корректный JSON-файл',
    invalidFileFormat: 'Неверный формат файла: отсутствует массив events',
    invalidJsonFormat: 'Неверный формат JSON-файла',
    failedToImport: 'Не удалось импортировать датасет',
    mustBeLoggedIn: 'Вы должны войти в систему для создания датасетов',
    failedToCreate: 'Не удалось создать датасет',
    failedToExport: 'Не удалось экспортировать датасет',
    failedToDelete: 'Не удалось удалить датасет',
    
    // Dataset Action Tooltips
    exportDatasetTitle: 'Экспортировать датасет:',
    deleteDatasetTitle: 'Удалить датасет:',
    unknown: 'Неизвестно'
  }
}

// Supported locales
const SUPPORTED_LOCALES = {
  en: {
    code: 'en',
    name: 'English',
    flag: '🇺🇸'
  },
  ru: {
    code: 'ru', 
    name: 'Русский',
    flag: '🇷🇺'
  }
}

// Global locale state with validation
const locale = ref(localStorage.getItem('historia-locale') || 'en')

// Validate and fix invalid locale in localStorage
if (!SUPPORTED_LOCALES[locale.value]) {
  locale.value = 'en'
  localStorage.setItem('historia-locale', 'en')
}

export function useLocale() {
  const currentLocale = computed(() => SUPPORTED_LOCALES[locale.value] || SUPPORTED_LOCALES.en)
  const supportedLocales = computed(() => Object.values(SUPPORTED_LOCALES))
  
  // Get translated text for current locale
  const t = (key) => {
    const translations = UI_TRANSLATIONS[locale.value] || UI_TRANSLATIONS.en
    return translations[key] || key
  }
  
  const setLocale = (newLocale) => {
    if (SUPPORTED_LOCALES[newLocale]) {
      locale.value = newLocale
      localStorage.setItem('historia-locale', newLocale)
    }
  }
  
  const getLocaleParam = () => {
    return `locale=${locale.value}`
  }
  
  const addLocaleToUrl = (url) => {
    const separator = url.includes('?') ? '&' : '?'
    return `${url}${separator}${getLocaleParam()}`
  }
  
  // Format date with localized month names and era labels
  // Input: ISO date string like "1453-05-29T00:00:00Z" or "-3500-01-01T00:00:00Z"
  // Output: "29 May 1453 AD" or "29 мая 1453 н.э." or "1 Jan 3500 BC" or "1 янв 3500 до н.э."
  const formatLocalizedDate = (isoDateString, era) => {
    const monthKeys = [
      'monthJan', 'monthFeb', 'monthMar', 'monthApr', 'monthMay', 'monthJun',
      'monthJul', 'monthAug', 'monthSep', 'monthOct', 'monthNov', 'monthDec'
    ]
    
    let year, month, day
    
    if (isoDateString.startsWith('-')) {
      // Negative year format: "-3500-01-01T00:00:00Z"
      const parts = isoDateString.substring(1).split('T')[0].split('-')
      year = parseInt(parts[0], 10)
      month = parseInt(parts[1], 10) - 1
      day = parseInt(parts[2], 10)
    } else {
      // Positive year format: "1453-05-29T00:00:00Z"
      // Use manual parsing to avoid timezone issues
      const parts = isoDateString.split('T')[0].split('-')
      year = parseInt(parts[0], 10)
      month = parseInt(parts[1], 10) - 1
      day = parseInt(parts[2], 10)
    }
    
    const monthName = t(monthKeys[month])
    const eraLabel = era === 'BC' ? t('eraBC') : t('eraAD')
    
    return `${day} ${monthName} ${year} ${eraLabel}`
  }

  // Format event display date with DD.MM.YYYY format and localized era labels
  // Input: ISO date string and era ('BC' or 'AD')
  // Output: "29.05.1453 AD" or "29.05.1453 н.э." or "01.01.3500 BC" or "01.01.3500 до н.э."
  const formatEventDisplayDate = (isoDateString, era) => {
    if (!isoDateString) return ''
    
    let year, month, day
    
    if (isoDateString.startsWith('-')) {
      // Negative year format: "-3500-01-01T00:00:00Z"
      const parts = isoDateString.substring(1).split('T')[0].split('-')
      year = parseInt(parts[0], 10)
      month = parseInt(parts[1], 10)
      day = parseInt(parts[2], 10)
    } else {
      // Positive year format: "1453-05-29T00:00:00Z"
      const parts = isoDateString.split('T')[0].split('-')
      year = parseInt(parts[0], 10)
      month = parseInt(parts[1], 10)
      day = parseInt(parts[2], 10)
    }
    
    const paddedDay = String(day).padStart(2, '0')
    const paddedMonth = String(month).padStart(2, '0')
    const eraLabel = era === 'BC' ? t('eraBC') : t('eraAD')
    
    return `${paddedDay}.${paddedMonth}.${year} ${eraLabel}`
  }
  
  return {
    locale: computed(() => locale.value),
    currentLocale,
    supportedLocales,
    t,
    setLocale,
    getLocaleParam,
    addLocaleToUrl,
    formatLocalizedDate,
    formatEventDisplayDate
  }
}