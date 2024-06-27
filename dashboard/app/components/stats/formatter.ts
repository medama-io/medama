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
	if (durationMs === 0) {
		return '0s';
	}

	const totalSeconds = Math.floor(durationMs / 1000);
	const hours = Math.floor(totalSeconds / 3600);
	const minutes = Math.floor((totalSeconds % 3600) / 60);
	const seconds = totalSeconds % 60;
	const milliseconds = durationMs % 1000;

	if (hours > 0) {
		return `${hours}h${minutes}m${seconds}s`;
	}

	if (minutes === 0) {
		if (milliseconds < 1000) {
			return `0.${Math.floor(milliseconds / 100)}s`;
		}

		return `${seconds}s`;
	}

	return `${minutes}m${seconds}s`;
};

// Format percentage value
export const formatPercentage: Formatter = (value = 0) => {
	return percentFormatter(value);
};

// Format count value
export const formatCount: Formatter = (value = 0) => {
	return countFormatter(value);
};
