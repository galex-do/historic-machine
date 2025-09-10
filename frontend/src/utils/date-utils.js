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
 * Formats: "500 BC", "1066 AD", "15.03.44 BC", "25.12.1066 AD", "500", "1066"
 * Returns: { year: number, era: 'BC'|'AD', day?: number, month?: number, isoDate: string }
 */
export function parseHistoricalDate(dateStr) {
  if (!dateStr || typeof dateStr !== 'string') {
    return null
  }
  
  // Clean up the input
  const cleanStr = dateStr.trim().toUpperCase()
  
  // Match patterns with DD.MM.YYYY format and BC/AD
  const ddmmyyyyBcMatch = cleanStr.match(/^(\d{1,2})\.(\d{1,2})\.(\d{1,4})\s*BC$/)
  const ddmmyyyyAdMatch = cleanStr.match(/^(\d{1,2})\.(\d{1,2})\.(\d{1,4})\s*AD$/)
  
  // Match patterns like "500 BC", "1066 AD", "500BC", "1066AD"
  const bcMatch = cleanStr.match(/^(\d{1,4})\s*BC$/)
  const adMatch = cleanStr.match(/^(\d{1,4})\s*AD$/)
  const yearOnlyMatch = cleanStr.match(/^(\d{1,4})$/)
  
  let year, era, day, month
  
  if (ddmmyyyyBcMatch) {
    day = parseInt(ddmmyyyyBcMatch[1], 10)
    month = parseInt(ddmmyyyyBcMatch[2], 10)
    year = parseInt(ddmmyyyyBcMatch[3], 10)
    era = 'BC'
  } else if (ddmmyyyyAdMatch) {
    day = parseInt(ddmmyyyyAdMatch[1], 10)
    month = parseInt(ddmmyyyyAdMatch[2], 10)
    year = parseInt(ddmmyyyyAdMatch[3], 10)
    era = 'AD'
  } else if (bcMatch) {
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
  
  // Validate day and month if provided
  if (day !== undefined && month !== undefined) {
    if (month < 1 || month > 12 || day < 1 || day > 31) {
      return null
    }
  }
  
  // Convert to ISO date (use actual year, no astronomical numbering)
  const isoYear = String(year).padStart(4, '0')
  const isoMonth = month ? String(month).padStart(2, '0') : '01'
  const isoDay = day ? String(day).padStart(2, '0') : '01'
  const isoDate = era === 'BC' ? `-${isoYear}-${isoMonth}-${isoDay}` : `${isoYear}-${isoMonth}-${isoDay}`
  
  const result = {
    year,
    era,
    isoDate,
    displayString: day && month ? `${String(day).padStart(2, '0')}.${String(month).padStart(2, '0')}.${year} ${era}` : `${year} ${era}`
  }
  
  if (day !== undefined && month !== undefined) {
    result.day = day
    result.month = month
  }
  
  return result
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
    let year, era, displayYear
    
    if (dateInput.startsWith('-')) {
      // BC date format: "-491-09-12T00:00:00Z"
      const yearMatch = dateInput.match(/^-(\d+)-/)
      if (yearMatch) {
        year = parseInt(yearMatch[1], 10)
        displayYear = year
        era = 'BC'
      } else {
        return ''
      }
    } else {
      // AD date format: "1453-05-29T00:00:00Z"
      year = parseInt(dateInput.split('-')[0], 10)
      era = year < 500 ? 'BC' : 'AD' // Legacy logic for old format
      displayYear = year
    }
    
    return `${displayYear} ${era}`
  }
  
  return ''
}

/**
 * Add years to a historical date
 * Used for step functionality - preserves day and month if present
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
  
  // Preserve day and month if they were in the original date
  if (parsed.day !== undefined && parsed.month !== undefined) {
    const paddedDay = String(parsed.day).padStart(2, '0')
    const paddedMonth = String(parsed.month).padStart(2, '0')
    return `${paddedDay}.${paddedMonth}.${newYear} ${newEra}`
  } else {
    return `${newYear} ${newEra}`
  }
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
  
  // Return the actual year as stored in database (no astronomical conversion)
  return year
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