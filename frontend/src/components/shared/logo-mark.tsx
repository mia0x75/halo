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
        d="M 265.488 300.839 A 24.000 24.000 0 0 1 289.228 280.362 L 425.839 280.362 A 24.000 24.000 0 0 1 449.577 307.901 L 425.045 472.440 A 24.000 24.000 0 0 0 448.783 499.979 L 530.287 499.979 A 24.000 24.000 0 0 0 554.028 479.498 L 599.836 170.473 A 24.000 24.000 0 0 1 623.577 149.992 L 760.030 149.992 A 24.000 24.000 0 0 1 783.770 177.514 L 675.158 909.581 A 24.000 24.000 0 0 1 651.418 930.059 L 515.044 930.059 A 24.000 24.000 0 0 1 491.313 902.476 L 535.836 607.583 A 24.000 24.000 0 0 0 512.105 580.000 L 429.766 580.000 A 24.000 24.000 0 0 0 406.026 600.477 L 379.458 779.523 A 24.000 24.000 0 0 1 355.718 800.000 L 219.243 800.000 A 24.000 24.000 0 0 1 195.503 772.477 Z"
        fill="var(--halo-logo-glyph)"
      />
    </svg>
  )
}
