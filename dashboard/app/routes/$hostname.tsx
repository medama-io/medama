import { useMemo } from 'react';
import { Outlet, type ShouldRevalidateFunctionArgs } from 'react-router';

import { websiteList } from '@/api/websites';
import { Chart } from '@/components/stats/Chart';
import { Filters } from '@/components/stats/Filter';
import { StatsHeader } from '@/components/stats/StatsHeader';
import type { ChartStat, StatHeaderData } from '@/components/stats/types';
import { useChartType } from '@/hooks/use-chart-type';
import { fetchStats } from '@/utils/stats';
import type { Route } from './+types/$hostname';

export const meta: Route.MetaFunction = () => {
	return [{ title: 'Dashboard | Medama' }];
};

export const clientLoader = async ({
	request,
	params,
}: Route.ClientLoaderArgs) => {
	// Check chart param for the chart data to display
	const searchParams = new URL(request.url).searchParams;
	const chart = searchParams.get('chart[stat]');

	const [stats, websites] = await Promise.all([
		fetchStats(request, params, {
			dataset: ['summary'],
			chartStat: chart || 'visitors',
		}),
		websiteList(),
	]);

	return { stats, websites: websites.data };
};

const LABEL_MAP: Record<ChartStat, string> = {
	visitors: 'Visitors',
	pageviews: 'Page Views',
	duration: 'Time Spent',
	bounces: 'Bounce Rate',
};

export default function Index({ loaderData }: Route.ComponentProps) {
	const { stats, websites } = loaderData;
	const { summary } = stats;
	if (!websites) throw new Error('Websites data is required');
	if (!summary) throw new Error('Summary data is required');
	const { current, previous } = summary;

	const { getChartStat, getChartType } = useChartType();
	const chart = getChartStat();
	const type = getChartType();

	const data: StatHeaderData[] = [
		{
			label: 'Visitors',
			chart: 'visitors',
			current: current.visitors,
			previous: previous?.visitors,
		},
		{
			label: 'Page Views',
			chart: 'pageviews',
			current: current.pageviews,
			previous: previous?.pageviews,
		},
		{
			label: 'Time Spent',
			chart: 'duration',
			current: current.duration,
			previous: previous?.duration,
		},
		{
			label: 'Bounce Rate',
			chart: 'bounces',
			current: current.bounce_percentage,
			previous: previous?.bounce_percentage,
		},
	];

	const chartData = useMemo(
		() =>
			summary.interval?.map((item) => {
				const valueMap: Record<ChartStat, number> = {
					visitors: item.visitors ?? 0,
					pageviews: item.pageviews ?? 0,
					bounces: item.bounce_percentage ?? 0,
					duration: item.duration ?? 0,
				};

				return {
					date: item.date,
					value: valueMap[chart],
				};
			}) ?? [],
		[summary.interval, chart],
	);

	const websiteList = websites.map((website) => website.hostname);

	const label = LABEL_MAP[chart];

	return (
		<>
			<StatsHeader stats={data} chart={chart} websites={websiteList} />
			<main>
				<Filters />
				{summary.interval && (
					<Chart type={type} label={label} data={chartData} />
				)}
				<Outlet />
			</main>
		</>
	);
}

export const shouldRevalidate = ({
	currentUrl,
	nextUrl,
	defaultShouldRevalidate,
}: ShouldRevalidateFunctionArgs) => {
	const currentParams = new URL(currentUrl).searchParams;
	const nextParams = new URL(nextUrl).searchParams;
	if (
		// We don't want to revalidate if the chart type or stat changes as
		// the data doesn't change, only the presentation.
		currentParams.get('chart[type]') !== nextParams.get('chart[type]') ||
		currentParams.get('chart[stat]') !== nextParams.get('chart[stat]')
	) {
		return false;
	}

	return defaultShouldRevalidate;
};
