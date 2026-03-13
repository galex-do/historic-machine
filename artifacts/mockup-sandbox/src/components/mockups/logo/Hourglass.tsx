export function Hourglass() {
  return (
    <div style={{
      minHeight: '100vh',
      background: '#0f0e1a',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      fontFamily: "'Space Grotesk', sans-serif",
      flexDirection: 'column',
      gap: '48px',
      padding: '32px'
    }}>

      {/* Stacked icon + wordmark */}
      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '28px' }}>
        <svg width="100" height="120" viewBox="0 0 100 120" fill="none">
          <defs>
            <linearGradient id="hg-top" x1="0" y1="0" x2="100" y2="60" gradientUnits="userSpaceOnUse">
              <stop offset="0%" stopColor="#6d63f5" />
              <stop offset="100%" stopColor="#4338ca" />
            </linearGradient>
            <linearGradient id="hg-bot" x1="0" y1="60" x2="100" y2="120" gradientUnits="userSpaceOnUse">
              <stop offset="0%" stopColor="#3730a3" />
              <stop offset="100%" stopColor="#4f46e5" />
            </linearGradient>
            <linearGradient id="hg-frame" x1="0" y1="0" x2="0" y2="120" gradientUnits="userSpaceOnUse">
              <stop offset="0%" stopColor="#7c73f8" />
              <stop offset="100%" stopColor="#4338ca" />
            </linearGradient>
          </defs>

          {/* Hourglass body top half */}
          <path d="M8 8 L92 8 L54 58 L46 58 Z" fill="url(#hg-top)" opacity="0.9" />
          {/* Hourglass body bottom half */}
          <path d="M46 62 L54 62 L92 112 L8 112 Z" fill="url(#hg-bot)" opacity="0.9" />

          {/* Frame lines */}
          <rect x="4" y="4" width="92" height="8" rx="4" fill="url(#hg-frame)" />
          <rect x="4" y="108" width="92" height="8" rx="4" fill="url(#hg-frame)" />

          {/* Dot grid in top half (sand falling) */}
          {[
            [22, 22], [38, 22], [54, 22], [70, 22],
            [28, 34], [44, 34], [60, 34],
            [36, 46], [52, 46],
          ].map(([cx, cy], i) => (
            <circle key={i} cx={cx} cy={cy} r="3.5" fill="rgba(255,255,255,0.5)" />
          ))}

          {/* Pinch — map pin */}
          <circle cx="50" cy="60" r="5.5" fill="#fff" />
          <circle cx="50" cy="60" r="2.5" fill="#4338ca" />

          {/* Dots falling in bottom half */}
          {[
            [44, 74], [56, 74],
            [36, 86], [50, 86], [64, 86],
            [28, 98], [42, 98], [58, 98], [72, 98],
          ].map(([cx, cy], i) => (
            <circle key={i} cx={cx} cy={cy} r="3.5" fill="rgba(255,255,255,0.3)" />
          ))}
        </svg>

        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '4px' }}>
          <span style={{
            fontSize: '50px',
            fontWeight: 700,
            letterSpacing: '-2px',
            color: '#fff',
            lineHeight: 1,
          }}>timediverr</span>
          <span style={{
            fontSize: '12px',
            fontWeight: 400,
            letterSpacing: '5px',
            color: '#6d63f5',
            textTransform: 'uppercase',
          }}>explore history</span>
        </div>
      </div>

      {/* Dark pill lockup */}
      <div style={{
        display: 'flex',
        alignItems: 'center',
        gap: '14px',
        border: '1px solid rgba(109, 99, 245, 0.4)',
        borderRadius: '14px',
        padding: '14px 24px',
        background: 'rgba(79,70,229,0.1)',
        backdropFilter: 'blur(12px)',
      }}>
        <svg width="28" height="34" viewBox="0 0 100 120" fill="none">
          <path d="M8 8 L92 8 L54 58 L46 58 Z" fill="#6d63f5" opacity="0.85" />
          <path d="M46 62 L54 62 L92 112 L8 112 Z" fill="#4338ca" opacity="0.85" />
          <rect x="4" y="4" width="92" height="8" rx="4" fill="#7c73f8" />
          <rect x="4" y="108" width="92" height="8" rx="4" fill="#4338ca" />
          <circle cx="50" cy="60" r="6" fill="#fff" />
          <circle cx="50" cy="60" r="2.8" fill="#4338ca" />
        </svg>
        <span style={{
          fontSize: '24px',
          fontWeight: 700,
          letterSpacing: '-0.5px',
          color: '#fff',
        }}>timediverr</span>
      </div>
    </div>
  )
}
