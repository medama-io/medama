import { userLoggedIn } from '@/api/user';
import { InnerHeader } from '@/components/layout/InnerHeader';
import { SettingsLayout } from '@/components/settings/Layout';
import { Outlet, useLoaderData } from '@remix-run/react';

export const clientLoader = async () => {
	await userLoggedIn();
	return null;
};

export default function Index() {
	useLoaderData();
	return (
		<>
			<InnerHeader>
				<h1>Settings</h1>
			</InnerHeader>
			<main>
				<SettingsLayout>
					<Outlet />
				</SettingsLayout>
			</main>
		</>
	);
}
