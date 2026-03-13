function GlobeClockMark({ size = 140, inverted = false }) {
  const C = size / 2
  const R = size / 2 - size * 0.04

  const stroke = inverted ? '#3730a3' : '#fff'
  const fill = inverted ? '#eef2ff' : 'url(#bg-mark)'
  const strokeOpacity = 1

  // Only 4 major clock ticks for bold readability at small sizes
  const majorTicks = [0, 3, 6, 9].map((i) => {
    const angleDeg = i * 30 - 90
    const angleRad = angleDeg * Math.PI / 180
    const inner = R - size * 0.18
    const outer = R - size * 0.055
    return {
      x1: C + inner * Math.cos(angleRad),
      y1: C + inner * Math.sin(angleRad),
      x2: C + outer * Math.cos(angleRad),
      y2: C + outer * Math.sin(angleRad),
    }
  })

  // Minor ticks between majors (less bold)
  const minorTicks = [1, 2, 4, 5, 7, 8, 10, 11].map((i) => {
    const angleDeg = i * 30 - 90
    const angleRad = angleDeg * Math.PI / 180
    const inner = R - size * 0.12
    const outer = R - size * 0.055
    return {
      x1: C + inner * Math.cos(angleRad),
      y1: C + inner * Math.sin(angleRad),
      x2: C + outer * Math.cos(angleRad),
      y2: C + outer * Math.sin(angleRad),
    }
  })

  // Hour hand at 2 o'clock
  const hourAngle = (2 * 30 - 90) * Math.PI / 180
  const hourLen = R * 0.45

  // Minute hand at 4 o'clock
  const minuteAngle = (4 * 30 - 90) * Math.PI / 180
  const minuteLen = R * 0.65

  const sw = size * 0.025 // base stroke width unit

  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`} fill="none">
      <defs>
        <linearGradient id="bg-mark" x1="0" y1="0" x2={size} y2={size} gradientUnits="userSpaceOnUse">
          <stop offset="0%" stopColor="#4f46e5" />
          <stop offset="100%" stopColor="#3730a3" />
        </linearGradient>
        <clipPath id={`lh-${size}`}>
          <path d={`M ${C} ${C - R} A ${R} ${R} 0 1 0 ${C} ${C + R} Z`} />
        </clipPath>
        <clipPath id={`rh-${size}`}>
          <path d={`M ${C} ${C - R} A ${R} ${R} 0 1 1 ${C} ${C + R} Z`} />
        </clipPath>
      </defs>

      {/* Base circle */}
      <circle cx={C} cy={C} r={R} fill={fill} />

      {/* LEFT HALF: Globe — bold lat/lon lines */}
      <g clipPath={`url(#lh-${size})`} stroke={stroke} fill="none" strokeOpacity="0.85">
        {/* 3 bold latitude lines */}
        {[C - R * 0.55, C, C + R * 0.55].map((y, i) => (
          <line key={i} x1={C - R} y1={y} x2={C} y2={y} strokeWidth={sw * 1.3} />
        ))}
        {/* 2 bold longitude arcs */}
        <ellipse cx={C} cy={C} rx={R * 0.42} ry={R} strokeWidth={sw * 1.3} />
        <ellipse cx={C} cy={C} rx={R * 0.76} ry={R} strokeWidth={sw * 1.0} strokeOpacity="0.5" />
      </g>

      {/* Globe equator ring (full arc, subtle) */}
      <g clipPath={`url(#lh-${size})`}>
        <line x1={C - R} y1={C} x2={C} y2={C} stroke={stroke} strokeWidth={sw * 1.5} strokeOpacity="0.9" />
      </g>

      {/* LEFT: Bold map pin */}
      <g clipPath={`url(#lh-${size})`}>
        <circle cx={C - R * 0.38} cy={C - R * 0.14} r={sw * 3} fill={stroke} />
        <circle cx={C - R * 0.38} cy={C - R * 0.14} r={sw * 1.3} fill={inverted ? '#eef2ff' : '#4338ca'} />
        <line
          x1={C - R * 0.38} y1={C - R * 0.14 + sw * 3}
          x2={C - R * 0.38} y2={C - R * 0.14 + sw * 6}
          stroke={stroke} strokeWidth={sw * 1.8} strokeLinecap="round"
        />
      </g>

      {/* RIGHT HALF: Clock */}
      <g clipPath={`url(#rh-${size})`}>
        {/* Major ticks */}
        {majorTicks.map((t, i) => (
          <line key={i} x1={t.x1} y1={t.y1} x2={t.x2} y2={t.y2}
            stroke={stroke} strokeWidth={sw * 2} strokeLinecap="round" strokeOpacity="0.9" />
        ))}
        {/* Minor ticks */}
        {minorTicks.map((t, i) => (
          <line key={i} x1={t.x1} y1={t.y1} x2={t.x2} y2={t.y2}
            stroke={stroke} strokeWidth={sw * 1.0} strokeLinecap="round" strokeOpacity="0.5" />
        ))}
        {/* Hour hand — thick */}
        <line
          x1={C} y1={C}
          x2={C + hourLen * Math.cos(hourAngle)} y2={C + hourLen * Math.sin(hourAngle)}
          stroke={stroke} strokeWidth={sw * 2.6} strokeLinecap="round" strokeOpacity="0.95"
        />
        {/* Minute hand — medium */}
        <line
          x1={C} y1={C}
          x2={C + minuteLen * Math.cos(minuteAngle)} y2={C + minuteLen * Math.sin(minuteAngle)}
          stroke={stroke} strokeWidth={sw * 1.6} strokeLinecap="round" strokeOpacity="0.75"
        />
      </g>

      {/* Bold divider line */}
      <line
        x1={C} y1={C - R + sw} x2={C} y2={C + R - sw}
        stroke={stroke} strokeWidth={sw * 1.8} strokeOpacity="0.5"
      />

      {/* Center pivot */}
      <circle cx={C} cy={C} r={sw * 2.8} fill={stroke} strokeOpacity="1" />
      <circle cx={C} cy={C} r={sw * 1.2} fill={inverted ? '#eef2ff' : '#4338ca'} />

      {/* Outer ring */}
      <circle cx={C} cy={C} r={R} stroke={stroke} strokeWidth={sw * 1.5} strokeOpacity="0.2" />
    </svg>
  )
}

export function GlobeClock() {
  return (
    <div style={{ fontFamily: "'Space Grotesk', sans-serif" }}>

      {/* === PRIMARY: Inverted on indigo === */}
      <div style={{
        background: 'linear-gradient(145deg, #4338ca 0%, #3730a3 100%)',
        padding: '56px 64px',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        gap: '40px',
      }}>
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '24px' }}>
          <GlobeClockMark size={140} inverted />
          <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '4px' }}>
            <span style={{
              fontSize: '58px', fontWeight: 700, letterSpacing: '-3px',
              color: '#fff', lineHeight: 1,
            }}>timediverr</span>
            <span style={{
              fontSize: '12px', fontWeight: 500, letterSpacing: '5px',
              color: '#a5b4fc', textTransform: 'uppercase',
            }}>map the past</span>
          </div>
        </div>

        {/* Horizontal lockup */}
        <div style={{
          display: 'flex', alignItems: 'center', gap: '14px',
          background: 'rgba(255,255,255,0.1)', borderRadius: '14px', padding: '14px 24px',
          border: '1px solid rgba(255,255,255,0.15)',
        }}>
          <GlobeClockMark size={40} inverted />
          <span style={{ fontSize: '26px', fontWeight: 700, letterSpacing: '-1px', color: '#fff' }}>
            timediverr
          </span>
        </div>

        {/* Favicon simulation */}
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '10px' }}>
          <span style={{ fontSize: '11px', color: '#818cf8', letterSpacing: '2px', textTransform: 'uppercase' }}>favicon sizes</span>
          <div style={{ display: 'flex', alignItems: 'flex-end', gap: '16px' }}>
            {[16, 24, 32, 48].map(s => (
              <div key={s} style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '6px' }}>
                <GlobeClockMark size={s} inverted />
                <span style={{ fontSize: '10px', color: '#6366f1' }}>{s}px</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* === SECONDARY: Standard on white === */}
      <div style={{
        background: '#f5f4ff',
        padding: '48px 64px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        gap: '40px',
        flexWrap: 'wrap',
      }}>
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '16px' }}>
          <GlobeClockMark size={100} />
          <span style={{
            fontSize: '40px', fontWeight: 700, letterSpacing: '-2px',
            background: 'linear-gradient(135deg, #4f46e5, #3730a3)',
            WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent',
          }}>timediverr</span>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
          {/* Dark pill */}
          <div style={{
            display: 'flex', alignItems: 'center', gap: '12px',
            background: 'linear-gradient(135deg, #4338ca, #3730a3)',
            borderRadius: '12px', padding: '12px 20px',
          }}>
            <GlobeClockMark size={36} inverted />
            <span style={{ fontSize: '22px', fontWeight: 700, color: '#fff', letterSpacing: '-0.5px' }}>timediverr</span>
          </div>
          {/* Light pill */}
          <div style={{
            display: 'flex', alignItems: 'center', gap: '12px',
            background: '#fff', borderRadius: '12px', padding: '12px 20px',
            border: '2px solid #e0e7ff',
          }}>
            <GlobeClockMark size={36} />
            <span style={{ fontSize: '22px', fontWeight: 700, color: '#3730a3', letterSpacing: '-0.5px' }}>timediverr</span>
          </div>
        </div>
      </div>
    </div>
  )
}
