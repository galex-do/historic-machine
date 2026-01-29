import { ref, watch } from 'vue'

const YEARS_RANGE = 3
const DISTANCE_KM = 100
const MAX_RESULTS = 3

const getDistanceKm = (lat1, lon1, lat2, lon2) => {
  const R = 6371
  const dLat = (lat2 - lat1) * Math.PI / 180
  const dLon = (lon2 - lon1) * Math.PI / 180
  const a = 
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) *
    Math.sin(dLon / 2) * Math.sin(dLon / 2)
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))
  return R * c
}

const getChronologicalValue = (dateString, era) => {
  let year, month, day
  
  if (dateString.startsWith('-')) {
    const parts = dateString.substring(1).split('T')[0].split('-')
    year = parseInt(parts[0], 10)
    month = parseInt(parts[1], 10) - 1
    day = parseInt(parts[2], 10)
  } else {
    const parts = dateString.split('T')[0].split('-')
    year = parseInt(parts[0], 10)
    month = parseInt(parts[1], 10) - 1
    day = parseInt(parts[2], 10)
  }
  
  if (era === 'BC') {
    return -(year - (month / 12) - (day / 365))
  } else {
    return year + (month / 12) + (day / 365)
  }
}

const getTimeDifferenceYears = (event1, event2) => {
  const val1 = getChronologicalValue(event1.event_date, event1.era)
  const val2 = getChronologicalValue(event2.event_date, event2.era)
  return Math.abs(val1 - val2)
}

const shuffleArray = (array) => {
  const result = [...array]
  for (let i = result.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[result[i], result[j]] = [result[j], result[i]]
  }
  return result
}

const hasSharedTag = (event1, event2) => {
  if (!event1.tags?.length || !event2.tags?.length) return false
  const tagIds1 = new Set(event1.tags.map(t => t.id))
  return event2.tags.some(t => tagIds1.has(t.id))
}

const sortByDate = (events) => {
  return [...events].sort((a, b) => {
    const aVal = getChronologicalValue(a.event_date, a.era)
    const bVal = getChronologicalValue(b.event_date, b.era)
    return aVal - bVal
  })
}

export function useRelatedEvents(currentEvent, allEvents) {
  const aroundSameTime = ref([])
  const sameTimeRegion = ref([])
  const nearByKind = ref([])
  const isComputed = ref(false)

  const computeRelatedEvents = () => {
    if (!currentEvent.value?.id || !allEvents.value?.length) {
      aroundSameTime.value = []
      sameTimeRegion.value = []
      nearByKind.value = []
      isComputed.value = false
      return
    }

    const current = currentEvent.value
    const currentLat = parseFloat(current.latitude)
    const currentLon = parseFloat(current.longitude)
    
    const candidates = allEvents.value.filter(e => {
      if (e.id === current.id) return false
      const yearDiff = getTimeDifferenceYears(current, e)
      return yearDiff <= YEARS_RANGE
    })

    const usedIds = new Set()

    const regionCandidates = candidates.filter(e => {
      const lat = parseFloat(e.latitude)
      const lon = parseFloat(e.longitude)
      if (isNaN(lat) || isNaN(lon) || isNaN(currentLat) || isNaN(currentLon)) return false
      const distance = getDistanceKm(currentLat, currentLon, lat, lon)
      return distance <= DISTANCE_KM
    })
    
    const shuffledRegion = shuffleArray(regionCandidates).slice(0, MAX_RESULTS)
    sameTimeRegion.value = sortByDate(shuffledRegion)
    shuffledRegion.forEach(e => usedIds.add(e.id))

    const timeCandidates = candidates.filter(e => !usedIds.has(e.id))
    const shuffledTime = shuffleArray(timeCandidates).slice(0, MAX_RESULTS)
    aroundSameTime.value = sortByDate(shuffledTime)
    shuffledTime.forEach(e => usedIds.add(e.id))

    const tagCandidates = candidates
      .filter(e => !usedIds.has(e.id) && hasSharedTag(current, e))
      .map(e => ({
        event: e,
        timeDiff: getTimeDifferenceYears(current, e)
      }))
      .sort((a, b) => a.timeDiff - b.timeDiff)
      .slice(0, MAX_RESULTS)
      .map(item => item.event)
    
    nearByKind.value = sortByDate(tagCandidates)

    isComputed.value = true
  }

  watch([currentEvent, allEvents], () => {
    computeRelatedEvents()
  }, { immediate: true, deep: false })

  return {
    aroundSameTime,
    sameTimeRegion,
    nearByKind,
    isComputed
  }
}
