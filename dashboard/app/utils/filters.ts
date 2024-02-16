import { add, format, sub } from 'date-fns';

interface FilterOptions {
	start?: string;
	end?: string;
	limit?: number;
}

const generatePeriods = (period: string) => {
	// TODO: More granular periods
	const currentDate = new Date();
	let startPeriod: Date;
	let endPeriod: Date;

	switch (period) {
		case 'today': {
			startPeriod = currentDate;
			endPeriod = add(currentDate, { days: 1 });
			break;
		}
		case 'yesterday': {
			startPeriod = sub(currentDate, { days: 1 });
			endPeriod = currentDate;
			break;
		}
		case '24h': {
			startPeriod = sub(currentDate, { hours: 24 });
			endPeriod = currentDate;
			break;
		}
		case '72h': {
			startPeriod = sub(currentDate, { hours: 72 });
			endPeriod = currentDate;
			break;
		}
		case '7d': {
			startPeriod = sub(currentDate, { days: 7 });
			endPeriod = currentDate;
			break;
		}
		case '14d': {
			startPeriod = sub(currentDate, { days: 14 });
			endPeriod = currentDate;
			break;
		}
		case '30d': {
			startPeriod = sub(currentDate, { days: 30 });
			endPeriod = currentDate;
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
			throw new Error('Invalid period');
		}
	}

	return {
		start: format(startPeriod, 'yyyy-MM-dd'),
		end: format(endPeriod, 'yyyy-MM-dd'),
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
