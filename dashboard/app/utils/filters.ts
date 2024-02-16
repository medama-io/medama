import {
	endOfToday,
	endOfYesterday,
	formatRFC3339,
	startOfToday,
	startOfYesterday,
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

	switch (period) {
		case 'today': {
			startPeriod = startOfToday();
			endPeriod = endOfToday();
			break;
		}
		case 'yesterday': {
			startPeriod = startOfYesterday();
			endPeriod = endOfYesterday();
			break;
		}
		case 'quarter': {
			startPeriod = sub(currentDate, { months: 3 });
			endPeriod = currentDate;
			break;
		}
		case 'half': {
			startPeriod = sub(currentDate, { months: 6 });
			endPeriod = currentDate;
			break;
		}
		case 'year': {
			startPeriod = sub(currentDate, { years: 1 });
			endPeriod = currentDate;
			break;
		}
		case 'all': {
			startPeriod = new Date(0);
			endPeriod = currentDate;
			break;
		}
		default: {
			// Manually parse periods like 24h, 14d, 30d, etc
			if (period.endsWith('d')) {
				const days = Number.parseInt(period, 10);
				startPeriod = sub(currentDate, { days });
				endPeriod = currentDate;
			} else if (period.endsWith('h')) {
				const hours = Number.parseInt(period, 10);
				startPeriod = sub(currentDate, { hours });
				endPeriod = currentDate;
			} else {
				throw new Error(`Invalid period: ${period}`);
			}
		}
	}

	return {
		start: formatRFC3339(startPeriod),
		end: formatRFC3339(endPeriod),
	};
};

export const generateFilters = (
	url: string,
	opts?: FilterOptions
): Record<string, string | number | undefined> => {
	// Convert search params to filters
	const searchParams = new URL(url).searchParams;

	// Convert period param to start and end
	const period = searchParams.get('period');
	const { start, end } = generatePeriods(period ?? 'today');

	const filters: Record<string, string> = {};
	for (const [key, value] of searchParams) {
		if (value !== null && key !== 'period') {
			filters[key] = value;
		}
	}

	return { start, end, limit: opts?.limit, ...filters };
};
