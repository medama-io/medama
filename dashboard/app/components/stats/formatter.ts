// Define types for formatters
type Formatter = (value: number) => string;

// Determine languages array
const languages: string[] =
	typeof document === 'undefined' ? ['en-US'] : [...navigator.languages];

// Intl formatters
const countFormatter: Formatter = Intl.NumberFormat(languages, {
	notation: 'compact',
	maximumFractionDigits: 2,
}).format;

const percentFormatter: Formatter = Intl.NumberFormat(languages, {
	style: 'percent',
	maximumFractionDigits: 1,
}).format;

// Convert duration in milliseconds to a human readable format
export const formatDuration: Formatter = (durationMs = 0) => {
	let milliseconds = durationMs;
	if (milliseconds === 0 || milliseconds < 50) return 'N/A';

	const hours = Math.floor(milliseconds / 3600000);
	milliseconds %= 3600000;
	const minutes = Math.floor(milliseconds / 60000);
	milliseconds %= 60000;
	const seconds = Math.round(milliseconds / 1000);

	if (hours > 0) return `${hours}h${minutes}m${seconds}s`;
	if (minutes > 0) return `${minutes}m${seconds}s`;
	if (seconds > 0) return `${seconds}s`;
	return `0.${Math.round(milliseconds / 100)}s`;
};

// Format percentage value
export const formatPercentage: Formatter = (value = 0) => {
	return percentFormatter(value);
};

// Format count value
export const formatCount: Formatter = (value = 0) => {
	return countFormatter(value);
};
