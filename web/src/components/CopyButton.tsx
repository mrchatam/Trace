type Props = {
  text: string
  label?: string
}

export function CopyButton({ text, label = 'Copy' }: Props) {
  async function onCopy() {
    try {
      await navigator.clipboard.writeText(text)
    } catch {
      // Fallback for older browsers / denied permission
      const ta = document.createElement('textarea')
      ta.value = text
      document.body.appendChild(ta)
      ta.select()
      document.execCommand('copy')
      document.body.removeChild(ta)
    }
  }

  return (
    <button type="button" className="btn" onClick={() => void onCopy()} aria-label={label}>
      {label}
    </button>
  )
}
