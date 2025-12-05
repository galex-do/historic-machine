import { ref } from 'vue'
import { useLocale } from './useLocale.js'

export function useGeolocation() {
  const { t } = useLocale()
  
  const coords = ref(null)
  const error = ref(null)
  const loading = ref(false)

  const is_supported = () => {
    return 'geolocation' in navigator
  }

  const get_error_message = (err) => {
    switch (err.code) {
      case err.PERMISSION_DENIED:
        return t('geolocationPermissionDenied')
      case err.POSITION_UNAVAILABLE:
        return t('geolocationUnavailable')
      case err.TIMEOUT:
        return t('geolocationTimeout')
      default:
        return t('geolocationError')
    }
  }

  const get_current_position = () => {
    return new Promise((resolve, reject) => {
      if (!is_supported()) {
        reject(new Error(t('geolocationNotSupported')))
        return
      }

      loading.value = true
      error.value = null
      coords.value = null

      navigator.geolocation.getCurrentPosition(
        (position) => {
          coords.value = {
            latitude: position.coords.latitude,
            longitude: position.coords.longitude,
            accuracy: position.coords.accuracy
          }
          loading.value = false
          resolve(coords.value)
        },
        (err) => {
          const message = get_error_message(err)
          error.value = message
          loading.value = false
          reject(new Error(message))
        },
        {
          enableHighAccuracy: true,
          timeout: 10000,
          maximumAge: 60000
        }
      )
    })
  }

  return {
    coords,
    error,
    loading,
    is_supported,
    get_current_position
  }
}
