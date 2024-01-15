// Format values into a more readable format
// navigator.languages has a readonly property so we need to copy it into a new array
export const languages =
	typeof document === 'undefined' ? ['en-US'] : [...navigator.languages];

export const countFormatter = Intl.NumberFormat(languages, {
	notation: 'compact',
	maximumFractionDigits: 2,
});

export const percentFormatter = Intl.NumberFormat(languages, {
	style: 'percent',
	maximumFractionDigits: 1,
});

// Convert a duration in milliseconds to a human readable format
// such as 2h1m30s, 1m30s or 30s
export const formatDuration = (durationMs: number) => {
	const totalSeconds = Math.floor(durationMs / 1000);
	const hours = Math.floor(totalSeconds / 3600);
	const minutes = Math.floor((totalSeconds % 3600) / 60);
	const seconds = totalSeconds % 60;

	if (hours > 0) {
		return `${hours}h${minutes}m${seconds}s`;
	}

	return minutes === 0 ? `${seconds}s` : `${minutes}m${seconds}s`;
};
