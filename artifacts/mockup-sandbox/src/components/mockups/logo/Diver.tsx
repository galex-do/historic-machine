export function Diver() {
  return (
    <div style={{
      minHeight: '100vh',
      background: 'linear-gradient(160deg, #eef2ff 0%, #e0e7ff 100%)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      fontFamily: "'Space Grotesk', sans-serif",
      flexDirection: 'column',
      gap: '48px',
      padding: '32px',
    }}>

      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '28px' }}>
        {/* Icon: figure diving into a timeline ring */}
        <svg width="120" height="120" viewBox="0 0 120 120" fill="none">
          <defs>
            <linearGradient id="ring-grad" x1="0" y1="0" x2="120" y2="120" gradientUnits="userSpaceOnUse">
              <stop offset="0%" stopColor="#4f46e5" />
              <stop offset="100%" stopColor="#3730a3" />
            </linearGradient>
            <linearGradient id="diver-grad" x1="40" y1="10" x2="80" y2="60" gradientUnits="userSpaceOnUse">
              <stop offset="0%" stopColor="#4338ca" />
              <stop offset="100%" stopColor="#3730a3" />
            </linearGradient>
          </defs>

          {/* Outer ring — timeline circle */}
          <circle cx="60" cy="68" r="42" stroke="url(#ring-grad)" strokeWidth="5" fill="none" strokeDasharray="8 5" />

          {/* Tick marks on ring */}
          {[0, 60, 120, 180, 240, 300].map((deg, i) => {
            const rad = (deg * Math.PI) / 180
            const x1 = 60 + 38 * Math.cos(rad)
            const y1 = 68 + 38 * Math.sin(rad)
            const x2 = 60 + 45 * Math.cos(rad)
            const y2 = 68 + 45 * Math.sin(rad)
            return <line key={i} x1={x1} y1={y1} x2={x2} y2={y2} stroke="#4f46e5" strokeWidth="2.5" strokeLinecap="round" />
          })}

          {/* Diver body — a stylized figure diving downward */}
          {/* Head */}
          <circle cx="60" cy="20" r="9" fill="url(#diver-grad)" />

          {/* Body */}
          <line x1="60" y1="29" x2="60" y2="50" stroke="url(#diver-grad)" strokeWidth="5" strokeLinecap="round" />

          {/* Arms outstretched like diving */}
          <path d="M60 35 L38 28" stroke="url(#diver-grad)" strokeWidth="4" strokeLinecap="round" />
          <path d="M60 35 L82 28" stroke="url(#diver-grad)" strokeWidth="4" strokeLinecap="round" />

          {/* Legs split like a dive */}
          <path d="M60 50 L48 66" stroke="url(#diver-grad)" strokeWidth="4" strokeLinecap="round" />
          <path d="M60 50 L72 66" stroke="url(#diver-grad)" strokeWidth="4" strokeLinecap="round" />

          {/* Motion trail dots */}
          <circle cx="60" cy="76" r="3" fill="#4f46e5" opacity="0.7" />
          <circle cx="60" cy="87" r="2" fill="#4f46e5" opacity="0.4" />
          <circle cx="60" cy="96" r="1.2" fill="#4f46e5" opacity="0.2" />

          {/* Small map pin where diver enters */}
          <circle cx="60" cy="26" r="0" fill="none" />
        </svg>

        {/* Wordmark */}
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '2px' }}>
          <div style={{ display: 'flex', alignItems: 'baseline', gap: 0 }}>
            <span style={{
              fontSize: '52px',
              fontWeight: 700,
              letterSpacing: '-2px',
              color: '#1e1b4b',
              lineHeight: 1,
            }}>time</span>
            <span style={{
              fontSize: '52px',
              fontWeight: 700,
              letterSpacing: '-2px',
              background: 'linear-gradient(135deg, #4f46e5, #4338ca)',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
              lineHeight: 1,
            }}>diverr</span>
          </div>
          <span style={{
            fontSize: '12px',
            fontWeight: 500,
            letterSpacing: '5px',
            color: '#6366f1',
            textTransform: 'uppercase',
          }}>dive into history</span>
        </div>
      </div>

      {/* Horizontal lockup */}
      <div style={{
        display: 'flex',
        alignItems: 'center',
        gap: '14px',
        background: '#1e1b4b',
        borderRadius: '14px',
        padding: '14px 26px',
      }}>
        <svg width="32" height="32" viewBox="0 0 120 120" fill="none">
          <circle cx="60" cy="68" r="42" stroke="#4f46e5" strokeWidth="6" fill="none" strokeDasharray="10 6" />
          <circle cx="60" cy="20" r="9" fill="#fff" />
          <line x1="60" y1="29" x2="60" y2="50" stroke="#fff" strokeWidth="6" strokeLinecap="round" />
          <path d="M60 35 L38 28" stroke="#fff" strokeWidth="5" strokeLinecap="round" />
          <path d="M60 35 L82 28" stroke="#fff" strokeWidth="5" strokeLinecap="round" />
          <path d="M60 50 L48 66" stroke="#fff" strokeWidth="5" strokeLinecap="round" />
          <path d="M60 50 L72 66" stroke="#fff" strokeWidth="5" strokeLinecap="round" />
        </svg>
        <div style={{ display: 'flex', alignItems: 'baseline', gap: 0 }}>
          <span style={{ fontSize: '24px', fontWeight: 700, letterSpacing: '-0.5px', color: '#c7d2fe', lineHeight: 1 }}>time</span>
          <span style={{ fontSize: '24px', fontWeight: 700, letterSpacing: '-0.5px', color: '#fff', lineHeight: 1 }}>diverr</span>
        </div>
      </div>
    </div>
  )
}
