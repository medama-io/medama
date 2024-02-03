import {
	json,
	type LoaderFunctionArgs,
	type MetaFunction,
} from '@remix-run/node';
import { Outlet, useLoaderData } from '@remix-run/react';
import invariant from 'tiny-invariant';

import { Filters } from '@/components/stats/Filter';
import { StatsHeader } from '@/components/stats/StatsHeader';
import { fetchStats } from '@/utils/stats';

export const meta: MetaFunction = () => {
	return [{ title: 'Dashboard | Medama' }];
};

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
	const stats = await fetchStats(request, params, {
		dataset: ['summary'],
	});

	return json(
		{
			status: 200,
			...stats,
		},
		{
			headers: {
				'Cache-Control': 'private, max-age=10',
			},
		}
	);
};

export default function Index() {
	const { summary } = useLoaderData<typeof loader>();

	invariant(summary, 'Summary data is required');
	return (
		<>
			<StatsHeader current={summary.current} previous={summary.previous} />
			<main>
				<Filters />
				<div>Chart</div>
				<Outlet />
			</main>
		</>
	);
}
