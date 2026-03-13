function Mark({ size = 140, inverted = false }) {
  const C = size / 2
  const R = size / 2 - size * 0.05
  const s = size / 100 // scale unit

  const fg = inverted ? '#3730a3' : '#fff'
  const circleFill = inverted ? '#e0e7ff' : 'url(#grad)'

  // Clock: only 4 major ticks at 12/3/6/9 — right half only
  const majorTicks = [0, 3, 6, 9].map((i) => {
    const a = (i * 30 - 90) * Math.PI / 180
    return {
      x1: C + (R - s * 18) * Math.cos(a), y1: C + (R - s * 18) * Math.sin(a),
      x2: C + (R - s * 6)  * Math.cos(a), y2: C + (R - s * 6)  * Math.sin(a),
    }
  })

  // Hands: hour at 2, minute at 4
  const hA = (2 * 30 - 90) * Math.PI / 180
  const mA = (4 * 30 - 90) * Math.PI / 180

  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`} fill="none">
      <defs>
        <linearGradient id="grad" x1="0" y1="0" x2={size} y2={size} gradientUnits="userSpaceOnUse">
          <stop offset="0%" stopColor="#4f46e5" />
          <stop offset="100%" stopColor="#3730a3" />
        </linearGradient>
        <clipPath id={`l${size}`}><path d={`M${C} ${C-R} A${R} ${R} 0 1 0 ${C} ${C+R}Z`}/></clipPath>
        <clipPath id={`r${size}`}><path d={`M${C} ${C-R} A${R} ${R} 0 1 1 ${C} ${C+R}Z`}/></clipPath>
      </defs>

      {/* Circle base */}
      <circle cx={C} cy={C} r={R} fill={circleFill} />

      {/* GLOBE — equator + one meridian arc + pin */}
      <g clipPath={`url(#l${size})`} stroke={fg} fill="none" strokeLinecap="round">
        <line x1={C-R} y1={C} x2={C} y2={C} strokeWidth={s * 3.2} />
        <path d={`M${C} ${C-R} A${R*0.5} ${R} 0 0 0 ${C} ${C+R}`} strokeWidth={s * 3.2} />
      </g>

      {/* Map pin — upper left */}
      <g clipPath={`url(#l${size})`}>
        <circle cx={C - R*0.54} cy={C - R*0.48} r={s*5.5} fill={fg} />
        <circle cx={C - R*0.54} cy={C - R*0.48} r={s*2.4} fill={inverted ? '#e0e7ff' : '#3730a3'} />
        <line
          x1={C - R*0.54} y1={C - R*0.48 + s*5.5}
          x2={C - R*0.54} y2={C - R*0.48 + s*10}
          stroke={fg} strokeWidth={s*3} strokeLinecap="round"
        />
      </g>

      {/* CLOCK — 4 ticks + 2 hands */}
      <g clipPath={`url(#r${size})`} strokeLinecap="round">
        {majorTicks.map((t, i) => (
          <line key={i} x1={t.x1} y1={t.y1} x2={t.x2} y2={t.y2}
            stroke={fg} strokeWidth={s * 3.2} />
        ))}
        <line x1={C} y1={C}
          x2={C + R*0.44*Math.cos(hA)} y2={C + R*0.44*Math.sin(hA)}
          stroke={fg} strokeWidth={s * 4.5} />
        <line x1={C} y1={C}
          x2={C + R*0.66*Math.cos(mA)} y2={C + R*0.66*Math.sin(mA)}
          stroke={fg} strokeWidth={s * 3} />
      </g>

      {/* Divider */}
      <line x1={C} y1={C-R+s*2} x2={C} y2={C+R-s*2}
        stroke={fg} strokeWidth={s*2} strokeOpacity="0.35" />

      {/* Center dot */}
      <circle cx={C} cy={C} r={s*4} fill={fg} />
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
