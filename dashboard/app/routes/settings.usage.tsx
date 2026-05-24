import { SimpleGrid } from '@mantine/core';
import { useInterval } from '@mantine/hooks';
import { data as json, useRevalidator } from 'react-router';

import { userUsageGet } from '@/api/user';
import {
	ResourcePanel,
	ResourcePanelCPU,
} from '@/components/settings/Resource';
import { SectionWrapper } from '@/components/settings/Section';
import type { Route } from './+types/settings.usage';

export const meta: Route.MetaFunction = () => {
	return [{ title: 'Usage Settings | Medama' }];
};

export const clientLoader = async () => {
	const { data } = await userUsageGet();

	if (!data) {
		throw json('Failed to get server usage metrics.', {
			status: 500,
		});
	}

	return {
		usage: data,
	};
};

export default function Index({ loaderData }: Route.ComponentProps) {
	const { usage } = loaderData;
	const revalidator = useRevalidator();
	useInterval(revalidator.revalidate, 2500, { autoInvoke: true });
	const { cpu, memory, disk } = usage;

	return (
		<SectionWrapper>
			<SimpleGrid cols={{ base: 1, lg: 2 }}>
				<ResourcePanelCPU title="CPU Usage" {...cpu} />
				<ResourcePanel title="Memory Usage" {...memory} />
				<ResourcePanel title="Disk Usage" {...disk} />
			</SimpleGrid>
		</SectionWrapper>
	);
}
