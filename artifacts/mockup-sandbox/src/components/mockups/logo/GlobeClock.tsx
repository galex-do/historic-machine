export function GlobeClock() {
  const C = 70
  const R = 58

  // Clock tick marks — all 12, clipped to right half
  const clockTicks = Array.from({ length: 12 }, (_, i) => {
    const angleDeg = i * 30 - 90
    const angleRad = angleDeg * Math.PI / 180
    const isMajor = i % 3 === 0
    const inner = isMajor ? R - 16 : R - 11
    const outer = R - 4
    return {
      x1: C + inner * Math.cos(angleRad),
      y1: C + inner * Math.sin(angleRad),
      x2: C + outer * Math.cos(angleRad),
      y2: C + outer * Math.sin(angleRad),
      major: isMajor,
    }
  })

  // Hour hand: 2 o'clock
  const hourAngle = (2 * 30 - 90) * Math.PI / 180
  const hourLen = 28
  // Minute hand: 4 o'clock
  const minuteAngle = (4 * 30 - 90) * Math.PI / 180
  const minuteLen = 40

  return (
    <div style={{
      minHeight: '100vh',
      background: '#f5f4ff',
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      gap: '44px',
      padding: '32px',
      fontFamily: "'Space Grotesk', sans-serif",
    }}>

      {/* Main mark + wordmark stacked */}
      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '28px' }}>

        {/* The icon mark */}
        <svg width={C * 2} height={C * 2} viewBox={`0 0 ${C * 2} ${C * 2}`} fill="none">
          <defs>
            <linearGradient id="bg" x1="0" y1="0" x2={C * 2} y2={C * 2} gradientUnits="userSpaceOnUse">
              <stop offset="0%" stopColor="#4f46e5" />
              <stop offset="100%" stopColor="#3730a3" />
            </linearGradient>

            {/* Left semicircle clip */}
            <clipPath id="lhalf">
              <path d={`M ${C} ${C - R} A ${R} ${R} 0 1 0 ${C} ${C + R} Z`} />
            </clipPath>

            {/* Right semicircle clip */}
            <clipPath id="rhalf">
              <path d={`M ${C} ${C - R} A ${R} ${R} 0 1 1 ${C} ${C + R} Z`} />
            </clipPath>

            {/* Full circle clip for inner elements */}
            <clipPath id="circle">
              <circle cx={C} cy={C} r={R} />
            </clipPath>
          </defs>

          {/* Base circle */}
          <circle cx={C} cy={C} r={R} fill="url(#bg)" />

          {/* === LEFT HALF: Globe === */}
          <g clipPath="url(#lhalf)" stroke="rgba(255,255,255,0.22)" fill="none">
            {/* Latitude lines */}
            {[C - 36, C - 20, C, C + 20, C + 36].map((y, i) => (
              <line key={i} x1={C - R} y1={y} x2={C} y2={y} strokeWidth="1.2" />
            ))}
            {/* Longitude arcs (vertical ellipses) */}
            <ellipse cx={C} cy={C} rx={14} ry={R} strokeWidth="1.2" />
            <ellipse cx={C} cy={C} rx={30} ry={R} strokeWidth="1.2" />
            <ellipse cx={C} cy={C} rx={46} ry={R} strokeWidth="1.2" />
          </g>

          {/* Map pin on globe (equator, left side) */}
          <g clipPath="url(#lhalf)">
            <circle cx={C - 22} cy={C - 6} r={6} fill="rgba(255,255,255,0.9)" />
            <circle cx={C - 22} cy={C - 6} r={2.8} fill="#4338ca" />
            <line
              x1={C - 22} y1={C - 0.5} x2={C - 22} y2={C + 7}
              stroke="rgba(255,255,255,0.9)" strokeWidth="2" strokeLinecap="round"
            />
          </g>

          {/* === RIGHT HALF: Clock === */}
          <g clipPath="url(#rhalf)">
            {/* Tick marks */}
            {clockTicks.map((t, i) => (
              <line
                key={i}
                x1={t.x1} y1={t.y1} x2={t.x2} y2={t.y2}
                stroke={t.major ? 'rgba(255,255,255,0.85)' : 'rgba(255,255,255,0.45)'}
                strokeWidth={t.major ? 2.5 : 1.5}
                strokeLinecap="round"
              />
            ))}

            {/* Hour hand */}
            <line
              x1={C} y1={C}
              x2={C + hourLen * Math.cos(hourAngle)}
              y2={C + hourLen * Math.sin(hourAngle)}
              stroke="rgba(255,255,255,0.95)"
              strokeWidth="4"
              strokeLinecap="round"
            />

            {/* Minute hand */}
            <line
              x1={C} y1={C}
              x2={C + minuteLen * Math.cos(minuteAngle)}
              y2={C + minuteLen * Math.sin(minuteAngle)}
              stroke="rgba(255,255,255,0.75)"
              strokeWidth="2.5"
              strokeLinecap="round"
            />
          </g>

          {/* Center pivot dot */}
          <circle cx={C} cy={C} r={4.5} fill="#fff" />
          <circle cx={C} cy={C} r={2} fill="#4338ca" />

          {/* Divider line */}
          <line
            x1={C} y1={C - R + 1}
            x2={C} y2={C + R - 1}
            stroke="rgba(255,255,255,0.35)"
            strokeWidth="1.5"
          />

          {/* Outer ring */}
          <circle cx={C} cy={C} r={R} stroke="rgba(255,255,255,0.18)" strokeWidth="1.5" />
        </svg>

        {/* Wordmark */}
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '3px' }}>
          <span style={{
            fontSize: '56px',
            fontWeight: 700,
            letterSpacing: '-2.5px',
            background: 'linear-gradient(135deg, #4f46e5 0%, #3730a3 100%)',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
            lineHeight: 1,
          }}>
            timediverr
          </span>
          <span style={{
            fontSize: '12px',
            fontWeight: 500,
            letterSpacing: '5px',
            color: '#818cf8',
            textTransform: 'uppercase',
          }}>
            map the past
          </span>
        </div>
      </div>

      {/* Horizontal lockups row */}
      <div style={{ display: 'flex', gap: '24px', flexWrap: 'wrap', justifyContent: 'center' }}>

        {/* Dark lockup */}
        <div style={{
          display: 'flex', alignItems: 'center', gap: '14px',
          background: 'linear-gradient(135deg, #4338ca, #3730a3)',
          borderRadius: '14px', padding: '14px 24px',
        }}>
          <svg width="38" height="38" viewBox={`0 0 ${C * 2} ${C * 2}`} fill="none">
            <defs>
              <clipPath id="lh2"><path d={`M ${C} ${C - R} A ${R} ${R} 0 1 0 ${C} ${C + R} Z`} /></clipPath>
              <clipPath id="rh2"><path d={`M ${C} ${C - R} A ${R} ${R} 0 1 1 ${C} ${C + R} Z`} /></clipPath>
            </defs>
            <circle cx={C} cy={C} r={R} fill="rgba(255,255,255,0.15)" />
            <g clipPath="url(#lh2)" stroke="rgba(255,255,255,0.3)" fill="none">
              {[C - 30, C, C + 30].map((y, i) => <line key={i} x1={C - R} y1={y} x2={C} y2={y} strokeWidth="1.5" />)}
              <ellipse cx={C} cy={C} rx={28} ry={R} strokeWidth="1.5" />
            </g>
            <g clipPath="url(#rh2)">
              {clockTicks.filter((_, i) => [0, 3, 6, 9].includes(i)).map((t, i) => (
                <line key={i} x1={t.x1} y1={t.y1} x2={t.x2} y2={t.y2} stroke="rgba(255,255,255,0.8)" strokeWidth="2.5" strokeLinecap="round" />
              ))}
              <line x1={C} y1={C} x2={C + hourLen * Math.cos(hourAngle)} y2={C + hourLen * Math.sin(hourAngle)} stroke="#fff" strokeWidth="4.5" strokeLinecap="round" />
              <line x1={C} y1={C} x2={C + minuteLen * Math.cos(minuteAngle)} y2={C + minuteLen * Math.sin(minuteAngle)} stroke="rgba(255,255,255,0.75)" strokeWidth="3" strokeLinecap="round" />
            </g>
            <circle cx={C} cy={C} r={5} fill="#fff" />
            <line x1={C} y1={C - R + 1} x2={C} y2={C + R - 1} stroke="rgba(255,255,255,0.35)" strokeWidth="1.5" />
            <circle cx={C} cy={C} r={R} stroke="rgba(255,255,255,0.2)" strokeWidth="1.5" />
          </svg>
          <span style={{ fontSize: '26px', fontWeight: 700, letterSpacing: '-1px', color: '#fff' }}>timediverr</span>
        </div>

        {/* Light lockup */}
        <div style={{
          display: 'flex', alignItems: 'center', gap: '14px',
          border: '2px solid #e0e7ff', borderRadius: '14px', padding: '14px 24px',
          background: '#fff',
        }}>
          <svg width="38" height="38" viewBox={`0 0 ${C * 2} ${C * 2}`} fill="none">
            <defs>
              <clipPath id="lh3"><path d={`M ${C} ${C - R} A ${R} ${R} 0 1 0 ${C} ${C + R} Z`} /></clipPath>
              <clipPath id="rh3"><path d={`M ${C} ${C - R} A ${R} ${R} 0 1 1 ${C} ${C + R} Z`} /></clipPath>
            </defs>
            <circle cx={C} cy={C} r={R} fill="url(#bg)" />
            <g clipPath="url(#lh3)" stroke="rgba(255,255,255,0.25)" fill="none">
              {[C - 30, C, C + 30].map((y, i) => <line key={i} x1={C - R} y1={y} x2={C} y2={y} strokeWidth="1.5" />)}
              <ellipse cx={C} cy={C} rx={28} ry={R} strokeWidth="1.5" />
            </g>
            <g clipPath="url(#rh3)">
              {clockTicks.filter((_, i) => [0, 3, 6, 9].includes(i)).map((t, i) => (
                <line key={i} x1={t.x1} y1={t.y1} x2={t.x2} y2={t.y2} stroke="rgba(255,255,255,0.8)" strokeWidth="2.5" strokeLinecap="round" />
              ))}
              <line x1={C} y1={C} x2={C + hourLen * Math.cos(hourAngle)} y2={C + hourLen * Math.sin(hourAngle)} stroke="#fff" strokeWidth="4.5" strokeLinecap="round" />
              <line x1={C} y1={C} x2={C + minuteLen * Math.cos(minuteAngle)} y2={C + minuteLen * Math.sin(minuteAngle)} stroke="rgba(255,255,255,0.75)" strokeWidth="3" strokeLinecap="round" />
            </g>
            <circle cx={C} cy={C} r={5} fill="#fff" />
            <line x1={C} y1={C - R + 1} x2={C} y2={C + R - 1} stroke="rgba(255,255,255,0.35)" strokeWidth="1.5" />
            <circle cx={C} cy={C} r={R} stroke="rgba(255,255,255,0.2)" strokeWidth="1.5" />
          </svg>
          <span style={{ fontSize: '26px', fontWeight: 700, letterSpacing: '-1px', color: '#3730a3' }}>timediverr</span>
        </div>
      </div>
    </div>
  )
}
