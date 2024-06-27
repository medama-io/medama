import { useSearchParams } from '@remix-run/react';
import { useCallback } from 'react';

import type { ChartStat, ChartType } from '@/components/stats/types';

const useChartType = () => {
	const [searchParams, setSearchParams] = useSearchParams();

	const setChartStat = useCallback(
		(stat: ChartStat) => {
			searchParams.set('chart[stat]', stat);
			setSearchParams(searchParams, {
				preventScrollReset: true,
			});
		},
		[searchParams, setSearchParams],
	);

	const setChartType = useCallback(
		(type: ChartType) => {
			searchParams.set('chart[type]', type);
			setSearchParams(searchParams, {
				preventScrollReset: true,
			});
		},
		[searchParams, setSearchParams],
	);

	const getChartStat = useCallback(() => {
		const value = searchParams.get('chart[stat]') as ChartStat;
		return value || 'visitors';
	}, [searchParams]);

	const getChartType = useCallback(() => {
		const value = searchParams.get('chart[type]') as ChartType;
		return value || 'area';
	}, [searchParams]);

	return {
		searchParams,
		setChartStat,
		setChartType,
		getChartStat,
		getChartType,
	};
};

export { useChartType };
