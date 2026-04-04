export function formatDate(date: string | undefined): string {
    if (!date) {
        return 'never'
    }
    const instant = Temporal.Instant.from(date);
    const zdt = instant.toZonedDateTimeISO(Temporal.Now.timeZoneId())
    return zdt.toLocaleString(undefined, { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

export function formatDateAgo(date: string | undefined): string {
    if (!date) {
        return 'never'
    }
    const instant = Temporal.Instant.from(date)
    const duration = Temporal.Now.zonedDateTimeISO().since(instant.toZonedDateTimeISO(Temporal.Now.timeZoneId()), { largestUnit: 'year' })

    if (duration.years > 0) return `${duration.years} year${duration.years > 1 ? 's' : ''} ago`
    if (duration.months > 0) return `${duration.months} month${duration.months > 1 ? 's' : ''} ago`
    if (duration.days > 0) return `${duration.days} day${duration.days > 1 ? 's' : ''} ago`
    if (duration.hours > 0) return `${duration.hours} hour${duration.hours > 1 ? 's' : ''} ago`
    if (duration.minutes > 0) return `${duration.minutes} minute${duration.minutes > 1 ? 's' : ''} ago`
    return 'just now'
}


export function toDateTimeLocal(date: Temporal.Instant) {
    if (!date) {
        return '';
    }
    return date.toZonedDateTimeISO(Temporal.Now.timeZoneId()).toPlainDateTime().toString({ smallestUnit: 'minute' });
}

export function fromDateTimeLocal(date: string) {
    return Temporal.PlainDateTime.from(date).toZonedDateTime(Temporal.Now.timeZoneId()).toInstant();
}