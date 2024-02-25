import {
	json,
	type LoaderFunctionArgs,
	type MetaFunction,
} from '@remix-run/node';
import { Outlet, useLoaderData, useSearchParams } from '@remix-run/react';
import { format, parseISO } from 'date-fns';
import invariant from 'tiny-invariant';

import { StackedBarChart } from '@/components/stats/Chart';
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

	return json(stats);
};

export default function Index() {
	const { summary } = useLoaderData<typeof loader>();
	const [searchParams] = useSearchParams();

	invariant(summary, 'Summary data is required');

	// eslint-disable-next-line unicorn/consistent-function-scoping
	let dateFormatter = (date: Date) => format(date, 'MMM, yyyy');
	switch (searchParams.get('period')) {
		case 'today':
		case 'yesterday':
		case '12h':
		case '24h':
		case '72h': {
			dateFormatter = (date: Date) => format(date, 'HH:mm');
			break;
		}
		case '7d':
		case '14d':
		case '30d': {
			dateFormatter = (date: Date) => format(date, 'EEEEEE MMM d');
			break;
		}
	}

	return (
		<>
			<StatsHeader current={summary.current} previous={summary.previous} />
			<main>
				<Filters />
				{summary.interval && (
					<StackedBarChart
						data={summary.interval.map((item) => ({
							label: dateFormatter(parseISO(item.date)),
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
