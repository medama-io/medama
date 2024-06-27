import {
	Outlet,
	json,
	useLoaderData,
	useSearchParams,
	type ClientLoaderFunctionArgs,
	type MetaFunction,
} from '@remix-run/react';
import { useMemo } from 'react';

import { userLoggedIn } from '@/api/user';
import { Chart } from '@/components/stats/Chart';
import { Filters } from '@/components/stats/Filter';
import { StatsHeader } from '@/components/stats/StatsHeader';
import type { StatHeaderData } from '@/components/stats/types';
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
	const chart = searchParams.get('chart');

	const stats = await fetchStats(request, params, {
		dataset: ['summary'],
		chartStat: chart || 'visitors',
	});

	return json(stats);
};

const LABEL_MAP = {
	visitors: 'Visitors',
	pageviews: 'Page Views',
	duration: 'Time Spent',
	bounces: 'Bounce Rate',
};

export default function Index() {
	const { summary } = useLoaderData<typeof clientLoader>();
	if (!summary) throw new Error('Summary data is required');

	const [searchParams] = useSearchParams();
	const chart = searchParams.get('chart') || 'visitors';

	const { current, previous } = summary;

	const stats: StatHeaderData[] = [
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
			current: current.bounces / current.visitors || 0,
			previous: previous
				? previous.bounces / previous.visitors || 0
				: undefined,
		},
	];

	const chartData = useMemo(
		() =>
			summary.interval?.map((item) => {
				const valueMap = {
					visitors: item.visitors ?? 0,
					pageviews: item.pageviews ?? 0,
					bounces: (item.bounces ?? 0) / (item.visitors ?? 0) || 0,
					duration: item.duration ?? 0,
				};

				return {
					date: item.date,
					value: valueMap[chart as keyof typeof valueMap],
				};
			}) ?? [],
		[summary.interval, chart],
	);

	const label = LABEL_MAP[chart as keyof typeof LABEL_MAP];

	return (
		<>
			<StatsHeader stats={stats} chart={chart} />
			<main>
				<Filters />
				{summary.interval && (
					<>
						<Chart type="bar" label={label} data={chartData} />
						<Chart type="area" label={label} data={chartData} />
					</>
				)}
				<Outlet />
			</main>
		</>
	);
}
