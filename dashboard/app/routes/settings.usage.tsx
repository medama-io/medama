import {
	json,
	useLoaderData,
	useRevalidator,
	type MetaFunction,
} from '@remix-run/react';

import { Section, SectionWrapper } from '@/components/settings/Section';
import { settingsResources } from '@/api/settings';
import { SimpleGrid } from '@mantine/core';
import {
	ResourcePanel,
	ResourcePanelCPU,
	ResourcePanelMetadata,
} from '@/components/settings/Resource';
import { useInterval } from '@mantine/hooks';
import { useEffect } from 'react';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Usage Settings | Medama' },
		{ name: 'description', content: 'Privacy focused web analytics.' },
	];
};

export const clientLoader = async () => {
	const { data } = await settingsResources();

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
	const interval = useInterval(revalidator.revalidate, 2000);
	const { cpu, memory, disk, metadata } = usage;

	useEffect(() => {
		interval.start();
		return interval.stop;
	}, [interval.start, interval.stop]);

	return (
		<>
			<SectionWrapper>
				<SimpleGrid cols={{ base: 1, lg: 2 }}>
					<ResourcePanelCPU title="CPU Usage" {...cpu} />
					<ResourcePanel title="Memory Usage" {...memory} />
					<ResourcePanel title="Disk Usage" {...disk} />
					<ResourcePanelMetadata
						title="Metadata"
						sqlite={metadata.meta_db_version}
						duckdb={metadata.analytics_db_version}
					/>
				</SimpleGrid>
			</SectionWrapper>
			<Section title="Resource management" description="Change program limits">
				<p>Test</p>
			</Section>
		</>
	);
}
