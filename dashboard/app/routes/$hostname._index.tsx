import { type LoaderFunctionArgs, type MetaFunction } from '@remix-run/node';
import { useLoaderData } from '@remix-run/react';

import { statsPages, statsSummary } from '@/api/stats';

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
	const { data: summary } = await statsSummary({
		cookie: request.headers.get('Cookie'),
		pathKey: params.hostname,
	});

	const { data: pages } = await statsPages({
		cookie: request.headers.get('Cookie'),
		pathKey: params.hostname,
	});

	if (!summary) {
		throw new Response('Failed to get stats.', {
			status: 500,
		});
	}

	return { status: 200, summary, pages };
};

export const meta: MetaFunction<typeof loader> = () => {
	return [{ title: 'Dashboard | Medama' }];
};

export default function Index() {
	const { summary, pages } = useLoaderData<typeof loader>();
	return (
		<div>
			<h1>Summary</h1>
			{JSON.stringify(summary)}
			<h1>Pages</h1>
			{JSON.stringify(pages)}
		</div>
	);
}
