import {
	json,
	type MetaFunction,
	type ClientLoaderFunctionArgs,
	Outlet,
	useLoaderData,
} from '@remix-run/react';
import invariant from 'tiny-invariant';

import { StackedBarChart } from '@/components/stats/Chart';
import { Filters } from '@/components/stats/Filter';
import { StatsHeader } from '@/components/stats/StatsHeader';
import { fetchStats } from '@/utils/stats';

export const meta: MetaFunction = () => {
	return [{ title: 'Dashboard | Medama' }];
};

export const clientLoader = async ({
	request,
	params,
}: ClientLoaderFunctionArgs) => {
	const stats = await fetchStats(request, params, {
		dataset: ['summary'],
	});

	return json(stats);
};

export default function Index() {
	const { summary } = useLoaderData<typeof clientLoader>();
	invariant(summary, 'Summary data is required');

	return (
		<>
			<StatsHeader current={summary.current} previous={summary.previous} />
			<main>
				<Filters />
				{summary.interval && (
					<StackedBarChart
						data={summary.interval.map((item) => ({
							date: item.date,
							value: item.visitors,
							stackValue: item.pageviews,
						}))}
					/>
				)}
				<Outlet />
			</main>
		</>
	);
}
