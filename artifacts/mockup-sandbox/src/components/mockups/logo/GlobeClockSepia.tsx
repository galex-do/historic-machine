// Variation B — "Antique Ink" : warm parchment + ink-style hatching + sienna palette
function Mark({ size = 140 }: { size?: number }) {
  const C = size / 2
  const R = size / 2 - size * 0.05
  const s = size / 100
  const id = `s${size}`

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

  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`} fill="none">
      <defs>
        {/* Parchment radial */}
        <radialGradient id={`pg${id}`} cx="42%" cy="36%" r="70%">
          <stop offset="0%" stopColor="#fdf6e3" />
          <stop offset="60%" stopColor="#f5ead0" />
          <stop offset="100%" stopColor="#e8d5a8" />
        </radialGradient>
        {/* Ocean hatching pattern */}
        <pattern id={`hatch${id}`} patternUnits="userSpaceOnUse" width={s * 4} height={s * 4} patternTransform="rotate(35)">
          <line x1="0" y1="0" x2="0" y2={s * 4} stroke="#b8956a" strokeWidth={s * 0.55} strokeOpacity="0.4"/>
        </pattern>
        <clipPath id={`gg${id}`}><path d={globeClip} /></clipPath>
        <clipPath id={`cc${id}`}><path d={clockClip} /></clipPath>
      </defs>

      {/* Parchment base */}
      <circle cx={C} cy={C} r={R} fill={`url(#pg${id})`} />
      {/* Ocean hatch fill in globe area */}
      <circle cx={C} cy={C} r={R} fill={`url(#hatch${id})`} clipPath={`url(#gg${id})`} />

      {/* Continent fills — sienna ink */}
      <g clipPath={`url(#gg${id})`} fill="#7c4a1e" fillOpacity="0.82" stroke="#5a3010" strokeWidth={s*0.8} strokeLinejoin="round">
        <path d={poly([[-0.56,-0.76],[-0.38,-0.84],[-0.22,-0.82],[-0.18,-0.72],[-0.26,-0.62],[-0.44,-0.60],[-0.58,-0.64],[-0.62,-0.72]])} />
        <path d={poly([[-0.86,-0.50],[-0.72,-0.70],[-0.55,-0.80],[-0.38,-0.82],[-0.22,-0.78],[-0.14,-0.70],[-0.06,-0.68],[-0.08,-0.58],[-0.14,-0.50],[-0.18,-0.42],[-0.20,-0.30],[-0.20,-0.18],[-0.22,-0.06],[-0.26,0.00],[-0.36,0.04],[-0.44,0.02],[-0.50,-0.06],[-0.54,-0.16],[-0.50,-0.28],[-0.44,-0.20],[-0.54,-0.06],[-0.62,-0.18],[-0.70,-0.30],[-0.78,-0.40],[-0.84,-0.48]])} />
        <path d={poly([[-0.30,0.06],[-0.14,0.00],[-0.06,0.06],[-0.02,0.16],[-0.02,0.30],[-0.04,0.44],[-0.10,0.56],[-0.18,0.66],[-0.26,0.76],[-0.34,0.76],[-0.40,0.68],[-0.46,0.54],[-0.50,0.38],[-0.50,0.24],[-0.46,0.10],[-0.40,0.04]])} />
        <path d={poly([[-0.20,-0.56],[-0.08,-0.62],[0.00,-0.60],[0.04,-0.50],[0.00,-0.40],[-0.06,-0.36],[-0.14,-0.40],[-0.18,-0.50]])} />
        <path d={poly([[-0.08,-0.28],[0.08,-0.30],[0.20,-0.24],[0.30,-0.10],[0.36,0.06],[0.34,0.22],[0.28,0.40],[0.20,0.54],[0.08,0.62],[-0.06,0.64],[-0.18,0.56],[-0.28,0.42],[-0.32,0.26],[-0.26,0.10],[-0.24,0.02],[-0.28,-0.08],[-0.22,-0.20],[-0.12,-0.26]])} />
      </g>

      {/* Globe latitude arcs — faint ink lines */}
      <g clipPath={`url(#gg${id})`} stroke="#b8956a" strokeWidth={s*0.7} strokeOpacity="0.35" fill="none">
        {[-0.5,-0.25,0,0.25,0.5].map((ny,i) => {
          const lat_r = Math.sqrt(1 - ny*ny) * R
          return <ellipse key={i} cx={C} cy={C+ny*R} rx={lat_r} ry={lat_r * 0.25} />
        })}
      </g>

      {/* Map pin — dark ink with sepia ring */}
      <g clipPath={`url(#gg${id})`}>
        <circle cx={pinX} cy={pinY} r={s*5.5} fill="#3d1f08" />
        <circle cx={pinX} cy={pinY} r={s*2.4} fill="#e8d5a8" />
        <line x1={pinX} y1={pinY+s*5.5} x2={pinX} y2={pinY+s*10} stroke="#3d1f08" strokeWidth={s*3} strokeLinecap="round"/>
      </g>

      {/* Clock sector — dark brown with cross ticks */}
      <g clipPath={`url(#cc${id})`} strokeLinecap="round">
        {/* Sector fill — slightly darker parchment */}
        <path d={clockClip} fill="#e8d2a0" fillOpacity="0.5" />
        {majorTicks.map((t,i) => (
          <line key={i} x1={t.x1} y1={t.y1} x2={t.x2} y2={t.y2} stroke="#3d1f08" strokeWidth={s*3.2} />
        ))}
        <line x1={C} y1={C} x2={C+R*0.44*Math.cos(hA)} y2={C+R*0.44*Math.sin(hA)} stroke="#3d1f08" strokeWidth={s*4.5} />
        <line x1={C} y1={C} x2={C+R*0.66*Math.cos(mA)} y2={C+R*0.66*Math.sin(mA)} stroke="#5a3010" strokeWidth={s*3} />
      </g>

      {/* Outer border ring — ink */}
      <circle cx={C} cy={C} r={R} stroke="#5a3010" strokeWidth={s*2} />

      {/* Divider */}
      <line x1={C} y1={C-R+s} x2={C} y2={C} stroke="#5a3010" strokeWidth={s*2} strokeOpacity="0.6" />
      <line x1={C} y1={C} x2={edgeX} y2={edgeY} stroke="#5a3010" strokeWidth={s*2} strokeOpacity="0.6" />

      {/* Center dot */}
      <circle cx={C} cy={C} r={s*4} fill="#3d1f08" />
    </svg>
  )
}

export function GlobeClockSepia() {
  return (
    <div style={{ fontFamily: "'Space Grotesk', sans-serif", background: '#f0e6cc', minHeight: '100%', padding: '48px', display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '40px' }}>
      <div style={{ color: '#a07848', fontSize: '11px', letterSpacing: '3px', textTransform: 'uppercase' }}>B — Antique Ink</div>
      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '24px' }}>
        <Mark size={140} />
        <span style={{ fontSize: '52px', fontWeight: 700, letterSpacing: '-2.5px', color: '#3d1f08' }}>timediverr</span>
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: '12px', background: '#fdf6e3', border: '1.5px solid #b8956a', borderRadius: '12px', padding: '12px 20px' }}>
        <Mark size={38} />
        <span style={{ fontSize: '22px', fontWeight: 700, color: '#3d1f08', letterSpacing: '-1px' }}>timediverr</span>
      </div>
      <div style={{ display: 'flex', alignItems: 'flex-end', gap: '20px' }}>
        {[16,24,32,48].map(n => (
          <div key={n} style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '8px' }}>
            <Mark size={n} />
            <span style={{ fontSize: '10px', color: '#a07848', letterSpacing: '1px' }}>{n}</span>
          </div>
        ))}
      </div>
    </div>
  )
}
