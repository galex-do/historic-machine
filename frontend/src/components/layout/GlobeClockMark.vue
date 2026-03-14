<template>
  <svg :width="size" :height="size" :viewBox="`0 0 ${size} ${size}`" fill="none" :aria-label="label">
    <defs>
      <linearGradient :id="`grad-${uid}`" x1="0" y1="0" :x2="size" :y2="size" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#4f46e5" />
        <stop offset="100%" stop-color="#3730a3" />
      </linearGradient>
      <clipPath :id="`g-${uid}`"><path :d="globeClip" /></clipPath>
      <clipPath :id="`c-${uid}`"><path :d="clockClip" /></clipPath>
    </defs>

    <!-- Circle base -->
    <circle :cx="C" :cy="C" :r="R" :fill="circleFill" />

    <!-- Globe — continent silhouettes -->
    <g :clip-path="`url(#g-${uid})`" :fill="fg" fill-opacity="0.88" stroke="none">
      <path v-for="(d, i) in continents" :key="i" :d="d" />
    </g>

    <!-- Map pin on the meridian -->
    <g :clip-path="`url(#g-${uid})`">
      <circle :cx="pinX" :cy="pinY" :r="s * 5.5" :fill="fg" />
      <circle :cx="pinX" :cy="pinY" :r="s * 2.4" :fill="inverted ? '#e0e7ff' : '#3730a3'" />
      <line :x1="pinX" :y1="pinY + s * 5.5" :x2="pinX" :y2="pinY + s * 10"
        :stroke="fg" :stroke-width="s * 3" stroke-linecap="round" />
    </g>

    <!-- Clock — ticks + hands -->
    <g :clip-path="`url(#c-${uid})`" stroke-linecap="round">
      <line v-for="(t, i) in majorTicks" :key="i"
        :x1="t.x1" :y1="t.y1" :x2="t.x2" :y2="t.y2"
        :stroke="fg" :stroke-width="s * 3.2" />
      <line :x1="C" :y1="C"
        :x2="C + R * 0.44 * Math.cos(hA)" :y2="C + R * 0.44 * Math.sin(hA)"
        :stroke="fg" :stroke-width="s * 4.5" />
      <line :x1="C" :y1="C"
        :x2="C + R * 0.66 * Math.cos(mA)" :y2="C + R * 0.66 * Math.sin(mA)"
        :stroke="fg" :stroke-width="s * 3" />
    </g>

    <!-- Divider lines -->
    <line :x1="C" :y1="C - R + s" :x2="C" :y2="C"
      :stroke="fg" :stroke-width="s * 2" stroke-opacity="0.4" />
    <line :x1="C" :y1="C" :x2="edgeX" :y2="edgeY"
      :stroke="fg" :stroke-width="s * 2" stroke-opacity="0.4" />

    <!-- Center dot -->
    <circle :cx="C" :cy="C" :r="s * 4" :fill="fg" />
  </svg>
</template>

<script>
import { computed } from 'vue'

let uid_counter = 0

export default {
  name: 'GlobeClockMark',
  props: {
    size:     { type: Number, default: 32 },
    inverted: { type: Boolean, default: false },
    label:    { type: String, default: 'timediverr logo' },
  },
  setup(props) {
    const uid = `gcm-${uid_counter++}`

    const C = computed(() => props.size / 2)
    const R = computed(() => props.size / 2 - props.size * 0.05)
    const s = computed(() => props.size / 100)

    const fg          = computed(() => props.inverted ? '#3730a3' : '#fff')
    const circleFill  = computed(() => props.inverted ? '#e0e7ff' : `url(#grad-${uid})`)

    // Minute hand at 4 o'clock
    const mA = (4 * 30 - 90) * Math.PI / 180
    // Hour hand at 2 o'clock
    const hA = (2 * 30 - 90) * Math.PI / 180

    const edgeX = computed(() => C.value + R.value * Math.cos(mA))
    const edgeY = computed(() => C.value + R.value * Math.sin(mA))

    const globeClip = computed(() =>
      `M${C.value} ${C.value - R.value} A${R.value} ${R.value} 0 1 0 ${edgeX.value} ${edgeY.value} L${C.value} ${C.value} Z`
    )
    const clockClip = computed(() =>
      `M${C.value} ${C.value - R.value} A${R.value} ${R.value} 0 0 1 ${edgeX.value} ${edgeY.value} L${C.value} ${C.value} Z`
    )

    const majorTicks = computed(() =>
      [0, 1, 2, 3, 4].map(i => {
        const a = (i * 30 - 90) * Math.PI / 180
        return {
          x1: C.value + (R.value - s.value * 18) * Math.cos(a),
          y1: C.value + (R.value - s.value * 18) * Math.sin(a),
          x2: C.value + (R.value - s.value * 6)  * Math.cos(a),
          y2: C.value + (R.value - s.value * 6)  * Math.sin(a),
        }
      })
    )

    const pinX = computed(() => C.value - R.value * 0.44)
    const pinY = computed(() => C.value - R.value * 0.48)

    // Continent polygons — orthographic projection, center 30°W 0°N
    const continents = computed(() => {
      const c = C.value, r = R.value
      const q = (nx, ny) => `${c + nx * r} ${c + ny * r}`
      const poly = pts => `M${pts.map(([x, y]) => q(x, y)).join('L')}Z`
      return [
        // Greenland
        poly([[-0.063,-0.970],[-0.021,-0.970],[0.027,-0.951],[0.068,-0.921],[0.034,-0.875],[-0.121,-0.866],[-0.152,-0.914],[-0.106,-0.951]]),
        // North America
        poly([[-0.358,-0.934],[-0.265,-0.956],[-0.324,-0.906],[-0.266,-0.731],[-0.354,-0.707],[-0.433,-0.695],[-0.589,-0.574],[-0.704,-0.423],[-0.718,-0.500],[-0.757,-0.485],[-0.827,-0.438],[-0.824,-0.309],[-0.871,-0.276],[-0.907,-0.391],[-0.799,-0.602],[-0.655,-0.755],[-0.528,-0.839]]),
        // South America
        poly([[-0.743,-0.017],[-0.724,-0.139],[-0.536,-0.174],[-0.374,-0.070],[-0.086,0.139],[-0.136,0.225],[-0.207,0.391],[-0.310,0.500],[-0.389,0.559],[-0.344,0.829],[-0.435,0.788],[-0.470,0.682],[-0.766,0.035]]),
        // Africa (shifted -0.25 in x so southern body is in globe sector)
        poly([[-0.032,-0.242],[0.328,-0.530],[0.509,-0.375],[0.578,0.407],[0.468,0.500],[0.378,0.574],[0.366,0.559],[0.394,0.375],[0.396,0.259],[0.418,0.035],[0.293,-0.087],[0.022,-0.174]]),
      ]
    })

    return {
      uid,
      C, R, s,
      fg, circleFill,
      mA, hA,
      edgeX, edgeY,
      globeClip, clockClip,
      majorTicks,
      pinX, pinY,
      continents,
      Math,
    }
  }
}
</script>
