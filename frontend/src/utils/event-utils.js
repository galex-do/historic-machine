/**
 * Event-related utility functions
 */

/**
 * Get emoji for event lens type
 */
export function getEventEmoji(lensType) {
  const emojiMap = {
    'historic': '📜',
    'political': '🏛️', 
    'cultural': '🎭',
    'military': '⚔️',
    'scientific': '🔬'
  }
  return emojiMap[lensType] || '📍'
}

/**
 * Get label for lens type
 */
export function getLensLabel(lensType) {
  const labelMap = {
    'historic': 'Historic',
    'political': 'Political',
    'cultural': 'Cultural', 
    'military': 'Military',
    'scientific': 'Scientific'
  }
  return labelMap[lensType] || lensType
}

/**
 * Get available lens types with labels
 */
export function getAvailableLensTypes() {
  return [
    { value: 'historic', label: '📜 Historic' },
    { value: 'political', label: '🏛️ Political' },
    { value: 'cultural', label: '🎭 Cultural' },
    { value: 'military', label: '⚔️ Military' },
    { value: 'scientific', label: '🔬 Scientific' }
  ]
}

/**
 * Validate event coordinates
 */
export function validateCoordinates(lat, lng) {
  if (lat < -90 || lat > 90) {
    return { valid: false, error: 'Latitude must be between -90 and 90 degrees' }
  }
  if (lng < -180 || lng > 180) {
    return { valid: false, error: 'Longitude must be between -180 and 180 degrees' }
  }
  return { valid: true }
}