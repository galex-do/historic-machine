function Mark({ size = 140, inverted = false }) {
  const C = size / 2
  const R = size / 2 - size * 0.05
  const s = size / 100

  const fg = inverted ? '#3730a3' : '#fff'
  const circleFill = inverted ? '#e0e7ff' : 'url(#grad)'

  // Minute hand at 4 o'clock — this defines the globe/clock boundary
  const mA = (4 * 30 - 90) * Math.PI / 180  // 30° from horizontal
  const edgeX = C + R * Math.cos(mA)         // circle edge at 4 o'clock
  const edgeY = C + R * Math.sin(mA)

  // Hour hand at 2 o'clock
  const hA = (2 * 30 - 90) * Math.PI / 180

  // Globe clip: large 240° sector — counterclockwise from top to 4 o'clock edge
  const globeClip = `M${C} ${C - R} A${R} ${R} 0 1 0 ${edgeX} ${edgeY} L${C} ${C} Z`
  // Clock clip: small 120° sector — clockwise from top to 4 o'clock edge
  const clockClip = `M${C} ${C - R} A${R} ${R} 0 0 1 ${edgeX} ${edgeY} L${C} ${C} Z`

  // Clock ticks: 12/3 visible in the small sector (6/9 fall in globe area, clipped away)
  const majorTicks = [0, 1, 2, 3, 4].map((i) => {
    const a = (i * 30 - 90) * Math.PI / 180
    return {
      x1: C + (R - s * 18) * Math.cos(a), y1: C + (R - s * 18) * Math.sin(a),
      x2: C + (R - s * 6)  * Math.cos(a), y2: C + (R - s * 6)  * Math.sin(a),
    }
  })

  // Pin sits exactly on the meridian arc at y = C - R*0.48
  // Meridian ellipse: cx=C, cy=C, rx=R/2, ry=R → at that y, x = C - R*0.44
  const pinX = C - R * 0.44
  const pinY = C - R * 0.48

  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`} fill="none">
      <defs>
        <linearGradient id="grad" x1="0" y1="0" x2={size} y2={size} gradientUnits="userSpaceOnUse">
          <stop offset="0%" stopColor="#4f46e5" />
          <stop offset="100%" stopColor="#3730a3" />
        </linearGradient>
        <clipPath id={`g${size}`}><path d={globeClip} /></clipPath>
        <clipPath id={`c${size}`}><path d={clockClip} /></clipPath>
      </defs>

      {/* Circle base */}
      <circle cx={C} cy={C} r={R} fill={circleFill} />

      {/* GLOBE — continent silhouettes */}
      <g clipPath={`url(#g${size})`} fill={fg} fillOpacity="0.88" stroke="none">
        {/* Helper: normalized (nx,ny) → SVG coords, nx/ny in [-1,1] */}
        {/* Greenland */}
        <path d={[
          `M${C-R*0.55},${C-R*0.80}`,
          `L${C-R*0.32},${C-R*0.84}`,
          `L${C-R*0.24},${C-R*0.68}`,
          `L${C-R*0.34},${C-R*0.58}`,
          `L${C-R*0.60},${C-R*0.58}`,
          `L${C-R*0.68},${C-R*0.70}`,
          `Z`
        ].join(' ')} />
        {/* North America */}
        <path d={[
          `M${C-R*0.72},${C-R*0.58}`,
          `L${C-R*0.24},${C-R*0.64}`,
          `L${C-R*0.10},${C-R*0.44}`,
          `L${C-R*0.16},${C-R*0.20}`,
          `L${C-R*0.46},${C-R*0.05}`,
          `L${C-R*0.68},${C-R*0.12}`,
          `L${C-R*0.80},${C-R*0.36}`,
          `Z`
        ].join(' ')} />
        {/* South America */}
        <path d={[
          `M${C-R*0.22},${C-R*0.02}`,
          `L${C-R*0.02},${C+R*0.08}`,
          `L${C-R*0.02},${C+R*0.35}`,
          `L${C-R*0.18},${C+R*0.68}`,
          `L${C-R*0.42},${C+R*0.64}`,
          `L${C-R*0.52},${C+R*0.44}`,
          `L${C-R*0.40},${C+R*0.06}`,
          `Z`
        ].join(' ')} />
        {/* Europe */}
        <path d={[
          `M${C-R*0.12},${C-R*0.60}`,
          `L${C+R*0.06},${C-R*0.62}`,
          `L${C+R*0.12},${C-R*0.44}`,
          `L${C+R*0.00},${C-R*0.34}`,
          `L${C-R*0.12},${C-R*0.40}`,
          `Z`
        ].join(' ')} />
        {/* Africa */}
        <path d={[
          `M${C-R*0.02},${C-R*0.28}`,
          `L${C+R*0.26},${C-R*0.24}`,
          `L${C+R*0.38},${C+R*0.05}`,
          `L${C+R*0.26},${C+R*0.48}`,
          `L${C+R*0.05},${C+R*0.52}`,
          `L${C-R*0.12},${C+R*0.38}`,
          `L${C-R*0.08},${C+R*0.00}`,
          `Z`
        ].join(' ')} />
      </g>

      {/* Map pin — sits on the meridian */}
      <g clipPath={`url(#g${size})`}>
        <circle cx={pinX} cy={pinY} r={s * 5.5} fill={fg} />
        <circle cx={pinX} cy={pinY} r={s * 2.4} fill={inverted ? '#e0e7ff' : '#3730a3'} />
        <line
          x1={pinX} y1={pinY + s * 5.5}
          x2={pinX} y2={pinY + s * 10}
          stroke={fg} strokeWidth={s * 3} strokeLinecap="round"
        />
      </g>

      {/* CLOCK — ticks + hands, clipped to the 120° sector */}
      <g clipPath={`url(#c${size})`} strokeLinecap="round">
        {majorTicks.map((t, i) => (
          <line key={i} x1={t.x1} y1={t.y1} x2={t.x2} y2={t.y2}
            stroke={fg} strokeWidth={s * 3.2} />
        ))}
        {/* Hour hand */}
        <line x1={C} y1={C}
          x2={C + R * 0.44 * Math.cos(hA)} y2={C + R * 0.44 * Math.sin(hA)}
          stroke={fg} strokeWidth={s * 4.5} />
        {/* Minute hand */}
        <line x1={C} y1={C}
          x2={C + R * 0.66 * Math.cos(mA)} y2={C + R * 0.66 * Math.sin(mA)}
          stroke={fg} strokeWidth={s * 3} />
      </g>

      {/* Divider — two lines forming the clock sector boundary */}
      <line x1={C} y1={C - R + s} x2={C} y2={C}
        stroke={fg} strokeWidth={s * 2} strokeOpacity="0.4" />
      <line x1={C} y1={C} x2={edgeX} y2={edgeY}
        stroke={fg} strokeWidth={s * 2} strokeOpacity="0.4" />

      {/* Center dot */}
      <circle cx={C} cy={C} r={s * 4} fill={fg} />
    </svg>
  )
}

export function GlobeClock() {
  return (
    <div style={{ fontFamily: "'Space Grotesk', sans-serif" }}>

      {/* PRIMARY: inverted on indigo */}
      <div style={{
        background: 'linear-gradient(145deg, #4338ca, #3730a3)',
        padding: '64px',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        gap: '48px',
      }}>
        {/* Stacked lockup */}
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '28px' }}>
          <Mark size={140} inverted />
          <span style={{
            fontSize: '60px', fontWeight: 700, letterSpacing: '-3px',
            color: '#fff', lineHeight: 1,
          }}>timediverr</span>
        </div>

        {/* Horizontal lockup */}
        <div style={{
          display: 'flex', alignItems: 'center', gap: '14px',
          background: 'rgba(255,255,255,0.1)',
          border: '1px solid rgba(255,255,255,0.18)',
          borderRadius: '14px', padding: '14px 22px',
        }}>
          <Mark size={42} inverted />
          <span style={{ fontSize: '26px', fontWeight: 700, color: '#fff', letterSpacing: '-1px' }}>
            timediverr
          </span>
        </div>

        {/* Favicon row */}
        <div style={{ display: 'flex', alignItems: 'flex-end', gap: '20px' }}>
          {[16, 24, 32, 48].map(n => (
            <div key={n} style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '8px' }}>
              <Mark size={n} inverted />
              <span style={{ fontSize: '10px', color: '#818cf8', letterSpacing: '1px' }}>{n}</span>
            </div>
          ))}
        </div>
      </div>

      {/* SECONDARY: on white */}
      <div style={{
        background: '#f5f4ff',
        padding: '48px 64px',
        display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '32px', flexWrap: 'wrap',
      }}>
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '16px' }}>
          <Mark size={90} />
          <span style={{
            fontSize: '38px', fontWeight: 700, letterSpacing: '-2px',
            background: 'linear-gradient(135deg, #4f46e5, #3730a3)',
            WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent',
          }}>timediverr</span>
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '12px',
          background: 'linear-gradient(135deg,#4338ca,#3730a3)',
          borderRadius: '12px', padding: '12px 20px',
        }}>
          <Mark size={34} inverted />
          <span style={{ fontSize: '21px', fontWeight: 700, color: '#fff', letterSpacing: '-0.5px' }}>timediverr</span>
        </div>
      </div>
    </div>
  )
}
