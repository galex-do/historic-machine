import { ref, computed } from 'vue'

// UI Text translations
const UI_TRANSLATIONS = {
  en: {
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
    years: 'years'
  },
  ru: {
    filters: 'Фильтры',
    historicalEvents: 'Исторические события',
    historicalPeriod: 'Исторический период:',
    specificPeriod: 'Конкретный период:',
    from: 'С:',
    to: 'По:',
    step: 'Шаг:',
    defaultPeriod: 'По умолчанию (1 н.э. - Сегодня)',
    customDateRange: 'Произвольный диапазон дат',
    selectSpecificPeriod: 'Выберите конкретный период...',
    years: 'лет'
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
  
  return {
    locale: computed(() => locale.value),
    currentLocale,
    supportedLocales,
    t,
    setLocale,
    getLocaleParam,
    addLocaleToUrl
  }
}