import type { MetaFunction } from '@remix-run/node';

import { Login } from '@/components/login/Login';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Login | Medama' },
		{ name: 'description', content: 'Login into Medama Analytics.' },
	];
};

export default function Index() {
	return <Login />;
}
