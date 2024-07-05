import type { DataRow } from './types';

const sortBy =
	// biome-ignore lint/suspicious/noExplicitAny: Generic function.
		<T extends Record<string, any>>(key: keyof T) =>
		(a: T, b: T) =>
			a[key] > b[key] ? 1 : b[key] > a[key] ? -1 : 0;

const aggregateStatsByLanguage = (stats: DataRow[]): DataRow[] => {
	// A record of aggregated stats for each language. e.g. "English", DataRow
	const languageStats: Record<string, DataRow> = {};

	for (const stat of stats) {
		const language = stat.language || 'Unknown';

		if (!languageStats[language]) {
			languageStats[language] = { ...stat };
		} else {
			// Add the stats together
			languageStats[language] = {
				language,
				// Calculate sums
				visitors:
					(languageStats[language].visitors || 0) + (stat.visitors || 0),
				pageviews:
					(languageStats[language].pageviews || 0) + (stat.pageviews || 0),
				bounces: (languageStats[language].bounces || 0) + (stat.bounces || 0),
				duration:
					(languageStats[language].duration || 0) + (stat.duration || 0),
			};
		}
	}

	return Object.values(languageStats);
};

export { sortBy, aggregateStatsByLanguage };
