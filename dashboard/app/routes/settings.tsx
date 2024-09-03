import { Outlet } from '@remix-run/react';

import { InnerHeader } from '@/components/layout/InnerHeader';
import { SettingsLayout } from '@/components/settings/Layout';

export default function Index() {
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
