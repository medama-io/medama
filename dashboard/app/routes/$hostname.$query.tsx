import {
	json,
	type LoaderFunctionArgs,
	type MetaFunction,
} from '@remix-run/node';
import { useLoaderData, useParams } from '@remix-run/react';
import invariant from 'tiny-invariant';

import { StatsTable } from '@/components/stats/StatsTable';
import { type DatasetItem, fetchStats, isDatasetItem } from '@/utils/stats';

export const meta: MetaFunction = () => {
	return [{ title: 'Dashboard | Medama' }];
};

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
	const query = params.query as DatasetItem;
	invariant(
		!query || (isDatasetItem(query) && params.query !== 'summary'),
		'Invalid dataset item'
	);

	const stats = await fetchStats(request, params, {
		dataset: [query],
	});

	return json(stats);
};

export default function Index() {
	const params = useParams();
	const data = useLoaderData<Omit<typeof loader, 'summary'>>();

	// We can safely assume that the dataset items are present as the loader function
	// has already validated the query parameter
	const query = params.query as keyof typeof data;
	const stats = data[query];

	if (!stats || typeof stats === 'number') {
		return <div>No data found.</div>;
	}

	return <StatsTable query={query} data={stats} />;
}
