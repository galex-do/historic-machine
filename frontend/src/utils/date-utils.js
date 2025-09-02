/**
 * Date formatting and parsing utilities
 */

/**
 * Format date to DD.MM.YYYY format
 */
export function formatDateToDDMMYYYY(date) {
  let day, month, year

  if (typeof date === 'string') {
    // Parse ISO string manually for very old dates
    const [yearStr, monthStr, dayStr] = date.split('-')
    day = parseInt(dayStr, 10)
    month = parseInt(monthStr, 10) 
    year = parseInt(yearStr, 10)
  } else {
    // Regular Date object
    day = date.getDate()
    month = date.getMonth() + 1
    year = date.getFullYear()
  }
  
  const paddedDay = String(day).padStart(2, '0')
  const paddedMonth = String(month).padStart(2, '0')
  
  return `${paddedDay}.${paddedMonth}.${year}`
}

/**
 * Parse DD.MM.YYYY format to ISO date string
 */
export function parseDDMMYYYYToISO(dateStr) {
  // Allow years from 1-4 digits (year 1 to 9999)
  if (!dateStr || !dateStr.match(/^\d{1,2}\.\d{1,2}\.\d{1,4}$/)) {
    return null
  }
  
  const [day, month, year] = dateStr.split('.')
  const yearNum = parseInt(year, 10)
  const monthNum = parseInt(month, 10)
  const dayNum = parseInt(day, 10)
  
  // Basic validation
  if (yearNum < 1 || yearNum > 9999) return null
  if (monthNum < 1 || monthNum > 12) return null
  if (dayNum < 1 || dayNum > 31) return null
  
  // For very old years, construct ISO string manually to avoid Date object issues
  const paddedYear = String(yearNum).padStart(4, '0')
  const paddedMonth = String(monthNum).padStart(2, '0')
  const paddedDay = String(dayNum).padStart(2, '0')
  
  return `${paddedYear}-${paddedMonth}-${paddedDay}`
}

/**
 * Format date string for display
 */
export function formatDate(dateString) {
  if (!dateString) return ''
  
  // Ensure dateString is a string before processing
  const dateStr = String(dateString)
  
  // For very old dates, parse manually to avoid Date object issues
  if (dateStr.startsWith('00') || dateStr.startsWith('01')) {
    return formatDateToDDMMYYYY(dateStr)
  }
  return formatDateToDDMMYYYY(new Date(dateStr))
}

/**
 * Get today's date in ISO format
 */
export function getTodayISO() {
  return new Date().toISOString().split('T')[0]
}

/**
 * Parse historical date string with BC/AD support
 * Formats: "500 BC", "1066 AD", "500BC", "1066AD", "500", "1066"
 * Returns: { year: number, era: 'BC'|'AD', isoDate: string }
 */
export function parseHistoricalDate(dateStr) {
  if (!dateStr || typeof dateStr !== 'string') {
    return null
  }
  
  // Clean up the input
  const cleanStr = dateStr.trim().toUpperCase()
  
  // Match patterns like "500 BC", "1066 AD", "500BC", "1066AD"
  const bcMatch = cleanStr.match(/^(\d{1,4})\s*BC$/)
  const adMatch = cleanStr.match(/^(\d{1,4})\s*AD$/)
  const yearOnlyMatch = cleanStr.match(/^(\d{1,4})$/)
  
  let year, era
  
  if (bcMatch) {
    year = parseInt(bcMatch[1], 10)
    era = 'BC'
  } else if (adMatch) {
    year = parseInt(adMatch[1], 10)
    era = 'AD'
  } else if (yearOnlyMatch) {
    year = parseInt(yearOnlyMatch[1], 10)
    // Default assumption: years below 1000 are likely BC for historical mapping
    era = year < 500 ? 'BC' : 'AD'
  } else {
    return null
  }
  
  // Validate year range
  if (year < 1 || year > 9999) {
    return null
  }
  
  // Convert to ISO date (astronomical year numbering for proper BC/AD comparison)
  const isoYear = era === 'BC' ? (year === 1 ? '0000' : String(year - 1).padStart(4, '0')) : String(year).padStart(4, '0')
  const isoDate = era === 'BC' ? `-${isoYear}-01-01` : `${isoYear}-01-01`
  
  return {
    year,
    era,
    isoDate,
    displayString: `${year} ${era}`
  }
}

/**
 * Format historical date for display
 * Input: ISO date string or parsed historical date object
 * Output: "500 BC", "1066 AD", etc.
 */
export function formatHistoricalDate(dateInput) {
  if (!dateInput) return ''
  
  // If it's already a parsed object, use it
  if (typeof dateInput === 'object' && dateInput.year && dateInput.era) {
    return `${dateInput.year} ${dateInput.era}`
  }
  
  // If it's an ISO string, parse it
  if (typeof dateInput === 'string') {
    const year = parseInt(dateInput.split('-')[0], 10)
    const era = dateInput.startsWith('-') || year < 500 ? 'BC' : 'AD'
    const displayYear = dateInput.startsWith('-') ? Math.abs(year) + 1 : year
    return `${displayYear} ${era}`
  }
  
  return ''
}

/**
 * Add years to a historical date
 * Used for step functionality
 */
export function addYearsToHistoricalDate(dateStr, years) {
  const parsed = parseHistoricalDate(dateStr)
  if (!parsed) return dateStr
  
  let newYear = parsed.year
  let newEra = parsed.era
  
  if (parsed.era === 'BC') {
    newYear = parsed.year - years
    if (newYear <= 0) {
      newYear = Math.abs(newYear) + 1
      newEra = 'AD'
    }
  } else {
    newYear = parsed.year + years
    if (newYear <= 0) {
      newYear = Math.abs(newYear) + 1
      newEra = 'BC'
    }
  }
  
  // Ensure year is within valid range
  if (newYear < 1) newYear = 1
  if (newYear > 9999) newYear = 9999
  
  return `${newYear} ${newEra}`
}

/**
 * Convert historical date to ISO format for API calls
 */
export function historicalDateToISO(dateStr) {
  const parsed = parseHistoricalDate(dateStr)
  return parsed ? parsed.isoDate : null
}

/**
 * Convert date to astronomical year for BC/AD calculations
 */
export function getAstronomicalYear(dateString, era) {
  if (!dateString) return null
  
  const year = parseInt(dateString.split('-')[0], 10)
  
  if (era === 'BC') {
    return (year * -1) + 1  // Convert BC year to astronomical year
  } else {
    return year  // AD years are already correct
  }
}

/**
 * Get era from date string (heuristic)
 */
export function getEraFromDate(dateString) {
  if (!dateString) return 'AD'
  
  const year = parseInt(dateString.split('-')[0], 10)
  // Default assumption based on year
  return year < 1000 ? 'BC' : 'AD'  // Heuristic for old dates
}