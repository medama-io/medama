import { userLoggedIn } from '@/api/user';
import { websiteCreate } from '@/api/websites';
import { Add } from '@/components/add/Add';
import {
	type ClientActionFunctionArgs,
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

export const clientLoader = async () => {
	await userLoggedIn();
	return null;
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
