import {
	json,
	useLoaderData,
	useParams,
	type ClientLoaderFunctionArgs,
	type MetaFunction,
} from '@remix-run/react';

import { StatsTable } from '@/components/stats/StatsTable';
import type { DataRow } from '@/components/stats/types';
import { fetchStats, isDatasetItem, type DatasetItem } from '@/utils/stats';

export const meta: MetaFunction = () => {
	return [{ title: 'Dashboard | Medama' }];
};

export const clientLoader = async ({
	request,
	params,
}: ClientLoaderFunctionArgs) => {
	const query = params.query as DatasetItem;
	if (!query || query === 'summary' || !isDatasetItem(query)) {
		throw new Error('Invalid dataset item');
	}

	const stats = await fetchStats(request, params, { dataset: [query] });

	return json(stats);
};

export default function Index() {
	const params = useParams();
	const data = useLoaderData<Omit<typeof clientLoader, 'summary'>>();

	// We can safely assume that the dataset items are present as the loader function
	// has already validated the query parameter
	const query = params.query as keyof typeof data;
	const stats = data[query] as DataRow[];

	return <StatsTable query={query} data={stats} />;
}
