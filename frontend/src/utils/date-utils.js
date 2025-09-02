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