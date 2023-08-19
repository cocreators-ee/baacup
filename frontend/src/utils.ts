const numberFormatter = new Intl.NumberFormat()
const dateTimeformatter = new Intl.DateTimeFormat()

export function formatNumber(num: number): string {
  return numberFormatter.format(num)
}

export function formatDateTime(dt: string): string {
  if (dt === "0001-01-01T00:00:00Z") {
    return "Never"
  }

  const d = Date.parse(dt)

  return dateTimeformatter.format(d) + " " + new Date(d).toLocaleTimeString()
}
