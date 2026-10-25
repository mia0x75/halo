export interface LogoMarkProps {
  /** Rendered size in pixels; the source viewBox is square. */
  size?: number
  className?: string
}

/**
 * Brand mark, inlined instead of loaded as an asset: it appears in the
 * sidebar, the sign-in screen and the setup wizard, and inlining keeps it out
 * of the request path so it paints with the first frame.
 *
 * Both fills come from theme variables so each theme picks its own pair. The
 * original near-black disc disappears against a dark background, hence the
 * light-mode and dark-mode values in index.css.
 */
export function LogoMark({ size = 32, className }: LogoMarkProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 1080 1080"
      width={size}
      height={size}
      role="img"
      aria-label="Halo"
      className={className}
    >
      <title>Halo</title>
      <path
        d="M10 540 C10 247.289 247.289 10 540 10 C832.711 10 1070 247.289 1070 540 C1070 832.711 832.711 1070 540 1070 C247.289 1070 10 832.711 10 540 Z"
        fill="var(--halo-logo-disc)"
      />
      <path
        d="M268.527 280.362 L453.683 280.362 L420.939 499.979 L550.992 499.979 L602.872 149.992 L787.853 149.992 L672.12 930.059 L487.149 930.059 L540 580 L409.065 580 L376.419 800 L191.419 800 Z"
        fill="var(--halo-logo-glyph)"
      />
    </svg>
  )
}
