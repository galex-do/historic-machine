// Variation C — "Glass Sphere" : radial gradient depth + translucent continent fills + metallic divider
function Mark({ size = 140 }: { size?: number }) {
  const C = size / 2
  const R = size / 2 - size * 0.05
  const s = size / 100
  const id = `r${size}`

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

  // Highlight center: upper-left at ~(-0.25, -0.30) of radius
  const hx = C - R * 0.25
  const hy = C - R * 0.30

  const q = (nx: number, ny: number) => `${C + nx * R} ${C + ny * R}`
  const poly = (pts: [number,number][]) => `M${pts.map(([x,y]) => q(x,y)).join('L')}Z`

  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`} fill="none">
      <defs>
        {/* Sphere base: deep indigo at center, darker at edges */}
        <radialGradient id={`sph${id}`} cx={`${hx/size*100}%`} cy={`${hy/size*100}%`} r="80%" gradientUnits="userSpaceOnUse">
          <stop offset="0%"  stopColor="#6366f1" />
          <stop offset="40%" stopColor="#4338ca" />
          <stop offset="100%" stopColor="#1e1b4b" />
        </radialGradient>
        {/* Specular highlight */}
        <radialGradient id={`spec${id}`} cx={`${hx/size*100}%`} cy={`${hy/size*100}%`} r="38%" gradientUnits="userSpaceOnUse">
          <stop offset="0%"  stopColor="#fff" stopOpacity="0.30" />
          <stop offset="100%" stopColor="#fff" stopOpacity="0" />
        </radialGradient>
        {/* Shadow rim */}
        <radialGradient id={`rim${id}`} cx="55%" cy="60%" r="55%">
          <stop offset="0%"  stopColor="#000" stopOpacity="0" />
          <stop offset="100%" stopColor="#000" stopOpacity="0.45" />
        </radialGradient>
        {/* Clock sector fill — slightly lighter sphere tone */}
        <radialGradient id={`clk${id}`} cx={`${hx/size*100}%`} cy={`${hy/size*100}%`} r="80%" gradientUnits="userSpaceOnUse">
          <stop offset="0%"  stopColor="#818cf8" />
          <stop offset="60%" stopColor="#4f46e5" />
          <stop offset="100%" stopColor="#1e1b4b" />
        </radialGradient>
        {/* Drop shadow filter for continents */}
        <filter id={`dsf${id}`} x="-10%" y="-10%" width="120%" height="120%">
          <feDropShadow dx={s*0.5} dy={s*0.5} stdDeviation={s*1.2} floodColor="#1e1b4b" floodOpacity="0.6"/>
        </filter>
        {/* Metallic divider gradient */}
        <linearGradient id={`met${id}`} x1="0" y1="0" x2="1" y2="1">
          <stop offset="0%" stopColor="#c7d2fe" stopOpacity="0.9"/>
          <stop offset="50%" stopColor="#fff" stopOpacity="0.5"/>
          <stop offset="100%" stopColor="#818cf8" stopOpacity="0.7"/>
        </linearGradient>
        <clipPath id={`gR${id}`}><path d={globeClip} /></clipPath>
        <clipPath id={`cR${id}`}><path d={clockClip} /></clipPath>
      </defs>

      {/* Sphere body */}
      <circle cx={C} cy={C} r={R} fill={`url(#sph${id})`} />
      {/* Rim darkening */}
      <circle cx={C} cy={C} r={R} fill={`url(#rim${id})`} />

      {/* Clock sector — subtly brighter tone */}
      <path d={clockClip} fill={`url(#clk${id})`} fillOpacity="0.65" />

      {/* Globe continents — white with soft shadow for depth */}
      <g clipPath={`url(#gR${id})`} fill="rgba(255,255,255,0.80)" stroke="rgba(255,255,255,0.25)" strokeWidth={s*0.6}
         filter={`url(#dsf${id})`} strokeLinejoin="round">
                         <path d={poly([[-0.063,-0.970],[-0.021,-0.970],[0.027,-0.951],[0.068,-0.921],[0.034,-0.875],[-0.121,-0.866],[-0.152,-0.914],[-0.106,-0.951]])} />
                 <path d={poly([[-0.358,-0.934],[-0.265,-0.956],[-0.324,-0.906],[-0.266,-0.731],[-0.354,-0.707],[-0.433,-0.695],[-0.589,-0.574],[-0.704,-0.423],[-0.718,-0.5],[-0.757,-0.485],[-0.827,-0.438],[-0.824,-0.309],[-0.871,-0.276],[-0.907,-0.391],[-0.799,-0.602],[-0.655,-0.755],[-0.528,-0.839]])} />
                 <path d={poly([[-0.743,-0.017],[-0.724,-0.139],[-0.536,-0.174],[-0.374,-0.07],[-0.086,0.139],[-0.136,0.225],[-0.207,0.391],[-0.31,0.5],[-0.389,0.559],[-0.344,0.829],[-0.435,0.788],[-0.47,0.682],[-0.766,0.035]])} />
                 <path d={poly([[-0.032,-0.242],[0.328,-0.53],[0.509,-0.375],[0.578,0.407],[0.468,0.5],[0.378,0.574],[0.366,0.559],[0.394,0.375],[0.396,0.259],[0.418,0.035],[0.293,-0.087],[0.022,-0.174]])} />
         <path d={poly([[0.16,-0.26],[0.26,-0.18],[0.36,-0.04],[0.40,0.12],[0.34,0.30],[0.26,0.46],[0.18,0.60],[0.06,0.68],[-0.06,0.70],[-0.14,0.60],[0.16,0.46],[0.18,0.30],[0.16,0.14],[0.16,-0.06],[0.14,-0.26]])} />
      </g>

      {/* Specular highlight over everything */}
      <circle cx={C} cy={C} r={R} fill={`url(#spec${id})`} />

      {/* Map pin */}
      <g clipPath={`url(#gR${id})`}>
        <circle cx={pinX} cy={pinY} r={s*5.5} fill="#fff" fillOpacity="0.95" />
        <circle cx={pinX} cy={pinY} r={s*2.4} fill="#4338ca" />
        <line x1={pinX} y1={pinY+s*5.5} x2={pinX} y2={pinY+s*10} stroke="#fff" strokeWidth={s*3} strokeLinecap="round" strokeOpacity="0.9"/>
      </g>

      {/* Clock sector ticks and hands */}
      <g clipPath={`url(#cR${id})`} strokeLinecap="round">
        {majorTicks.map((t,i) => (
          <line key={i} x1={t.x1} y1={t.y1} x2={t.x2} y2={t.y2} stroke="#e0e7ff" strokeWidth={s*3.2} strokeOpacity="0.9" />
        ))}
        <line x1={C} y1={C} x2={C+R*0.44*Math.cos(hA)} y2={C+R*0.44*Math.sin(hA)} stroke="#fff" strokeWidth={s*4.5} />
        <line x1={C} y1={C} x2={C+R*0.66*Math.cos(mA)} y2={C+R*0.66*Math.sin(mA)} stroke="#c7d2fe" strokeWidth={s*3} />
      </g>

      {/* Metallic divider */}
      <line x1={C} y1={C-R+s} x2={C} y2={C} stroke={`url(#met${id})`} strokeWidth={s*2.5} />
      <line x1={C} y1={C} x2={edgeX} y2={edgeY} stroke={`url(#met${id})`} strokeWidth={s*2.5} />

      {/* Outer gloss ring */}
      <circle cx={C} cy={C} r={R-s*0.5} stroke="rgba(255,255,255,0.20)" strokeWidth={s*1.5} />

      {/* Center dot */}
      <circle cx={C} cy={C} r={s*4} fill="#e0e7ff" />
    </svg>
  )
}

export function GlobeClockRadial() {
  return (
    <div style={{ fontFamily: "'Space Grotesk', sans-serif", background: '#e0e7ff', minHeight: '100%', padding: '48px', display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '40px' }}>
      <div style={{ color: '#6366f1', fontSize: '11px', letterSpacing: '3px', textTransform: 'uppercase' }}>C — Glass Sphere</div>
      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '24px' }}>
        <Mark size={140} />
        <span style={{ fontSize: '52px', fontWeight: 700, letterSpacing: '-2.5px', background: 'linear-gradient(135deg,#3730a3,#6366f1)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>timediverr</span>
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: '12px', background: '#fff', border: '1px solid #c7d2fe', borderRadius: '12px', padding: '12px 20px', boxShadow: '0 4px 20px rgba(99,102,241,0.15)' }}>
        <Mark size={38} />
        <span style={{ fontSize: '22px', fontWeight: 700, color: '#3730a3', letterSpacing: '-1px' }}>timediverr</span>
      </div>
      <div style={{ display: 'flex', alignItems: 'flex-end', gap: '20px' }}>
        {[16,24,32,48].map(n => (
          <div key={n} style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '8px' }}>
            <Mark size={n} />
            <span style={{ fontSize: '10px', color: '#6366f1', letterSpacing: '1px' }}>{n}</span>
          </div>
        ))}
      </div>
    </div>
  )
}
