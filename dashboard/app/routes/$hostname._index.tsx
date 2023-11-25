import { type LoaderFunctionArgs, type MetaFunction } from '@remix-run/node';
import { useLoaderData } from '@remix-run/react';

import { statsPages, statsSummary, statsTime } from '@/api/stats';

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
	const { data: summary } = await statsSummary({
		cookie: request.headers.get('Cookie'),
		pathKey: params.hostname,
	});

	const { data: pages } = await statsPages({
		cookie: request.headers.get('Cookie'),
		pathKey: params.hostname,
	});
	const { data: pagesSummary } = await statsPages({
		cookie: request.headers.get('Cookie'),
		pathKey: params.hostname,
		query: {
			summary: true,
		},
	});

	const { data: time } = await statsTime({
		cookie: request.headers.get('Cookie'),
		pathKey: params.hostname,
	});
	const { data: timeSummary } = await statsTime({
		cookie: request.headers.get('Cookie'),
		pathKey: params.hostname,
		query: {
			summary: true,
		},
	});

	if (!summary) {
		throw new Response('Failed to get stats.', {
			status: 500,
		});
	}

	return { status: 200, summary, pages, pagesSummary, time, timeSummary };
};

export const meta: MetaFunction<typeof loader> = () => {
	return [{ title: 'Dashboard | Medama' }];
};

export default function Index() {
	const { summary, pages, pagesSummary, time, timeSummary } =
		useLoaderData<typeof loader>();
	return (
		<div>
			<h1>Summary</h1>
			{JSON.stringify(summary, undefined, 2)}
			<h1>Pages</h1>
			<p>Summary</p>
			{JSON.stringify(pagesSummary, undefined, 2)}
			<p>Full</p>
			{JSON.stringify(pages, undefined, 2)}
			<h1>Time</h1>
			<p>Summary</p>
			{JSON.stringify(timeSummary, undefined, 2)}
			<p>Full</p>
			{JSON.stringify(time, undefined, 2)}
		</div>
	);
}
