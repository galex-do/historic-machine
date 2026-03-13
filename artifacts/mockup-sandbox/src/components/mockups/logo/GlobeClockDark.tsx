// Variation A — "Starfield" : deep void + dual-tone (cyan globe / amber clock) + glow
function Mark({ size = 140 }: { size?: number }) {
  const C = size / 2
  const R = size / 2 - size * 0.05
  const s = size / 100

  const mA = (4 * 30 - 90) * Math.PI / 180
  const edgeX = C + R * Math.cos(mA)
  const edgeY = C + R * Math.sin(mA)
  const hA = (2 * 30 - 90) * Math.PI / 180

  const globeClip = `M${C} ${C - R} A${R} ${R} 0 1 0 ${edgeX} ${edgeY} L${C} ${C} Z`
  const clockClip = `M${C} ${C - R} A${R} ${R} 0 0 1 ${edgeX} ${edgeY} L${C} ${C} Z`

  const majorTicks = [0, 1, 2, 3, 4].map((i) => {
    const a = (i * 30 - 90) * Math.PI / 180
    return {
      x1: C + (R - s * 18) * Math.cos(a), y1: C + (R - s * 18) * Math.sin(a),
      x2: C + (R - s * 6)  * Math.cos(a), y2: C + (R - s * 6)  * Math.sin(a),
    }
  })

  const pinX = C - R * 0.44
  const pinY = C - R * 0.48

  const q = (nx: number, ny: number) => `${C + nx * R} ${C + ny * R}`
  const poly = (pts: [number,number][]) => `M${pts.map(([x,y]) => q(x,y)).join('L')}Z`

  // Starfield dots
  const stars = [
    [-0.62,0.28],[-0.75,0.12],[-0.52,0.55],[-0.80,0.42],[-0.35,0.60],
    [-0.68,-0.22],[-0.82,-0.18],[-0.38,0.72],[-0.25,0.45],[-0.72,0.65],
    [-0.88,0.00],[-0.40,-0.62],[-0.18,0.80],
  ] as [number,number][]

  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`} fill="none">
      <defs>
        {/* Deep space background */}
        <radialGradient id={`bg${size}`} cx="40%" cy="35%" r="65%">
          <stop offset="0%" stopColor="#0f172a" />
          <stop offset="100%" stopColor="#020617" />
        </radialGradient>
        {/* Cyan glow */}
        <filter id={`glow${size}`} x="-30%" y="-30%" width="160%" height="160%">
          <feGaussianBlur stdDeviation={s * 1.5} result="blur"/>
          <feMerge><feMergeNode in="blur"/><feMergeNode in="SourceGraphic"/></feMerge>
        </filter>
        {/* Amber glow */}
        <filter id={`aglow${size}`} x="-40%" y="-40%" width="180%" height="180%">
          <feGaussianBlur stdDeviation={s * 1.8} result="blur"/>
          <feMerge><feMergeNode in="blur"/><feMergeNode in="SourceGraphic"/></feMerge>
        </filter>
        <clipPath id={`gd${size}`}><path d={globeClip} /></clipPath>
        <clipPath id={`cd${size}`}><path d={clockClip} /></clipPath>
      </defs>

      {/* Base circle — deep space */}
      <circle cx={C} cy={C} r={R} fill={`url(#bg${size})`} />
      <circle cx={C} cy={C} r={R} stroke="#1e293b" strokeWidth={s * 1.5} />

      {/* Star dots in globe area */}
      <g clipPath={`url(#gd${size})`}>
        {stars.map(([nx,ny],i) => (
          <circle key={i} cx={C+nx*R} cy={C+ny*R} r={s*(0.6+Math.sin(i*1.7)*0.4)} fill="#94a3b8" fillOpacity="0.6" />
        ))}
      </g>

      {/* Globe — cyan continent outlines */}
      <g clipPath={`url(#gd${size})`} fill="none" stroke="#22d3ee" strokeWidth={s * 1.8}
         filter={`url(#glow${size})`} strokeLinejoin="round">
                <path d={poly([[-0.52,-0.76],[-0.34,-0.86],[-0.20,-0.84],[-0.16,-0.72],[-0.26,-0.62],[-0.44,-0.60],[-0.58,-0.66]])} />
        <path d={poly([[-0.88,-0.48],[-0.70,-0.68],[-0.50,-0.80],[-0.32,-0.86],[-0.18,-0.82],[-0.06,-0.72],[-0.12,-0.58],[-0.16,-0.46],[-0.18,-0.32],[-0.20,-0.16],[-0.22,-0.04],[-0.26,0.04],[-0.38,0.06],[-0.46,0.04],[-0.54,-0.04],[-0.60,-0.16],[-0.66,-0.26],[-0.74,-0.36],[-0.82,-0.44]])} />
        <path d={poly([[-0.36,0.08],[-0.18,0.02],[-0.08,0.06],[0.02,0.12],[0.08,0.28],[0.06,0.46],[0.00,0.56],[-0.10,0.64],[-0.20,0.72],[-0.28,0.78],[-0.34,0.82],[-0.40,0.76],[-0.46,0.62],[-0.50,0.46],[-0.48,0.28],[-0.42,0.10]])} />
        <path d={poly([[-0.22,-0.60],[-0.06,-0.66],[0.04,-0.60],[0.08,-0.50],[0.04,-0.40],[-0.04,-0.36],[-0.12,-0.42],[-0.18,-0.52]])} />
        <path d={poly([[0.16,-0.26],[0.26,-0.18],[0.36,-0.04],[0.40,0.12],[0.34,0.30],[0.26,0.46],[0.18,0.60],[0.06,0.68],[-0.06,0.70],[-0.14,0.60],[0.16,0.46],[0.18,0.30],[0.16,0.14],[0.16,-0.06],[0.14,-0.26]])} />
      </g>

      {/* Map pin — glowing red */}
      <g clipPath={`url(#gd${size})`} filter={`url(#aglow${size})`}>
        <circle cx={pinX} cy={pinY} r={s*5.5} fill="#f87171" />
        <circle cx={pinX} cy={pinY} r={s*2.4} fill="#7f1d1d" />
        <line x1={pinX} y1={pinY+s*5.5} x2={pinX} y2={pinY+s*10} stroke="#f87171" strokeWidth={s*3} strokeLinecap="round"/>
      </g>

      {/* Clock sector — amber */}
      <g clipPath={`url(#cd${size})`} strokeLinecap="round" filter={`url(#aglow${size})`}>
        {majorTicks.map((t,i) => (
          <line key={i} x1={t.x1} y1={t.y1} x2={t.x2} y2={t.y2} stroke="#fbbf24" strokeWidth={s*3.2} />
        ))}
        <line x1={C} y1={C} x2={C+R*0.44*Math.cos(hA)} y2={C+R*0.44*Math.sin(hA)} stroke="#fbbf24" strokeWidth={s*4.5} />
        <line x1={C} y1={C} x2={C+R*0.66*Math.cos(mA)} y2={C+R*0.66*Math.sin(mA)} stroke="#f59e0b" strokeWidth={s*3} />
      </g>

      {/* Divider */}
      <line x1={C} y1={C-R+s} x2={C} y2={C} stroke="#334155" strokeWidth={s*2} />
      <line x1={C} y1={C} x2={edgeX} y2={edgeY} stroke="#334155" strokeWidth={s*2} />

      {/* Center dot */}
      <circle cx={C} cy={C} r={s*4} fill="#fbbf24" />
    </svg>
  )
}

export function GlobeClockDark() {
  return (
    <div style={{ fontFamily: "'Space Grotesk', sans-serif", background: '#020617', minHeight: '100%', padding: '48px', display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '40px' }}>
      <div style={{ color: '#64748b', fontSize: '11px', letterSpacing: '3px', textTransform: 'uppercase' }}>A — Starfield / Dual-tone</div>
      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '24px' }}>
        <Mark size={140} />
        <span style={{ fontSize: '52px', fontWeight: 700, letterSpacing: '-2.5px', background: 'linear-gradient(90deg,#22d3ee,#f59e0b)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>timediverr</span>
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: '12px', background: '#0f172a', border: '1px solid #1e293b', borderRadius: '12px', padding: '12px 20px' }}>
        <Mark size={38} />
        <span style={{ fontSize: '22px', fontWeight: 700, color: '#e2e8f0', letterSpacing: '-1px' }}>timediverr</span>
      </div>
      <div style={{ display: 'flex', alignItems: 'flex-end', gap: '20px' }}>
        {[16,24,32,48].map(n => (
          <div key={n} style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '8px' }}>
            <Mark size={n} />
            <span style={{ fontSize: '10px', color: '#475569', letterSpacing: '1px' }}>{n}</span>
          </div>
        ))}
      </div>
    </div>
  )
}
