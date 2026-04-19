import { useSearchParams } from '@remix-run/react';
import { useCallback } from 'react';

import type { ChartStat, ChartType } from '@/components/stats/types';

const useChartType = () => {
	const [searchParams, setSearchParams] = useSearchParams();

	const setChartStat = useCallback(
		(stat: ChartStat) => {
			setSearchParams(
				(params) => {
					params.set('chart[stat]', stat);
					return params;
				},
				{
					preventScrollReset: true,
				},
			);
		},
		[setSearchParams],
	);

	const setChartType = useCallback(
		(type: ChartType) => {
			setSearchParams(
				(params) => {
					params.set('chart[type]', type);
					return params;
				},
				{
					preventScrollReset: true,
				},
			);
		},
		[setSearchParams],
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
