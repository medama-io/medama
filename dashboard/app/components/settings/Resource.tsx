import { Group, Progress, Stack, Text } from '@mantine/core';
import byteSize from 'byte-size';

import classes from './Resource.module.css';

interface PanelProps {
	title: string;
	usage: number;
}

interface ResourcePanelProps {
	title: string;
	used: number;
	total: number;
}

interface ResourcePanelCPUProps extends PanelProps {
	cores: number;
	threads: number;
}

interface ResourcePanelMetadataProps {
	title: string;
	sqlite?: string;
	duckdb?: string;
}

const Panel = ({
	title,
	usage,
	children,
}: React.PropsWithChildren<PanelProps>) => (
	<div className={classes.panel}>
		<Group justify="space-between">
			<Text fw={500}>{title}</Text>
			<Text fz={22}>{usage.toFixed(2)}%</Text>
		</Group>
		<Progress.Root size="xl" mb={16}>
			<Progress.Section
				value={Math.round(usage)}
				color="#9d5def"
				aria-label={`${title} Progress`}
			/>
		</Progress.Root>
		{children}
	</div>
);

export const ResourcePanel = ({ title, used, total }: ResourcePanelProps) => {
	const usage = (used / total) * 100;
	return (
		<Panel title={title} usage={usage}>
			<Group justify="space-between">
				<Text>
					Used: <span>{String(byteSize(used))}</span>
				</Text>
				<Text>
					Capacity: <span>{String(byteSize(total))}</span>
				</Text>
			</Group>
		</Panel>
	);
};

export const ResourcePanelCPU = ({
	title,
	usage,
	cores,
	threads,
}: ResourcePanelCPUProps) => (
	<Panel title={title} usage={usage}>
		<Group justify="space-between">
			<Text>
				Cores: <span>{String(cores)}</span>
			</Text>
			<Text>
				Threads: <span>{String(threads)}</span>
			</Text>
		</Group>
	</Panel>
);

export const ResourcePanelMetadata = ({
	title,
	sqlite,
	duckdb,
}: ResourcePanelMetadataProps) => (
	<div className={classes.panel}>
		<Text fw={500}>{title}</Text>
		<Group justify="space-between">
			<Text>
				SQLite: <span>{sqlite}</span>
			</Text>
			<Text>
				DuckDB: <span>{duckdb}</span>
			</Text>
		</Group>
	</div>
);
