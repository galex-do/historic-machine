/**
 * Event-related utility functions
 */

/**
 * Get first custom emoji set on any of the event's tags (ordered by tag weight).
 * Returns null when no tag carries an emoji.
 */
export function getTagEmoji(tags) {
  if (!tags || !Array.isArray(tags)) return null
  for (const tag of tags) {
    if (tag.emoji) return tag.emoji
  }
  return null
}

/**
 * Resolve the display emoji for an event: tag emoji takes priority over lens-type emoji.
 */
export function resolveEventEmoji(event) {
  return getTagEmoji(event?.tags) || getEventEmoji(event?.lens_type)
}

/**
 * Get emoji for event lens type
 */
export function getEventEmoji(lensType) {
  const emojiMap = {
    'historic': '📜',
    'political': '🏛️', 
    'cultural': '🎭',
    'military': '⚔️',
    'scientific': '🔬',
    'religious': '⛪'
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
    'scientific': 'Scientific',
    'religious': 'Religious'
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
    { value: 'scientific', label: '🔬 Scientific' },
    { value: 'religious', label: '⛪ Religious' }
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