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

      {/* GLOBE — continent silhouettes (orthographic, center 30°W 0°N)
           All coordinates computed via:
             nx = cos(lat)*sin(lon - lon0),  ny = -sin(lat)
           where lon0 = -30°.  GlobeClip handles circular boundary. */}
      <g clipPath={`url(#g${size})`} fill={fg} fillOpacity="0.88" stroke="none">
        {(() => {
          const q = (nx: number, ny: number) => `${C + nx * R} ${C + ny * R}`
          const poly = (pts: [number,number][]) =>
            `M${pts.map(([x,y]) => q(x,y)).join('L')}Z`
          return <>
            {/* Greenland — western lobe visible in globe sector */}
            <path d={poly([
              [-0.063,-0.970], [-0.021,-0.970], [ 0.027,-0.951],
              [ 0.068,-0.921], [ 0.034,-0.875], [-0.121,-0.866],
              [-0.152,-0.914], [-0.106,-0.951],
            ])} />
            {/* North America — clockwise
                NW-Canada → Arctic coast → Hudson Bay → Labrador → Nova Scotia
                → Cape Hatteras → Florida → Gulf coast → Yucatan → Baja tip
                → California → Vancouver → Alaska → back */}
            <path d={poly([
              [-0.358,-0.934], [-0.265,-0.956], [-0.324,-0.906],
              [-0.266,-0.731], [-0.354,-0.707], [-0.433,-0.695],
              [-0.589,-0.574], [-0.704,-0.423], [-0.718,-0.500],
              [-0.757,-0.485], [-0.827,-0.438], [-0.824,-0.309],
              [-0.871,-0.276], [-0.907,-0.391],
              [-0.799,-0.602], [-0.655,-0.755], [-0.528,-0.839],
            ])} />
            {/* South America — clockwise
                Colombia-W → Colombia-N → Venezuela → Guyana → Brazil-NE (bulge)
                → Salvador → Rio → Buenos-Aires → Cape Horn → Patagonia-W
                → Chile → Ecuador/Peru */}
            <path d={poly([
              [-0.743,-0.017], [-0.724,-0.139], [-0.536,-0.174],
              [-0.374,-0.070], [-0.086, 0.139], [-0.136, 0.225],
              [-0.207, 0.391], [-0.310, 0.500], [-0.389, 0.559],
              [-0.344, 0.829], [-0.435, 0.788], [-0.470, 0.682],
              [-0.766, 0.035],
            ])} />
            {/* Africa — clock-sector clips northern half; southern tip visible.
                Clockwise: Senegal → Libya/Egypt → E-Africa → Cape → Angola → back */}
            <path d={poly([
              [ 0.218,-0.242], [ 0.578,-0.530], [ 0.759,-0.375],
              [ 0.828, 0.407], [ 0.718, 0.500], [ 0.628, 0.574],
              [ 0.616, 0.559], [ 0.644, 0.375], [ 0.646, 0.259],
              [ 0.668, 0.035], [ 0.543,-0.087], [ 0.272,-0.174],
            ])} />
          </>
        })()}
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
