import {
	type ClientLoaderFunctionArgs,
	type MetaFunction,
	Outlet,
	type ShouldRevalidateFunctionArgs,
	json,
	useLoaderData,
} from '@remix-run/react';
import { useMemo } from 'react';

import { userLoggedIn } from '@/api/user';
import { websiteList } from '@/api/websites';
import { Chart } from '@/components/stats/Chart';
import { Filters } from '@/components/stats/Filter';
import { StatsHeader } from '@/components/stats/StatsHeader';
import type { StatHeaderData } from '@/components/stats/types';
import { useChartType } from '@/hooks/use-chart-type';
import { fetchStats } from '@/utils/stats';

export const meta: MetaFunction = () => {
	return [{ title: 'Dashboard | Medama' }];
};

export const clientLoader = async ({
	request,
	params,
}: ClientLoaderFunctionArgs) => {
	await userLoggedIn();

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

	return json({ stats, websites: websites.data });
};

const LABEL_MAP = {
	visitors: 'Visitors',
	pageviews: 'Page Views',
	duration: 'Time Spent',
	bounces: 'Bounce Rate',
};

export default function Index() {
	const { stats, websites } = useLoaderData<typeof clientLoader>();
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
				const valueMap = {
					visitors: item.visitors ?? 0,
					pageviews: item.pageviews ?? 0,
					bounces: item.bounce_percentage ?? 0,
					duration: item.duration ?? 0,
				};

				return {
					date: item.date,
					value: valueMap[chart as keyof typeof valueMap],
				};
			}) ?? [],
		[summary.interval, chart],
	);

	const websiteList = websites.map((website) => website.hostname);

	const label = LABEL_MAP[chart as keyof typeof LABEL_MAP];

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
