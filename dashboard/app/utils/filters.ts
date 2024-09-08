import {
	endOfDay,
	endOfHour,
	endOfMonth,
	endOfWeek,
	formatRFC3339,
	intervalToDuration,
	parseISO,
	startOfDay,
	startOfHour,
	startOfMonth,
	startOfWeek,
	sub,
} from 'date-fns';

import { StatusError } from '@/components/layout/Error';

interface FilterOptions {
	start?: string;
	end?: string;
	limit?: number;
}

const generatePeriods = (searchParams: URLSearchParams) => {
	const period = searchParams.get('period') ?? 'today';
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
			endPeriod = endOfWeek(currentDate);
			interval = 'week';
			break;
		}
		case 'half': {
			startPeriod = startOfMonth(sub(currentDate, { months: 6 }));
			endPeriod = endOfMonth(currentDate);
			interval = 'month';
			break;
		}
		case 'year': {
			startPeriod = startOfMonth(sub(currentDate, { years: 1 }));
			endPeriod = endOfMonth(currentDate);
			interval = 'month';
			break;
		}
		case 'all': {
			startPeriod = new Date(2024, 0);
			endPeriod = endOfDay(currentDate);
			interval = 'month';
			break;
		}
		case 'custom': {
			const start = searchParams.get('start');
			const end = searchParams.get('end');
			if (start && end) {
				startPeriod = parseISO(start);
				endPeriod = parseISO(end);
			} else {
				throw new Error('Invalid custom period');
			}

			const diff = intervalToDuration({ start: startPeriod, end: endPeriod });
			if (diff.months && diff.months > 1) {
				interval = 'month';
			} else if (diff.days && diff.days > 1) {
				interval = 'day';
			} else {
				interval = 'hour';
			}

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
				throw new StatusError(400, `Invalid time period: ${period}`);
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
	const { start, end, defaultInterval } = generatePeriods(searchParams);
	const interval = searchParams.get('interval') ?? defaultInterval;

	const filters: Record<string, string | number | undefined> = {
		start,
		end,
		limit: opts?.limit,
	};

	for (const [key, value] of searchParams) {
		if (value && !['period', 'start', 'end', 'interval'].includes(key)) {
			filters[key] = value;
		}
	}

	return [filters, interval];
};
