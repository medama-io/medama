import {
	endOfDay,
	endOfHour,
	formatRFC3339,
	startOfDay,
	startOfHour,
	startOfMonth,
	startOfWeek,
	sub,
} from 'date-fns';

interface FilterOptions {
	start?: string;
	end?: string;
	limit?: number;
}

const generatePeriods = (period: string) => {
	const currentDate = new Date();
	let startPeriod: Date;
	let endPeriod: Date;
	let interval = 'hour';

	switch (period) {
		case 'today': {
			startPeriod = startOfDay(currentDate);
			endPeriod = endOfDay(currentDate);
			break;
		}
		case 'yesterday': {
			startPeriod = startOfDay(sub(currentDate, { days: 1 }));
			endPeriod = endOfDay(sub(currentDate, { days: 1 }));
			break;
		}
		case 'quarter': {
			startPeriod = startOfWeek(sub(currentDate, { months: 3 }));
			endPeriod = endOfDay(currentDate);
			interval = 'week';
			break;
		}
		case 'half': {
			startPeriod = startOfMonth(sub(currentDate, { months: 6 }));
			endPeriod = endOfDay(currentDate);
			interval = 'month';
			break;
		}
		case 'year': {
			startPeriod = startOfMonth(sub(currentDate, { years: 1 }));
			endPeriod = endOfDay(currentDate);
			interval = 'month';
			break;
		}
		case 'all': {
			startPeriod = new Date(2024, 0);
			endPeriod = endOfDay(currentDate);
			interval = 'month';
			break;
		}
		default: {
			// Manually parse periods like 24h, 14d, 30d, etc
			if (period.endsWith('d')) {
				const days = Number.parseInt(period, 10);
				startPeriod = startOfDay(sub(currentDate, { days }));
				endPeriod = endOfHour(currentDate);
				interval = 'day';
			} else if (period.endsWith('h')) {
				const hours = Number.parseInt(period, 10);
				startPeriod = startOfHour(sub(currentDate, { hours }));
				endPeriod = endOfHour(currentDate);
			} else {
				throw new Error(`Invalid period: ${period}`);
			}
		}
	}

	return {
		start: formatRFC3339(startPeriod),
		end: formatRFC3339(endPeriod),
		defaultInterval: interval,
	};
};

export const generateFilters = (
	searchParams: URLSearchParams,
	opts?: FilterOptions,
): [Record<string, string | number | undefined>, string | undefined] => {
	// Convert period param to start and end
	const period = searchParams.get('period');
	const { start, end, defaultInterval } = generatePeriods(period ?? 'today');

	// Get interval param
	const interval = searchParams.get('interval') ?? defaultInterval;

	const filters: Record<string, string> = {};
	for (const [key, value] of searchParams) {
		if (value !== null && key !== 'period') {
			filters[key] = value;
		}
	}

	return [{ start, end, limit: opts?.limit, ...filters }, interval];
};
