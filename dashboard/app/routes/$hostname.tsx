import {
	json,
	type MetaFunction,
	type ClientLoaderFunctionArgs,
	Outlet,
	useLoaderData,
} from '@remix-run/react';

import { StackedBarChart } from '@/components/stats/Chart';
import { Filters } from '@/components/stats/Filter';
import { StatsHeader } from '@/components/stats/StatsHeader';
import { fetchStats } from '@/utils/stats';
import { userLoggedIn } from '@/api/user';

export const meta: MetaFunction = () => {
	return [{ title: 'Dashboard | Medama' }];
};

export const clientLoader = async ({
	request,
	params,
}: ClientLoaderFunctionArgs) => {
	await userLoggedIn();
	const stats = await fetchStats(request, params, {
		dataset: ['summary'],
	});

	return json(stats);
};

export default function Index() {
	const { summary } = useLoaderData<typeof clientLoader>();
	if (!summary) throw new Error('Summary data is required');

	const chartData =
		summary.interval?.map((item) => ({
			date: item.date,
			value: item.visitors,
			stackValue: item.pageviews,
		})) || [];

	return (
		<>
			<StatsHeader current={summary.current} previous={summary.previous} />
			<main>
				<Filters />
				{summary.interval && <StackedBarChart data={chartData} />}
				<Outlet />
			</main>
		</>
	);
}
