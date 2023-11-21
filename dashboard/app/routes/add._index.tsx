import {
	type ActionFunctionArgs,
	json,
	type LoaderFunctionArgs,
	type MetaFunction,
	redirect,
} from '@remix-run/node';

import { userGet } from '@/api/user';
import { websiteCreate } from '@/api/websites';
import { Add } from '@/components/add/Add';
import { hasSession } from '@/utils/cookies';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Add Website | Medama' },
		{ name: 'description', content: 'Add a website to Medama Analytics.' },
	];
};

export const loader = async ({ request }: LoaderFunctionArgs) => {
	// If the user is already logged in, redirect them to the dashboard.
	if (hasSession(request)) {
		// Check if session hasn't been revoked
		await userGet({ cookie: request.headers.get('Cookie') });
	}

	return { status: 200 };
};

export const action = async ({ request }: ActionFunctionArgs) => {
	const body = await request.formData();

	const name = body.get('name')?.toString();
	const hostname = body.get('hostname')?.toString();

	if (!hostname) {
		throw json('Missing hostname', {
			status: 400,
		});
	}

	const { data, res } = await websiteCreate({
		cookie: request.headers.get('Cookie'),
		body: {
			name: name ?? hostname,
			hostname,
		},
	});

	if (!data) {
		throw json('Failed to create website.', {
			status: res.status,
		});
	}

	return redirect(`/websites/${data.hostname}`);
};

export default function Index() {
	return <Add />;
}
