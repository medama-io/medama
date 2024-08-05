import { SimpleGrid } from '@mantine/core';
import {
	type MetaFunction,
	json,
	useLoaderData,
	useRevalidator,
} from '@remix-run/react';
import { useEffect } from 'react';

import { userUsageGet } from '@/api/user';
import {
	ResourcePanel,
	ResourcePanelCPU,
} from '@/components/settings/Resource';
import { SectionWrapper } from '@/components/settings/Section';
import { useInterval } from '@/hooks/use-interval';

export const meta: MetaFunction = () => {
	return [{ title: 'Usage Settings | Medama' }];
};

export const clientLoader = async () => {
	const { data } = await userUsageGet();

	if (!data) {
		throw json('Failed to get server usage metrics.', {
			status: 500,
		});
	}

	return json({
		usage: data,
	});
};

export default function Index() {
	const { usage } = useLoaderData<typeof clientLoader>();
	const revalidator = useRevalidator();
	const interval = useInterval(revalidator.revalidate, 2500);
	const { cpu, memory, disk } = usage;

	useEffect(() => {
		interval.start();
		return interval.stop;
	}, [interval.start, interval.stop]);

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
