import { ref, computed } from 'vue'

// UI Text translations
const UI_TRANSLATIONS = {
  en: {
    // Filters
    filters: 'Filters',
    historicalEvents: 'Historical Events',
    historicalPeriod: 'Historical Period:',
    specificPeriod: 'Specific Period:',
    from: 'From:',
    to: 'To:',
    step: 'Step:',
    defaultPeriod: 'Default (1 AD - Today)',
    customDateRange: 'Custom Date Range',
    selectSpecificPeriod: 'Select specific period...',
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
    eraAD: 'AD'
  },
  ru: {
    // Filters
    filters: 'Фильтры',
    historicalEvents: 'Исторические события',
    historicalPeriod: 'Исторический период:',
    specificPeriod: 'Сегмент:',
    from: 'С:',
    to: 'По:',
    step: 'Шаг:',
    defaultPeriod: 'По умолчанию (1 н.э. - Сегодня)',
    customDateRange: 'Произвольный диапазон дат',
    selectSpecificPeriod: 'Выберите конкретный период...',
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
    eraAD: 'н.э.'
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
  
  return {
    locale: computed(() => locale.value),
    currentLocale,
    supportedLocales,
    t,
    setLocale,
    getLocaleParam,
    addLocaleToUrl,
    formatLocalizedDate
  }
}