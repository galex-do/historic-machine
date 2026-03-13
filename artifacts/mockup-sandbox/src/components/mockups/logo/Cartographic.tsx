export function Cartographic() {
  return (
    <div style={{
      minHeight: '100vh',
      background: '#f8f7ff',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      fontFamily: "'Space Grotesk', sans-serif",
      flexDirection: 'column',
      gap: '48px',
      padding: '32px'
    }}>
      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '32px' }}>

        {/* Mark: Globe with meridians */}
        <svg width="120" height="120" viewBox="0 0 120 120" fill="none">
          <defs>
            <linearGradient id="globe-grad" x1="0" y1="0" x2="120" y2="120" gradientUnits="userSpaceOnUse">
              <stop offset="0%" stopColor="#4338ca" />
              <stop offset="100%" stopColor="#3730a3" />
            </linearGradient>
            <clipPath id="globe-clip">
              <circle cx="60" cy="60" r="48" />
            </clipPath>
          </defs>

          {/* Globe base */}
          <circle cx="60" cy="60" r="48" fill="url(#globe-grad)" />

          {/* Grid lines clipped to globe */}
          <g clipPath="url(#globe-clip)" stroke="rgba(255,255,255,0.2)" strokeWidth="1.2" fill="none">
            {/* Latitude lines */}
            <ellipse cx="60" cy="60" rx="48" ry="16" />
            <ellipse cx="60" cy="60" rx="48" ry="32" />
            <line x1="12" y1="60" x2="108" y2="60" />
            <line x1="12" y1="36" x2="108" y2="36" />
            <line x1="12" y1="84" x2="108" y2="84" />
            {/* Longitude lines */}
            <line x1="60" y1="12" x2="60" y2="108" />
            <ellipse cx="60" cy="60" rx="20" ry="48" />
            <ellipse cx="60" cy="60" rx="38" ry="48" />
          </g>

          {/* Globe border */}
          <circle cx="60" cy="60" r="48" stroke="rgba(255,255,255,0.15)" strokeWidth="1.5" fill="none" />

          {/* Pin marker */}
          <circle cx="60" cy="44" r="8" fill="#fff" opacity="0.95" />
          <circle cx="60" cy="44" r="4" fill="#4338ca" />
          <line x1="60" y1="52" x2="60" y2="62" stroke="#fff" strokeWidth="2" opacity="0.95" />
        </svg>

        {/* Wordmark */}
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '4px' }}>
          <span style={{
            fontSize: '52px',
            fontWeight: 700,
            letterSpacing: '-2px',
            background: 'linear-gradient(135deg, #4338ca, #3730a3)',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
            lineHeight: 1,
          }}>timediverr</span>
          <span style={{
            fontSize: '13px',
            fontWeight: 400,
            letterSpacing: '4px',
            color: '#6b7280',
            textTransform: 'uppercase',
          }}>historical map</span>
        </div>
      </div>

      {/* Horizontal lockup */}
      <div style={{
        display: 'flex',
        alignItems: 'center',
        gap: '16px',
        background: 'linear-gradient(135deg, #4338ca, #3730a3)',
        borderRadius: '16px',
        padding: '16px 28px',
      }}>
        <svg width="36" height="36" viewBox="0 0 120 120" fill="none">
          <defs>
            <clipPath id="gc2">
              <circle cx="60" cy="60" r="48" />
            </clipPath>
          </defs>
          <circle cx="60" cy="60" r="48" fill="rgba(255,255,255,0.15)" />
          <g clipPath="url(#gc2)" stroke="rgba(255,255,255,0.35)" strokeWidth="2" fill="none">
            <ellipse cx="60" cy="60" rx="48" ry="18" />
            <line x1="12" y1="60" x2="108" y2="60" />
            <line x1="60" y1="12" x2="60" y2="108" />
            <ellipse cx="60" cy="60" rx="22" ry="48" />
          </g>
          <circle cx="60" cy="60" r="48" stroke="rgba(255,255,255,0.25)" strokeWidth="2" fill="none" />
          <circle cx="60" cy="44" r="9" fill="#fff" opacity="0.9" />
          <circle cx="60" cy="44" r="4.5" fill="#4338ca" />
          <line x1="60" y1="53" x2="60" y2="64" stroke="#fff" strokeWidth="2.5" opacity="0.9" />
        </svg>
        <span style={{
          fontSize: '26px',
          fontWeight: 700,
          letterSpacing: '-1px',
          color: '#fff',
          lineHeight: 1,
        }}>timediverr</span>
      </div>
    </div>
  )
}
