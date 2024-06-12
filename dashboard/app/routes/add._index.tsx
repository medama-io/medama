import { userGet } from '@/api/user';
import { websiteCreate } from '@/api/websites';
import { Add } from '@/components/add/Add';
import { hasSession } from '@/utils/cookies';
import {
	type ClientActionFunctionArgs,
	type ClientLoaderFunctionArgs,
	json,
	type MetaFunction,
	redirect,
} from '@remix-run/react';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Add Website | Medama' },
		{ name: 'description', content: 'Add a website to Medama Analytics.' },
	];
};

export const clientLoader = async ({ request }: ClientLoaderFunctionArgs) => {
	// If the user is already logged in, redirect them to the dashboard.
	if (hasSession()) {
		// Check if session hasn't been revoked
		await userGet();
	} else {
		// Otherwise, redirect them to the login page.
		return redirect('/login');
	}

	return { status: 200 };
};

export const clientAction = async ({ request }: ClientActionFunctionArgs) => {
	const body = await request.formData();

	const hostname = body.get('hostname')
		? String(body.get('hostname'))
		: undefined;
	const name = body.get('name') ? String(body.get('name')) : hostname;

	if (!hostname) {
		throw json('Missing hostname', {
			status: 400,
		});
	}

	const { data, res } = await websiteCreate({
		body: {
			name,
			hostname,
		},
	});

	if (!data) {
		throw json('Failed to create website.', {
			status: res.status,
		});
	}

	return redirect(`/${data.hostname}`);
};

export default function Index() {
	return (
		<main>
			<Add />
		</main>
	);
}
