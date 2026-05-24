import { Table } from '@/components/stats/Table';
import type { Dataset } from '@/components/stats/types';
import { fetchStats, isDatasetItem } from '@/utils/stats';
import type { Route } from './+types/$hostname.$query';

export const meta: Route.MetaFunction = () => {
	return [{ title: 'Dashboard | Medama' }];
};

const isDataset = (value: string | undefined): value is Dataset =>
	Boolean(value && value !== 'summary' && isDatasetItem(value));

export const clientLoader = async ({
	request,
	params,
}: Route.ClientLoaderArgs) => {
	const query = params.query;
	if (!isDataset(query)) {
		throw new Error('Invalid dataset item');
	}

	const stats = await fetchStats(request, params, { dataset: [query] });

	return {
		query,
		data: stats[query] ?? [],
	};
};

export default function Index({ loaderData }: Route.ComponentProps) {
	return <Table query={loaderData.query} data={loaderData.data} />;
}
