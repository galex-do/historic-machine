export function getContrastColor(hexColor) {
  if (!hexColor) return '#000000'
  const color = hexColor.replace('#', '')
  const r = parseInt(color.substr(0, 2), 16)
  const g = parseInt(color.substr(2, 2), 16)
  const b = parseInt(color.substr(4, 2), 16)
  const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255
  return luminance > 0.5 ? '#000000' : '#ffffff'
}

export function getTagBorderStyle(borderColor) {
  if (!borderColor) return {}
  return { boxShadow: `inset 0 0 0 2.5px ${borderColor}` }
}

export function getTagStyle(tag, options = {}) {
  const bgColor = tag.color || '#6366f1'
  const style = {
    backgroundColor: bgColor,
    color: getContrastColor(bgColor)
  }
  if (tag.border_color) {
    const shadows = [`inset 0 0 0 2.5px ${tag.border_color}`]
    if (options.outerShadow) {
      shadows.push(options.outerShadow)
    }
    style.boxShadow = shadows.join(', ')
  }
  return style
}
