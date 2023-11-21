import { type LoaderFunctionArgs, type MetaFunction } from '@remix-run/node';
import { useLoaderData } from '@remix-run/react';

import { statsSummary } from '@/api/stats';

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
	const { data } = await statsSummary({
		cookie: request.headers.get('Cookie'),
		pathKey: params.hostname,
	});

	if (!data) {
		throw new Response('Failed to get stats.', {
			status: 500,
		});
	}

	return { status: 200, data };
};

export const meta: MetaFunction<typeof loader> = () => {
	return [{ title: 'Dashboard | Medama' }];
};

export default function Index() {
	const { data } = useLoaderData<typeof loader>();
	return <>{JSON.stringify(data)}</>;
}
