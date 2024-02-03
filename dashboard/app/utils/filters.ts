import { add, format } from 'date-fns';

interface FilterOptions {
	start?: string;
	end?: string;
	limit?: number;
}

export const generateFilters = (
	url: string,
	opts?: FilterOptions
): Record<string, string | number | undefined> => {
	// Current time period truncated to YYYY-MM-DD
	const currentDate = new Date();
	const endPeriod = format(add(currentDate, { days: 1 }), 'yyyy-MM-dd');
	// Start time period is 24 hours before the current time period
	const startPeriod = format(currentDate, 'yyyy-MM-dd');

	// Convert search params to filters
	const searchParams = new URL(url).searchParams;
	const filters: Record<string, string> = {};
	for (const [key, value] of searchParams) {
		if (value !== null) {
			filters[key] = value;
		}
	}

	return { start: startPeriod, end: endPeriod, limit: opts?.limit, ...filters };
};
